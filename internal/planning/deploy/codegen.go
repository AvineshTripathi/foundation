// Copyright 2022 Namespace Labs Inc; All rights reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.

package deploy

import (
	"context"
	"io/fs"

	"namespacelabs.dev/foundation/internal/codegen/genpackage"
	"namespacelabs.dev/foundation/internal/compute"
	"namespacelabs.dev/foundation/internal/fnfs"
	"namespacelabs.dev/foundation/internal/planning"
	"namespacelabs.dev/foundation/internal/wscontents"
	"namespacelabs.dev/foundation/schema"
	"namespacelabs.dev/foundation/std/cfg"
	"namespacelabs.dev/foundation/std/execution"
	"namespacelabs.dev/foundation/std/pkggraph"
	"namespacelabs.dev/foundation/std/tasks"
)

type codegenWorkspace struct {
	srv planning.Server
}

func (cw codegenWorkspace) ModuleName() string             { return cw.srv.Module().ModuleName() }
func (cw codegenWorkspace) Abs() string                    { return cw.srv.Module().Abs() }
func (cw codegenWorkspace) ReadOnlyFS(rel ...string) fs.FS { return cw.srv.Module().ReadOnlyFS(rel...) }
func (cw codegenWorkspace) ChangeTrigger(rel string) compute.Computable[compute.Versioned] {
	return cw.srv.Module().ChangeTrigger(rel)
}
func (cw codegenWorkspace) Snapshot(rel string) compute.Computable[wscontents.Versioned] {
	if cw.srv.Module().IsExternal() {
		return cw.srv.Module().Snapshot(rel)
	}

	return &codegenThenSnapshot{srv: cw.srv, rel: rel}
}

type codegenThenSnapshot struct {
	srv planning.Server
	rel string
	compute.LocalScoped[wscontents.Versioned]
}

func (cd *codegenThenSnapshot) Action() *tasks.ActionEvent {
	return tasks.Action("workspace.codegen-and-snapshot").Scope(cd.srv.PackageName())
}
func (cd *codegenThenSnapshot) Inputs() *compute.In {
	return compute.Inputs().Indigestible("srv", cd.srv).Str("rel", cd.rel)
}
func (cd *codegenThenSnapshot) Compute(ctx context.Context, _ compute.Resolved) (wscontents.Versioned, error) {
	if err := codegenServer(ctx, cd.srv); err != nil {
		return nil, err
	}

	// Codegen is only run once; if codegen is required again, then it will be triggered
	// by a recomputation of the graph.

	return wscontents.MakeVersioned(ctx, cd.srv.Module().Abs(), cd.rel, true, nil)
}

type codegenEnv struct {
	config   cfg.Configuration
	root     *pkggraph.Module
	packages pkggraph.PackageLoader
	env      *schema.Environment
	fs       fnfs.ReadWriteFS
}

var _ pkggraph.ContextWithMutableModule = codegenEnv{}

func (ce codegenEnv) ErrorLocation() string            { return ce.root.ErrorLocation() }
func (ce codegenEnv) Environment() *schema.Environment { return ce.env }
func (ce codegenEnv) ModuleName() string               { return ce.root.ModuleName() }
func (ce codegenEnv) ReadWriteFS() fnfs.ReadWriteFS    { return ce.fs }
func (ce codegenEnv) Configuration() cfg.Configuration { return ce.config }
func (ce codegenEnv) Workspace() cfg.Workspace         { return ce.root.WorkspaceData }

func (ce codegenEnv) Resolve(ctx context.Context, pkg schema.PackageName) (pkggraph.Location, error) {
	return ce.packages.Resolve(ctx, pkg)
}

func (ce codegenEnv) LoadByName(ctx context.Context, packageName schema.PackageName) (*pkggraph.Package, error) {
	return ce.packages.LoadByName(ctx, packageName)
}

func codegenServer(ctx context.Context, srv planning.Server) error {
	// XXX we should be able to disable codegen for pure builds.
	if srv.Module().IsExternal() {
		return nil
	}

	codegen, err := genpackage.ForServerAndDeps(ctx, srv)
	if err != nil {
		return err
	}

	if len(codegen) == 0 {
		return nil
	}

	r := execution.NewPlan(codegen...)

	return execution.Execute(ctx, "workspace.codegen", r, nil,
		execution.FromContext(srv.SealedContext()),
		pkggraph.MutableModuleInjection.With(srv.Module()),
		pkggraph.PackageLoaderInjection.With(srv.SealedContext()),
	)
}
