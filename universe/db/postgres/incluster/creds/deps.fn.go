// This file was automatically generated by Namespace.
// DO NOT EDIT. To update, re-run `ns generate`.

package creds

import (
	"context"
	"namespacelabs.dev/foundation/std/go/core"
	"namespacelabs.dev/foundation/std/secrets"
)

// Dependencies that are instantiated once for the lifetime of the extension.
type ExtensionDeps struct {
	Password *secrets.Value
}

type _checkProvideCreds func(context.Context, *CredsRequest, ExtensionDeps) (*Creds, error)

var _ _checkProvideCreds = ProvideCreds

var (
	Package__9gpcgr = &core.Package{
		PackageName: "namespacelabs.dev/foundation/universe/db/postgres/incluster/creds",
	}

	Provider__9gpcgr = core.Provider{
		Package:     Package__9gpcgr,
		Instantiate: makeDeps__9gpcgr,
	}
)

func makeDeps__9gpcgr(ctx context.Context, di core.Dependencies) (_ interface{}, err error) {
	var deps ExtensionDeps

	// name: "postgres-password-file"
	if deps.Password, err = secrets.ProvideSecret(ctx, core.MustUnwrapProto("ChZwb3N0Z3Jlcy1wYXNzd29yZC1maWxl", &secrets.Secret{}).(*secrets.Secret)); err != nil {
		return nil, err
	}

	return deps, nil
}
