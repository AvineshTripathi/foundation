// Copyright 2022 Namespace Labs Inc; All rights reserved.
// Licensed under the EARLY ACCESS SOFTWARE LICENSE AGREEMENT
// available at http://github.com/namespacelabs/foundation

package devworkflow

import (
	"context"
	"fmt"

	"github.com/gorilla/mux"
	"namespacelabs.dev/foundation/engine/compute"
	"namespacelabs.dev/foundation/internal/artifacts/oci"
	"namespacelabs.dev/foundation/languages/web"
	"namespacelabs.dev/foundation/schema"
)

const (
	WebPackage schema.PackageName = "namespacelabs.dev/foundation/devworkflow/web"

	baseRepository = "us-docker.pkg.dev/foundation-344819/prebuilts"
	prebuilt       = "sha256:7c23463fd307825ab082152a527e4863773e3d513699afd8affb1868aa0172c4"
)

func PrebuiltWebUI(ctx context.Context) (*mux.Router, error) {
	image := oci.ImageP(fmt.Sprintf("%s/%s@%s", baseRepository, WebPackage, prebuilt), nil, oci.ResolveOpts{PublicImage: true})

	return compute.GetValue(ctx, web.ServeFS(image, true))
}
