// Copyright 2022 Namespace Labs Inc; All rights reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.

package dockerfile

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/moby/buildkit/client/llb"
	dockerfile "github.com/moby/buildkit/frontend/dockerfile/builder"
	"github.com/moby/buildkit/frontend/dockerfile/dockerignore"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
	"namespacelabs.dev/foundation/internal/artifacts/oci"
	"namespacelabs.dev/foundation/internal/build"
	"namespacelabs.dev/foundation/internal/build/buildkit"
	"namespacelabs.dev/foundation/internal/compute"
	"namespacelabs.dev/foundation/internal/fnerrors"
	"namespacelabs.dev/foundation/internal/parsing/platform"
	"namespacelabs.dev/foundation/internal/wscontents"
	"namespacelabs.dev/foundation/std/pkggraph"
	"namespacelabs.dev/foundation/std/tasks"
)

// A Dockerfile build is always relative to the module it lives in. Within that
// module, what's the relative path to the context, and within that context,
// what's the relative path to the Dockerfile.
func Build(rel, dockerfile string, isFocus bool) (build.Spec, error) {
	return dockerfileBuild{rel, dockerfile, isFocus}, nil
}

type dockerfileBuild struct {
	ContextRel string // Relative to the workspace.
	Dockerfile string // Relative to ContextRel.
	IsFocus    bool
}

var _ build.Spec = dockerfileBuild{}

func makeDockerfileState(sourceLabel string, contents []byte) llb.State {
	return llb.Scratch().
		File(llb.Mkfile("/Dockerfile", 0644,
			contents,
			llb.WithCreatedTime(build.FixedPoint)),
			llb.WithCustomName(fmt.Sprintf("Dockerfile (%s)", sourceLabel)))
}

func (df dockerfileBuild) BuildImage(ctx context.Context, env pkggraph.SealedContext, conf build.Configuration) (compute.Computable[oci.Image], error) {
	// There's a compromise here: we go through a non-snapshot to fetch
	// .dockerignore, to avoid creating two snapshots.
	dfignore, err := fs.ReadFile(conf.Workspace().ReadOnlyFS(), filepath.Join(df.ContextRel, ".dockerignore"))
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return nil, fnerrors.InternalError("failed to check if a .dockerignore exists: %w", err)
		}
	}

	excludes, err := dockerignore.ReadAll(bytes.NewReader(dfignore))
	if err != nil {
		return nil, fnerrors.New("failed to parse dockerignore: %w", err)
	}

	generatedRequest := &generateRequest{
		// Setting observeChanges to true will yield a new solve() on changes to the workspace.
		// Also importantly we scope observe changes to ContextRel.
		workspace:  conf.Workspace().Snapshot(df.ContextRel, df.IsFocus),
		contextRel: df.ContextRel,
		dockerfile: df.Dockerfile,
		conf:       conf,
		excludes:   excludes,
	}

	return buildkit.MakeImage(
		buildkit.DeferClient(env.Configuration(), conf.TargetPlatform()),
		conf,
		generatedRequest,
		[]buildkit.LocalContents{
			dockerContext(conf, df.ContextRel, excludes),
		}, conf.PublishName()), nil
}

func (df dockerfileBuild) PlatformIndependent() bool { return false }

type generateRequest struct {
	workspace              compute.Computable[wscontents.Versioned] // Used as an input so we trigger new requests on changes to the Dockerfile.
	contextRel, dockerfile string
	conf                   build.Configuration
	excludes               []string
	compute.LocalScoped[*buildkit.FrontendRequest]
}

var _ compute.Computable[*buildkit.FrontendRequest] = &generateRequest{}

func (g *generateRequest) Action() *tasks.ActionEvent {
	return tasks.Action("buildkit.make-dockerfile-request").
		Arg("module_name", g.conf.Workspace().ModuleName()).
		Arg("rel", g.contextRel).
		LogLevel(1)
}
func (g *generateRequest) Inputs() *compute.In {
	return compute.Inputs().
		Computable("workspace", g.workspace).
		Str("contextRel", g.contextRel).
		Str("dockerfile", g.dockerfile).
		Indigestible("conf", g.conf)
}
func (g *generateRequest) Output() compute.Output { return compute.Output{NotCacheable: true} }
func (g *generateRequest) Compute(ctx context.Context, deps compute.Resolved) (*buildkit.FrontendRequest, error) {
	workspace := compute.MustGetDepValue(deps, g.workspace, "workspace").FS()

	dfcontents, err := fs.ReadFile(workspace, g.dockerfile)
	if err != nil {
		return nil, err
	}

	req := &buildkit.FrontendRequest{
		Frontend: "dockerfile.v0",
		FrontendInputs: map[string]llb.State{
			dockerfile.DefaultLocalNameDockerfile: makeDockerfileState(g.conf.SourceLabel(), dfcontents),
			dockerfile.DefaultLocalNameContext:    buildkit.MakeLocalState(dockerContext(g.conf, g.contextRel, g.excludes)),
		},
	}

	if g.conf.TargetPlatform() != nil {
		req.FrontendOpt = makeDockerOpts([]specs.Platform{*g.conf.TargetPlatform()})
	}

	return req, nil
}

func dockerContext(conf build.Configuration, contextRel string, excludes []string) buildkit.LocalContents {
	return buildkit.LocalContents{
		Module:          conf.Workspace(),
		Path:            contextRel,
		ObserveChanges:  false, // We don't need to re-observe changes because changes to the workspace will already yield a new frontendReq.
		ExcludePatterns: excludes,
	}
}

func makeDockerOpts(platforms []specs.Platform) map[string]string {
	return map[string]string{
		"platform": formatPlatforms(platforms),
	}
}

func formatPlatforms(ps []specs.Platform) string {
	strs := make([]string, len(ps))
	for k, p := range ps {
		strs[k] = platform.FormatPlatform(p)
	}
	return strings.Join(strs, ",")
}
