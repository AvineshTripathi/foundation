// This file was automatically generated by Foundation.
// DO NOT EDIT. To update, re-run `fn generate`.

package localstacks3

import (
	"context"
	"namespacelabs.dev/foundation/std/go/core"
	"namespacelabs.dev/foundation/std/go/server"
	"namespacelabs.dev/foundation/universe/aws/s3"
	s31 "namespacelabs.dev/foundation/universe/storage/s3"
)

// Dependencies that are instantiated once for the lifetime of the service.
type ServiceDeps struct {
	Bucket *s3.Bucket
}

// Verify that WireService is present and has the appropriate type.
type checkWireService func(context.Context, server.Registrar, ServiceDeps)

var _ checkWireService = WireService

var (
	Package__g72tjm = &core.Package{
		PackageName: "namespacelabs.dev/foundation/std/testdata/service/localstacks3",
	}

	Provider__g72tjm = core.Provider{
		Package:     Package__g72tjm,
		Instantiate: makeDeps__g72tjm,
	}
)

func makeDeps__g72tjm(ctx context.Context, di core.Dependencies) (_ interface{}, err error) {
	var deps ServiceDeps

	if err := di.Instantiate(ctx, s31.Provider__4pkegf, func(ctx context.Context, v interface{}) (err error) {
		// bucket_name: "test-foo-bucket"
		if deps.Bucket, err = s31.ProvideBucket(ctx, core.MustUnwrapProto("Cg90ZXN0LWZvby1idWNrZXQ=", &s31.BucketArgs{}).(*s31.BucketArgs), v.(s31.ExtensionDeps)); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return deps, nil
}
