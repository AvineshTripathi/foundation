// Copyright 2022 Namespace Labs Inc; All rights reserved.
// Licensed under the EARLY ACCESS SOFTWARE LICENSE AGREEMENT
// available at http://github.com/namespacelabs/foundation

package snapshot

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"namespacelabs.dev/foundation/internal/bytestream"
	"namespacelabs.dev/foundation/internal/compute"
	"namespacelabs.dev/foundation/internal/console"
	"namespacelabs.dev/foundation/internal/filewatcher"
	"namespacelabs.dev/foundation/internal/fnerrors"
	"namespacelabs.dev/foundation/internal/fnfs"
	"namespacelabs.dev/foundation/internal/parsing"
	"namespacelabs.dev/foundation/internal/planning"
	"namespacelabs.dev/foundation/internal/wscontents"
	"namespacelabs.dev/foundation/schema"
	"namespacelabs.dev/foundation/std/cfg"
	"namespacelabs.dev/foundation/std/pkggraph"

	"namespacelabs.dev/foundation/std/tasks"
)

func RequireServers(env cfg.Context, servers ...schema.PackageName) compute.Computable[*ServerSnapshot] {
	return &requiredServers{env: env, packages: servers}
}

type requiredServers struct {
	env      cfg.Context
	packages []schema.PackageName

	compute.LocalScoped[*ServerSnapshot]
}

type ServerSnapshot struct {
	servers []planning.Server
	sealed  pkggraph.SealedPackageLoader
	// Used in Observe()
	env      cfg.Context
	packages []schema.PackageName
}

var _ compute.Versioned = &ServerSnapshot{}

func (rs *requiredServers) Action() *tasks.ActionEvent {
	return tasks.Action("planning.Require-servers")
}

func (rs *requiredServers) Inputs() *compute.In {
	return compute.Inputs().Indigestible("env", rs.env).Strs("packages", schema.Strs(rs.packages...))
}

func (rs *requiredServers) Output() compute.Output {
	return compute.Output{NotCacheable: true}
}

func (rs *requiredServers) Compute(ctx context.Context, _ compute.Resolved) (*ServerSnapshot, error) {
	return computeSnapshot(ctx, rs.env, rs.packages)
}

func computeSnapshot(ctx context.Context, env cfg.Context, packages []schema.PackageName) (*ServerSnapshot, error) {
	pl := parsing.NewPackageLoader(env)

	var servers []planning.Server
	for _, pkg := range packages {
		server, err := planning.RequireServerWith(ctx, env, pl, schema.PackageName(pkg))
		if err != nil {
			return nil, err
		}

		servers = append(servers, server)
	}

	return &ServerSnapshot{servers: servers, sealed: pl.Seal(), env: env, packages: packages}, nil
}

func (snap *ServerSnapshot) Get(pkgs ...schema.PackageName) ([]planning.Server, error) {
	var servers []planning.Server

start:
	for _, pkg := range pkgs {
		for _, srv := range snap.servers {
			if srv.PackageName() == pkg {
				servers = append(servers, srv)
				continue start
			}
		}
		return nil, fnerrors.InternalError("%s: not present in the snapshot", pkg)
	}

	return servers, nil
}

func (snap *ServerSnapshot) Env() pkggraph.Context {
	return pkggraph.MakeSealedContext(snap.env, snap.sealed)
}

func (snap *ServerSnapshot) Equals(rhs *ServerSnapshot) bool {
	return false // XXX optimization.
}

func (snap *ServerSnapshot) Observe(ctx context.Context, onChange func(compute.ResultWithTimestamp[any], bool)) (func(), error) {
	obs := obsState{onChange: onChange}
	if err := obs.prepare(ctx, snap); err != nil {
		return nil, err
	}
	return obs.cancel, nil
}

type obsState struct {
	cancelWatcher func()
	onChange      func(compute.ResultWithTimestamp[any], bool)
}

func (p *obsState) prepare(ctx context.Context, snapshot *ServerSnapshot) error {
	cancel, err := observe(ctx, snapshot, func(newSnapshot *ServerSnapshot) {
		if p.cancelWatcher != nil {
			p.cancelWatcher()
			p.cancelWatcher = nil
		}

		r := compute.ResultWithTimestamp[any]{
			Completed: time.Now(),
		}
		r.Value = newSnapshot

		p.onChange(r, false)

		if err := p.prepare(ctx, newSnapshot); err != nil {
			compute.Stop(ctx, err)
		}
	})
	p.cancelWatcher = cancel
	return err
}

func (p *obsState) cancel() {
	if p.cancelWatcher != nil {
		p.cancelWatcher()
	}
}

func observe(ctx context.Context, snap *ServerSnapshot, onChange func(*ServerSnapshot)) (func(), error) {
	logger := console.TypedOutput(ctx, "observepackages", console.CatOutputUs)

	watcher, err := filewatcher.NewFactory(ctx)
	if err != nil {
		return nil, err
	}

	expected := map[string][]byte{}

	for _, mod := range snap.sealed.Modules() {
		// Don't monitor changes to external modules, assume they're immutable.
		if mod.IsExternal() {
			continue
		}

		// XXX we don't watch directories, which may end up being a miss.
		if err := fnfs.VisitFiles(ctx, mod.ReadOnlyFS(), func(path string, blob bytestream.ByteStream, de fs.DirEntry) error {
			if de.IsDir() {
				return nil
			}
			contents, err := bytestream.ReadAll(blob)
			if err != nil {
				return err
			}
			p := filepath.Join(mod.Abs(), path)
			expected[p] = contents // Don't really care about permissions etc, only contents.
			return watcher.AddFile(p)
		}); err != nil {
			watcher.Close()
			return nil, err
		}
	}

	bufferCh := make(chan []fsnotify.Event)
	go func() {
		for buffer := range bufferCh {
			dirty := false
			for _, ev := range buffer {
				contents, err := os.ReadFile(ev.Name)
				if err == nil {
					if !bytes.Equal(contents, expected[ev.Name]) {
						err = fmt.Errorf("%s: contents differ", ev.Name)
					}
				}
				if err != nil {
					dirty = true
					break
				}
			}

			if !dirty {
				continue
			}

			newSnapshot, err := computeSnapshot(ctx, snap.env, snap.packages)
			if err != nil {
				if msg, ok := fnerrors.IsExpected(err); ok {
					fmt.Fprintf(console.Stderr(ctx), "\n  %s\n\n", msg)
					continue // Swallow the error in case it is expected.
				}
				compute.Stop(ctx, err)
				break
			}

			if newSnapshot.Equals(snap) {
				continue
			}

			fmt.Fprintln(logger, "Recomputed graph.")

			onChange(newSnapshot)
			break // Don't send any more events.
		}

		for range bufferCh {
			// Drain the channel in case we escaped the loop above, so the go-routine below
			// has a chance to observe a canceled context and close the channel.
		}
	}()

	w, err := watcher.StartWatching(ctx)
	if err != nil {
		return nil, err
	}

	go func() {
		wscontents.AggregateFSEvents(w, console.Debug(ctx), logger, bufferCh)
	}()

	return func() {
		w.Close()
	}, nil
}
