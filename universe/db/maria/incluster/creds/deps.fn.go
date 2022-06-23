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
	Package__bihnv9 = &core.Package{
		PackageName: "namespacelabs.dev/foundation/universe/db/maria/incluster/creds",
	}

	Provider__bihnv9 = core.Provider{
		Package:     Package__bihnv9,
		Instantiate: makeDeps__bihnv9,
	}
)

func makeDeps__bihnv9(ctx context.Context, di core.Dependencies) (_ interface{}, err error) {
	var deps ExtensionDeps

	// name: "mariadb-password-file"
	if deps.Password, err = secrets.ProvideSecret(ctx, core.MustUnwrapProto("ChVtYXJpYWRiLXBhc3N3b3JkLWZpbGU=", &secrets.Secret{}).(*secrets.Secret)); err != nil {
		return nil, err
	}

	return deps, nil
}
