// Copyright 2022 Namespace Labs Inc; All rights reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.

package deploy

import (
	"context"
	"fmt"

	"namespacelabs.dev/foundation/internal/compute"
	"namespacelabs.dev/foundation/internal/console"
	"namespacelabs.dev/foundation/internal/executor"
	"namespacelabs.dev/foundation/internal/fnerrors"
	"namespacelabs.dev/foundation/internal/planning"
	"namespacelabs.dev/foundation/internal/runtime"
	"namespacelabs.dev/foundation/schema"
	"namespacelabs.dev/foundation/std/cfg"
	"namespacelabs.dev/foundation/std/tasks"
)

type ComputeIngressResult struct {
	Fragments []*schema.IngressFragment

	rootenv cfg.Context
	stack   *schema.Stack
}

func ComputeIngress(planner planning.Planner, stack *schema.Stack, plans compute.Computable[[]*schema.IngressFragment], allocate bool) compute.Computable[*ComputeIngressResult] {
	return &computeIngress{rootenv: planner.Context, planner: planner.Runtime, stack: stack, fragments: plans, allocate: allocate}
}

func PlanIngressDeployment(rc runtime.Planner, c compute.Computable[*ComputeIngressResult]) compute.Computable[*runtime.DeploymentPlan] {
	return compute.Transform("plan ingress", c, func(ctx context.Context, res *ComputeIngressResult) (*runtime.DeploymentPlan, error) {
		return rc.PlanIngress(ctx, res.stack, res.Fragments)
	})
}

type computeIngress struct {
	rootenv   cfg.Context
	planner   runtime.Planner
	stack     *schema.Stack
	fragments compute.Computable[[]*schema.IngressFragment]
	allocate  bool // Actually fetch SSL certificates etc.

	compute.LocalScoped[*ComputeIngressResult]
}

func (ci *computeIngress) Action() *tasks.ActionEvent { return tasks.Action("deploy.compute-ingress") }
func (ci *computeIngress) Inputs() *compute.In {
	return compute.Inputs().
		Indigestible("cluster", ci.planner).
		Indigestible("rootenv", ci.rootenv).
		Proto("stack", ci.stack).
		Computable("fragments", ci.fragments)
}

func (ci *computeIngress) Output() compute.Output {
	return compute.Output{NotCacheable: true}
}

func (ci *computeIngress) Compute(ctx context.Context, deps compute.Resolved) (*ComputeIngressResult, error) {
	allFragments, err := computeDeferredIngresses(ctx, ci.rootenv, ci.planner, ci.stack)
	if err != nil {
		return nil, err
	}

	computed := compute.MustGetDepValue(deps, ci.fragments, "fragments")
	allFragments = append(allFragments, computed...)

	eg := executor.New(ctx, "compute.ingress")
	for _, fragment := range allFragments {
		fragment := fragment // Close fragment.

		eg.Go(func(ctx context.Context) error {
			sch := ci.stack.GetServer(schema.PackageName(fragment.Owner))
			if sch == nil {
				return fnerrors.BadInputError("%s: not present in the stack", fragment.Owner)
			}

			if ci.allocate {
				fragment.DomainCertificate, err = runtime.MaybeAllocateDomainCertificate(ctx, ci.rootenv.Environment(), sch, fragment.Domain)
				if err != nil {
					return err
				}
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return &ComputeIngressResult{
		rootenv:   ci.rootenv,
		stack:     ci.stack,
		Fragments: allFragments,
	}, nil
}

func computeDeferredIngresses(ctx context.Context, env cfg.Context, planner runtime.Planner, stack *schema.Stack) ([]*schema.IngressFragment, error) {
	var fragments []*schema.IngressFragment

	// XXX parallelism.
	for _, srv := range stack.Entry {
		frags, err := runtime.ComputeIngress(ctx, env, planner, srv, stack.Endpoint)
		if err != nil {
			return nil, err
		}
		fragments = append(fragments, frags...)
	}

	return fragments, nil
}

func computeIngressWithHandlerResult(planner planning.Planner, stack *planning.Stack, def compute.Computable[*handlerResult]) compute.Computable[*ComputeIngressResult] {
	computedIngressFragments := compute.Transform("parse computed ingress", def, func(ctx context.Context, h *handlerResult) ([]*schema.IngressFragment, error) {
		var fragments []*schema.IngressFragment

		for _, computed := range h.MergedComputedConfigurations().GetEntry() {
			for _, conf := range computed.Configuration {
				p := &schema.IngressFragment{}
				if !conf.Impl.MessageIs(p) {
					continue
				}

				if err := conf.Impl.UnmarshalTo(p); err != nil {
					return nil, err
				}

				fmt.Fprintf(console.Debug(ctx), "%s: received domain: %+v\n", conf.Owner, p.Domain)

				fragments = append(fragments, p)
			}
		}

		return fragments, nil
	})

	return ComputeIngress(planner, stack.Proto(), computedIngressFragments, AlsoDeployIngress)
}
