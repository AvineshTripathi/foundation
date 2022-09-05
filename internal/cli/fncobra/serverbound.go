// Copyright 2022 Namespace Labs Inc; All rights reserved.
// Licensed under the EARLY ACCESS SOFTWARE LICENSE AGREEMENT
// available at http://github.com/namespacelabs/foundation

package fncobra

import (
	"context"

	"github.com/spf13/cobra"
	"namespacelabs.dev/foundation/internal/fnerrors"
	"namespacelabs.dev/foundation/provision"
	"namespacelabs.dev/foundation/std/pkggraph"
	"namespacelabs.dev/foundation/std/planning"
	"namespacelabs.dev/foundation/workspace"
	"namespacelabs.dev/foundation/workspace/tasks"
)

type Servers struct {
	Servers        []provision.Server
	SealedPackages pkggraph.SealedPackageLoader
}

func ParseServers(serversOut *Servers, env *planning.Context, locs *Locations) *ServersParser {
	return &ServersParser{
		serversOut: serversOut,
		locs:       locs,
		env:        env,
	}
}

type ServersParser struct {
	serversOut *Servers
	locs       *Locations
	env        *planning.Context
}

func (p *ServersParser) AddFlags(cmd *cobra.Command) {}

func (p *ServersParser) Parse(ctx context.Context, args []string) error {
	if p.serversOut == nil {
		return fnerrors.InternalError("serversOut must be set")
	}
	if p.locs == nil {
		return fnerrors.InternalError("locs must be set")
	}
	if p.env == nil {
		return fnerrors.InternalError("env must be set")
	}

	var servers []provision.Server
	pl := workspace.NewPackageLoader(*p.env)
	for _, loc := range p.locs.Locs {
		if err := tasks.Action("package.load-server").Scope(loc.AsPackageName()).Run(ctx, func(ctx context.Context) error {
			pp, err := pl.LoadByName(ctx, loc.AsPackageName())
			if err != nil {
				return fnerrors.Wrap(loc, err)
			}

			if pp.Server == nil {
				if p.locs.AreSpecified {
					return fnerrors.UserError(loc, "expected a server")
				}

				return nil
			}

			server, err := provision.RequireServerWith(ctx, *p.env, pl, loc.AsPackageName())
			if err != nil {
				return err
			}

			if !p.locs.AreSpecified {
				if server.Package.Server.Testonly || server.Package.Server.ClusterAdmin {
					return nil
				}
			}

			servers = append(servers, server)
			return nil
		}); err != nil {
			return err
		}
	}

	*p.serversOut = Servers{
		Servers:        servers,
		SealedPackages: pl.Seal(),
	}

	return nil
}
