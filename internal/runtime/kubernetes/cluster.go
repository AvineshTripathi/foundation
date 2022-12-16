// Copyright 2022 Namespace Labs Inc; All rights reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.

package kubernetes

import (
	"context"
	"sync"

	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"namespacelabs.dev/foundation/framework/kubernetes/kubedef"
	"namespacelabs.dev/foundation/internal/runtime"
	"namespacelabs.dev/foundation/internal/runtime/kubernetes/client"
	"namespacelabs.dev/foundation/internal/runtime/kubernetes/networking/ingress"
	"namespacelabs.dev/foundation/internal/tcache"
	"namespacelabs.dev/foundation/std/cfg"
	"namespacelabs.dev/foundation/std/tasks"
)

type UnwrapCluster interface {
	KubernetesCluster() *Cluster
}

type Cluster struct {
	cli            *k8s.Clientset
	computedClient client.Prepared
	Configuration  cfg.Configuration

	FetchSystemInfo func(context.Context) (*kubedef.SystemInfo, error)

	ClusterAttachedState
}

type ClusterAttachedState struct {
	mu            sync.Mutex
	attachedState map[string]*state
}

type state struct {
	mu       sync.Mutex
	resolved bool
	value    any
	err      error
}

var _ kubedef.KubeCluster = &Cluster{}

func ConnectToCluster(ctx context.Context, config cfg.Configuration) (*Cluster, error) {
	cli, err := client.NewClient(ctx, config)
	if err != nil {
		return nil, err
	}

	deferredSystemInfo := tcache.NewDeferred(tasks.Action("kubernetes.fetch-system-info"), func(ctx context.Context) (*kubedef.SystemInfo, error) {
		return computeSystemInfo(ctx, cli.Clientset)
	})

	return NewCluster(cli, config, deferredSystemInfo.Get), nil
}

func NewCluster(cli *client.Prepared, config cfg.Configuration, fetchSystem FetchSystemInfoFunc) *Cluster {
	return &Cluster{
		cli:             cli.Clientset,
		computedClient:  *cli,
		Configuration:   config,
		FetchSystemInfo: fetchSystem,
	}
}

func (u *Cluster) Class() runtime.Class {
	return kubernetesClass{}
}

func (u *Cluster) Ingress() kubedef.IngressClass {
	return ingress.FromConfig(u.Configuration)
}

func (u *Cluster) RESTConfig() *rest.Config {
	return u.computedClient.RESTConfig
}

func (u *Cluster) KubernetesCluster() *Cluster { return u }

func (u *Cluster) PreparedClient() client.Prepared {
	return u.computedClient
}

func (u *Cluster) Bind(ctx context.Context, env cfg.Context) (runtime.ClusterNamespace, error) {
	return NewClusterNamespace(env, u, u), nil
}

func (r *Cluster) EnsureState(ctx context.Context, key string) (any, error) {
	return r.ClusterAttachedState.EnsureState(ctx, key, r.Configuration, r, nil)
}

func NewClusterNamespace(env cfg.Context, parent runtime.Cluster, u *Cluster) *ClusterNamespace {
	return &ClusterNamespace{parent: parent, underlying: u, target: newTarget(env)}
}

func (r *ClusterAttachedState) EnsureState(ctx context.Context, stateKey string, config cfg.Configuration, cluster runtime.Cluster, key *string) (any, error) {
	r.mu.Lock()
	if r.attachedState == nil {
		r.attachedState = map[string]*state{}
	}

	computedKey := stateKey
	if key != nil {
		computedKey += ":" + *key
	}

	if r.attachedState[computedKey] == nil {
		r.attachedState[computedKey] = &state{}
	}
	state := r.attachedState[computedKey]
	r.mu.Unlock()

	state.mu.Lock()
	defer state.mu.Unlock()

	if !state.resolved {
		if key != nil {
			state.value, state.err = runtime.PrepareKeyed(ctx, stateKey, config, cluster, *key)
		} else {
			state.value, state.err = runtime.Prepare(ctx, stateKey, config, cluster)
		}
		state.resolved = true
	}

	return state.value, state.err
}
