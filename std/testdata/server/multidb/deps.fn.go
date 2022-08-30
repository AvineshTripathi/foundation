// This file was automatically generated by Namespace.
// DO NOT EDIT. To update, re-run `ns generate`.

package main

import (
	"context"
	"namespacelabs.dev/foundation/std/go/core"
	"namespacelabs.dev/foundation/std/go/grpc/metrics"
	"namespacelabs.dev/foundation/std/go/server"
	"namespacelabs.dev/foundation/std/monitoring/tracing"
	"namespacelabs.dev/foundation/std/testdata/service/multidb"
)

func RegisterInitializers(di *core.DependencyGraph) {
	di.AddInitializers(metrics.Initializers__so2f3v...)
	di.AddInitializers(tracing.Initializers__70o2mm...)
}

func WireServices(ctx context.Context, srv server.Server, depgraph core.Dependencies) []error {
	var errs []error

	if err := depgraph.Instantiate(ctx, multidb.Provider__7cco3b, func(ctx context.Context, v interface{}) error {
		multidb.WireService(ctx, srv.Scope(multidb.Package__7cco3b), v.(multidb.ServiceDeps))
		return nil
	}); err != nil {
		errs = append(errs, err)
	}

	return errs
}
