// Copyright 2022 Namespace Labs Inc; All rights reserved.
// Licensed under the EARLY ACCESS SOFTWARE LICENSE AGREEMENT
// available at http://github.com/namespacelabs/foundation

package kubernetes

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/anypb"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	applycorev1 "k8s.io/client-go/applyconfigurations/core/v1"
	"namespacelabs.dev/foundation/internal/artifacts/oci"
	"namespacelabs.dev/foundation/internal/console"
	"namespacelabs.dev/foundation/internal/fnerrors"
	"namespacelabs.dev/foundation/internal/frontend"
	"namespacelabs.dev/foundation/runtime"
	"namespacelabs.dev/foundation/runtime/kubernetes/client"
	"namespacelabs.dev/foundation/runtime/kubernetes/kubedef"
	"namespacelabs.dev/foundation/runtime/kubernetes/kubetool"
	"namespacelabs.dev/foundation/runtime/rtypes"
	"namespacelabs.dev/foundation/schema"
	"namespacelabs.dev/foundation/std/planning"
	"namespacelabs.dev/foundation/workspace/tasks"
)

var (
	ObserveInitContainerLogs = false
)

func Register() {
	runtime.Register("kubernetes", func(ctx context.Context, env planning.Context) (runtime.Class, error) {
		hostConfig, err := client.ComputeHostConfig(env.Configuration())
		if err != nil {
			return nil, err
		}

		fmt.Fprintf(console.Debug(ctx), "kubernetes: selected %+v for %q\n", hostConfig.HostEnv, env.Environment().Name)

		p, err := client.MakeDeferredRuntime(ctx, hostConfig)
		if err != nil {
			return nil, err
		}

		if p != nil {
			return p, nil
		}

		return deferredRuntime{env.Configuration()}, nil
	})

	frontend.RegisterPrepareHook("namespacelabs.dev/foundation/std/runtime/kubernetes.ApplyServerExtensions", prepareApplyServerExtensions)
}

type deferredRuntime struct {
	cfg planning.Configuration
}

var _ runtime.Class = deferredRuntime{}

func (d deferredRuntime) Namespace(env planning.Context) runtime.Namespace {
	return RuntimeClass(env)
}

func (d deferredRuntime) AttachToCluster(ctx context.Context, ns runtime.Namespace) (runtime.Cluster, error) {
	unbound, err := NewCluster(ctx, d.cfg)
	if err != nil {
		return nil, err
	}

	return unbound.Bind(ns)
}

func (d deferredRuntime) EnsureCluster(ctx context.Context) (runtime.DeferredCluster, error) {
	unbound, err := NewCluster(ctx, d.cfg)
	if err != nil {
		return nil, err
	}

	return unbound, nil
}

func RuntimeClass(env planning.Context) planner {
	ns := ModuleNamespace(env.Workspace().Proto(), env.Environment())

	conf := &kubetool.KubernetesEnv{}
	if env.Configuration().Get(conf) {
		ns = conf.Namespace
	}

	return planner{clusterTarget{env: env.Environment(), namespace: ns}}
}

func MakeNamespace(env *schema.Environment, ns string) *applycorev1.NamespaceApplyConfiguration {
	return applycorev1.Namespace(ns).
		WithLabels(kubedef.MakeLabels(env, nil)).
		WithAnnotations(kubedef.MakeAnnotations(env, nil))
}

func (r ClusterNamespace) PrepareProvision(ctx context.Context, _ planning.Context) (*rtypes.ProvisionProps, error) {
	systemInfo, err := r.SystemInfo(ctx)
	if err != nil {
		return nil, err
	}

	return PrepareProvisionWith(r.env, r.namespace, systemInfo)
}

func PrepareProvisionWith(env *schema.Environment, ns string, systemInfo *kubedef.SystemInfo) (*rtypes.ProvisionProps, error) {
	packedHostEnv, err := anypb.New(&kubetool.KubernetesEnv{Namespace: ns})
	if err != nil {
		return nil, err
	}

	packedSystemInfo, err := anypb.New(systemInfo)
	if err != nil {
		return nil, err
	}

	// Ensure the namespace exist, before we go and apply definitions to it. Also, deployServer
	// assumes that a namespace already exists.
	def, err := (kubedef.Apply{
		Description: fmt.Sprintf("Namespace for %q", env.Name),
		Resource:    MakeNamespace(env, ns),
	}).ToDefinition()
	if err != nil {
		return nil, err
	}

	// Pass the computed namespace to the provisioning tool.
	return &rtypes.ProvisionProps{
		ProvisionInput: []*anypb.Any{packedHostEnv, packedSystemInfo},
		Invocation:     []*schema.SerializedInvocation{def},
	}, nil
}

func (r ClusterNamespace) DeployedConfigImageID(ctx context.Context, server *schema.Server) (oci.ImageID, error) {
	return tasks.Return(ctx, tasks.Action("kubernetes.resolve-config-image-id").Scope(schema.PackageName(server.PackageName)),
		func(ctx context.Context) (oci.ImageID, error) {
			// XXX need a StatefulSet variant.
			d, err := r.cli.AppsV1().Deployments(r.namespace).Get(ctx, kubedef.MakeDeploymentId(server), metav1.GetOptions{})
			if err != nil {
				// XXX better error messages.
				return oci.ImageID{}, err
			}

			cfgimage, ok := d.Annotations[kubedef.K8sConfigImage]
			if !ok {
				return oci.ImageID{}, fnerrors.BadInputError("%s: %q is missing as an annotation in %q",
					server.GetPackageName(), kubedef.K8sConfigImage, kubedef.MakeDeploymentId(server))
			}

			imgid, err := oci.ParseImageID(cfgimage)
			if err != nil {
				return imgid, err
			}

			tasks.Attachments(ctx).AddResult("ref", imgid.ImageRef())

			return imgid, nil
		})
}

func (r ClusterNamespace) StartTerminal(ctx context.Context, server *schema.Server, rio runtime.TerminalIO, command string, rest ...string) error {
	cmd := append([]string{command}, rest...)

	return r.startTerminal(ctx, r.cli, server, rio, cmd)
}

func (r ClusterNamespace) AttachTerminal(ctx context.Context, reference *runtime.ContainerReference, rio runtime.TerminalIO) error {
	cpr := &kubedef.ContainerPodReference{}
	if err := reference.Opaque.UnmarshalTo(cpr); err != nil {
		return fnerrors.InternalError("invalid reference: %w", err)
	}

	return r.attachTerminal(ctx, r.cli, cpr, rio)
}
