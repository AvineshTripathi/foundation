// Copyright 2022 Namespace Labs Inc; All rights reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.

package docker

import (
	"context"

	"namespacelabs.dev/foundation/internal/artifacts/oci"
	"namespacelabs.dev/foundation/internal/compute"
	"namespacelabs.dev/foundation/internal/fnerrors"
	"namespacelabs.dev/foundation/std/tasks"
)

func PublishImage(tag compute.Computable[oci.RepositoryWithParent], image compute.Computable[oci.ResolvableImage]) compute.Computable[oci.ImageID] {
	return &publishImage{tag: tag, image: image}
}

type publishImage struct {
	tag   compute.Computable[oci.RepositoryWithParent]
	image compute.Computable[oci.ResolvableImage]

	compute.LocalScoped[oci.ImageID]
}

func (pi *publishImage) Inputs() *compute.In {
	return compute.Inputs().Computable("tag", pi.tag).Computable("image", pi.image)
}

func (pi *publishImage) Output() compute.Output {
	return compute.Output{NotCacheable: true} // XXX capture more explicitly that there are side-effects.
}

func (pi *publishImage) Action() *tasks.ActionEvent {
	return tasks.Action("docker.publish")
}

func (pi *publishImage) Compute(ctx context.Context, deps compute.Resolved) (oci.ImageID, error) {
	tag := compute.MustGetDepValue(deps, pi.tag, "tag")
	resolvable := compute.MustGetDepValue(deps, pi.image, "image")

	tasks.Attachments(ctx).AddResult("repository", tag.Repository)

	img, err := resolvable.ImageForPlatform(HostPlatform())
	if err != nil {
		return oci.ImageID{}, fnerrors.InternalError("docker: %w", err)
	}

	digest, err := img.Digest()
	if err != nil {
		return oci.ImageID{}, err
	}

	ref, err := oci.ParseTag(tag.RepositoryWithAccess, digest)
	if err != nil {
		return oci.ImageID{}, err
	}

	if err := WriteImage(ctx, img, ref, true); err != nil {
		return oci.ImageID{}, err
	}

	return oci.ImageID{Repository: tag.Repository, Digest: digest.String()}, nil
}
