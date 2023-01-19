// Copyright 2022 Namespace Labs Inc; All rights reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.

package kubedef

import (
	"fmt"
	"strings"

	"namespacelabs.dev/foundation/framework/kubernetes/kubenaming"
	"namespacelabs.dev/foundation/internal/runtime"
	schema "namespacelabs.dev/foundation/schema"
)

func MakeDeploymentId(srv runtime.Deployable) string {
	if srv.GetName() == "" {
		return srv.GetId()
	}

	// k8s doesn't accept uppercase names.
	return fmt.Sprintf("%s-%s", strings.ToLower(srv.GetName()), srv.GetId())
}

func MakeVolumeName(v *schema.Volume) string {
	if v.Inline {
		return kubenaming.LabelLike("vi", v.Name)
	}

	return kubenaming.LabelLike("v", v.Name)
}

func MakeResourceName(deploymentId string, suffix ...string) string {
	return kubenaming.DomainFragLike(append([]string{deploymentId}, suffix...)...)
}
