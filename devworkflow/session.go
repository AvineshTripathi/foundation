// Copyright 2022 Namespace Labs Inc; All rights reserved.
// Licensed under the EARLY ACCESS SOFTWARE LICENSE AGREEMENT
// available at http://github.com/namespacelabs/foundation

package devworkflow

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

	"google.golang.org/protobuf/proto"
	"namespacelabs.dev/foundation/internal/console"
	"namespacelabs.dev/foundation/internal/fnerrors"
	"namespacelabs.dev/foundation/internal/logoutput"
	"namespacelabs.dev/foundation/internal/runtime/endpointfwd"
	"namespacelabs.dev/foundation/internal/syncbuffer"
	"namespacelabs.dev/foundation/provision"
	"namespacelabs.dev/foundation/runtime"
	"namespacelabs.dev/foundation/schema"
	"namespacelabs.dev/foundation/workspace/module"
	"namespacelabs.dev/foundation/workspace/tasks"
	"namespacelabs.dev/foundation/workspace/tasks/protocol"
)

var AlsoOutputToStderr = false
var AlsoOutputBuildToStderr = false
var TaskOutputBuildkitJsonLog = tasks.Output("buildkit.json", "application/json+fn.buildkit")

type Session struct {
	Ch chan *DevWorkflowRequest

	Console   io.Writer
	Errors    io.Writer
	setSticky func([]byte)

	localHostname string
	obs           *Observers
	sink          *tasks.StatefulSink

	commandOutput *syncbuffer.ByteBuffer // XXX cap the size
	buildOutput   *syncbuffer.ByteBuffer // XXX cap the size
	buildkitJSON  *syncbuffer.ByteBuffer

	mu        sync.Mutex // Protect below.
	requested struct {
		absRoot string
		envName string
		servers []string
	}
	cancelWorkspace func()
	currentStack    *Stack
	currentEnv      runtime.Selector
	pfw             *endpointfwd.PortForward
}

func NewSession(ctx context.Context, sink *tasks.StatefulSink, localHostname string, stickies []string) (*Session, error) {
	setSticky := func(b []byte) {
		var out bytes.Buffer
		for _, sticky := range stickies {
			fmt.Fprintf(&out, " %s\n", sticky)
		}
		if len(b) > 0 && len(stickies) > 0 {
			fmt.Fprintln(&out)
			out.Write(b)
		}

		console.SetStickyContent(ctx, "stack", out.Bytes())
	}

	setSticky(nil)

	return &Session{
		Console:       console.TypedOutput(ctx, "fn dev", console.CatOutputUs),
		Errors:        console.Errors(ctx),
		setSticky:     setSticky,
		localHostname: localHostname,
		obs:           NewObservers(ctx),
		Ch:            make(chan *DevWorkflowRequest, 1),
		commandOutput: syncbuffer.NewByteBuffer(),
		buildOutput:   syncbuffer.NewByteBuffer(),
		buildkitJSON:  syncbuffer.NewByteBuffer(),
		sink:          sink,
	}, nil
}

func (s *Session) Close() {
	close(s.Ch)
	s.obs.Close()
}

func (s *Session) NewClient() (chan *Update, func()) {
	ch := make(chan *Update, 1)

	const maxTaskUpload = 1000
	protos := s.sink.History(maxTaskUpload, func(t *protocol.Task) bool {
		return true
	})

	s.mu.Lock()
	// When a new client connects, send them the latest information immediately.
	// XXX keep latest computed stack in `s`.
	tu := &Update{TaskUpdate: protos, StackUpdate: proto.Clone(s.currentStack).(*Stack)}
	s.mu.Unlock()

	ch <- tu

	s.obs.Add(ch)
	return ch, func() {
		s.obs.Remove(ch)
		close(ch)
	}
}

func (s *Session) CommandOutput() io.ReadCloser   { return s.commandOutput.Reader() }
func (s *Session) BuildOutput() io.ReadCloser     { return s.buildOutput.Reader() }
func (s *Session) BuildJSONOutput() io.ReadCloser { return s.buildkitJSON.Reader() }

func (s *Session) ResolveServer(ctx context.Context, serverID string) (runtime.Selector, *schema.Server, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry := s.currentStack.GetStack().GetServerByID(serverID)
	if entry != nil {
		return s.currentEnv, entry.Server, nil
	}

	return nil, nil, fnerrors.UserError(nil, "%s: no such server in the current session", serverID)
}

func (s *Session) handleSetWorkspace(parentCtx context.Context, absRoot, envName string, servers []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cancelWorkspace != nil {
		s.cancelWorkspace() // Cancel whatever it is doing.
		s.cancelWorkspace = nil
	}

	previousPortFwds := s.currentStack.GetForwardedPort()
	s.currentStack = &Stack{ForwardedPort: previousPortFwds}

	s.requested.absRoot = absRoot
	s.requested.envName = envName
	s.requested.servers = servers

	if len(servers) > 0 {
		ctx, newCancel := context.WithCancel(parentCtx)
		s.cancelWorkspace = newCancel

		// Reset the banner.
		s.setSticky(nil)

		env, err := loadWorkspace(ctx, absRoot, envName)
		if err != nil {
			s.cancelPortForward()
			return err
		}

		resetStack(s.currentStack, env)
		pfw := s.setEnvironment(env)

		go func() {
			err := setWorkspace(ctx, env, servers[0], servers[1:], s, pfw)

			if err != nil && !errors.Is(err, context.Canceled) {
				fnerrors.Format(console.Stderr(parentCtx), true, err)
			}
		}()
	}

	return nil
}

func loadWorkspace(ctx context.Context, absRoot, envName string) (provision.Env, error) {
	// Re-create loc/root here, to dump the cache.
	root, err := module.FindRoot(ctx, absRoot)
	if err != nil {
		return provision.Env{}, err
	}

	return provision.RequireEnv(root, envName)
}

type sinkObserver struct{ s *Session }

func (so *sinkObserver) pushUpdate(ra *tasks.RunningAction) {
	p := ra.Proto()

	so.s.obs.Publish(&Update{TaskUpdate: []*protocol.Task{p}})
}

func (so *sinkObserver) OnStart(ra *tasks.RunningAction)  { so.pushUpdate(ra) }
func (so *sinkObserver) OnUpdate(ra *tasks.RunningAction) { so.pushUpdate(ra) }
func (so *sinkObserver) OnDone(ra *tasks.RunningAction)   { so.pushUpdate(ra) }

func (s *Session) Run(ctx context.Context) error {
	cancel := s.sink.Observe(&sinkObserver{s})
	defer cancel()

	writers := []io.Writer{s.commandOutput}

	if AlsoOutputToStderr {
		writers = append(writers, console.Stderr(ctx))
	}

	var w io.Writer
	if len(writers) != 1 {
		w = io.MultiWriter(writers...)
	} else {
		w = writers[0]
	}

	ctx = logoutput.WithOutput(ctx, logoutput.OutputTo{
		Writer:     w,
		WithColors: true, // Assume xterm.js can handle color.
	})

	defer func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		s.cancelPortForward()
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case req, ok := <-s.Ch:
			if !ok {
				return nil
			}

			switch x := req.Type.(type) {
			case *DevWorkflowRequest_SetWorkspace_:
				set := x.SetWorkspace
				servers := append([]string{set.GetPackageName()}, set.GetAdditionalServers()...)
				if err := s.handleSetWorkspace(ctx, set.GetAbsRoot(), set.GetEnvName(), servers); err != nil {
					fmt.Fprintln(console.Errors(ctx), "failed to load workspace", err)
					return err
				}

			case *DevWorkflowRequest_ReloadWorkspace:
				if x.ReloadWorkspace {
					s.mu.Lock()
					absRoot := s.requested.absRoot
					envName := s.requested.envName
					servers := s.requested.servers
					s.mu.Unlock()
					if err := s.handleSetWorkspace(ctx, absRoot, envName, servers); err != nil {
						fmt.Fprintln(console.Errors(ctx), "failed to load workspace", err)
						return err
					}
				}
			}
		}
	}
}

func (s *Session) TaskLogByName(taskID, name string) io.ReadCloser {
	return s.sink.HistoricReaderByName(taskID, name)
}

func (s *Session) setEnvironment(env runtime.Selector) *endpointfwd.PortForward {
	if s.pfw != nil && proto.Equal(s.currentEnv.Proto(), env.Proto()) {
		// Nothing to do.
		return s.pfw
	}

	s.cancelPortForward()

	s.pfw = newPortFwd(s, env, s.localHostname)
	s.currentEnv = env
	return s.pfw
}

func (s *Session) cancelPortForward() {
	if s.pfw != nil {
		if err := s.pfw.Cleanup(); err != nil {
			fmt.Fprintln(s.Errors, "Failed to cleanup port forwarding resources", err)
		}
		s.pfw = nil
	}
}

func (s *Session) updateStackInPlace(f func(stack *Stack)) {
	s.mu.Lock()
	f(s.currentStack)
	s.currentStack.Revision++
	copy := proto.Clone(s.currentStack).(*Stack)
	s.mu.Unlock()

	s.obs.Publish(&Update{StackUpdate: copy})
}

func resetStack(out *Stack, env provision.Env) {
	workspace := proto.Clone(env.Root().Workspace).(*schema.Workspace)

	// XXX handling broken web ui builds.
	if workspace.Env == nil {
		workspace.Env = provision.EnvsOrDefault(workspace)
	}

	out.AbsRoot = env.Root().Abs()
	out.Env = env.Proto()
	out.Workspace = workspace
	out.AvailableEnv = workspace.Env
}
