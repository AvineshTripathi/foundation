package nsingress

import (
	"fmt"

	"namespacelabs.dev/foundation/schema"
)

const (
	LocalBaseDomain = "nslocal.host"
	CloudBaseDomain = "nscloud.dev"
)

func ComputeNaming(env *schema.Environment, source *schema.Naming) (*schema.ComputedNaming, error) {
	if env.Purpose != schema.Environment_PRODUCTION {
		return &schema.ComputedNaming{
			Source:     source,
			BaseDomain: LocalBaseDomain,
			Managed:    schema.Domain_LOCAL_MANAGED,
		}, nil
	}

	if !source.GetEnableNamespaceManaged() {
		return &schema.ComputedNaming{}, nil
	}

	org := source.GetWithOrg()
	if org == "" {
		return &schema.ComputedNaming{}, nil
	}

	return &schema.ComputedNaming{
		Source:     source,
		BaseDomain: fmt.Sprintf("%s.%s", org, CloudBaseDomain),
		Managed:    schema.Domain_CLOUD_MANAGED,
	}, nil
}
