// Copyright 2022 Namespace Labs Inc; All rights reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.

package oci

import (
	"context"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"namespacelabs.dev/foundation/internal/build/registry"
	"namespacelabs.dev/foundation/schema"
)

type Keychain interface {
	Resolve(context.Context, authn.Resource) (authn.Authenticator, error)
}

type RepositoryWithParent struct {
	Parent interface{}
	RepositoryWithAccess
}

type RegistryAccess struct {
	InsecureRegistry bool
	PublicImage      bool
	Keychain         Keychain
	Transport        *registry.RegistryTransport
}

type RepositoryWithAccess struct {
	RegistryAccess
	Repository string
}

func (t RepositoryWithParent) ComputeDigest(context.Context) (schema.Digest, error) {
	return schema.DigestOf("insecureRegistry", t.InsecureRegistry, "repository", t.Repository, "transport", t.Transport)
}

func defaultTag(digest v1.Hash) string {
	// Registry protocol requires a tag. go-containerregistry uses "latest" by default.
	// We configute tags in ECR as immutable to harden versioning in deployments.
	// This does not combine well, so we infer a stable tag here.
	// Inferring the tag from the digest also helps to avoid AWS tag limits.
	// https://docs.aws.amazon.com/AmazonECR/latest/userguide/service-quotas.html
	return strings.TrimPrefix(digest.String(), "sha256:")
}

func ParseTag(tag RepositoryWithAccess, digest v1.Hash) (name.Tag, error) {
	var opts []name.Option
	if tag.InsecureRegistry {
		opts = append(opts, name.Insecure)
	}

	opts = append(opts, name.WithDefaultTag(defaultTag(digest)))

	return name.NewTag(tag.Repository, opts...)
}
