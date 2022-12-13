// Copyright 2022 Namespace Labs Inc; All rights reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.

package planning

import (
	"sort"
	"strings"
	"sync"

	"namespacelabs.dev/foundation/schema"
	"namespacelabs.dev/foundation/std/pkggraph"
)

type stackBuilder struct {
	mu      sync.Mutex
	servers map[schema.PackageName]*PlannedServer
}

func (stack *stackBuilder) claim(pkgname schema.PackageName) (*PlannedServer, bool) {
	stack.mu.Lock()
	defer stack.mu.Unlock()

	existing, has := stack.servers[pkgname]
	if has {
		return existing, true
	}

	if stack.servers == nil {
		stack.servers = map[schema.PackageName]*PlannedServer{}
	}

	ps := &PlannedServer{}
	stack.servers[pkgname] = ps
	return ps, false
}

func (stack *stackBuilder) buildStack(computedResources map[string][]pkggraph.ResourceInstance, focusPackages ...schema.PackageName) *Stack {
	stack.mu.Lock()
	defer stack.mu.Unlock()

	var focus schema.PackageList
	for _, pkg := range focusPackages {
		focus.Add(pkg)
	}

	s := &Stack{
		Focus:             focus,
		ComputedResources: computedResources,
	}

	for _, sb := range stack.servers {
		s.Servers = append(s.Servers, *sb)
	}

	sort.Slice(s.Servers, func(i, j int) bool {
		return order(focus, s.Servers[i].Server.PackageName(), s.Servers[j].Server.PackageName())
	})

	var endpoints []*schema.Endpoint
	var internal []*schema.InternalEndpoint

	for _, srv := range s.Servers {
		endpoints = append(endpoints, srv.Endpoints...)
		internal = append(internal, srv.InternalEndpoints...)
	}

	sort.Slice(endpoints, func(i, j int) bool {
		e_i := endpoints[i]
		e_j := endpoints[j]

		if e_i.ServerOwner == e_j.ServerOwner {
			return strings.Compare(e_i.AllocatedName, e_j.AllocatedName) < 0
		}
		return order(focus, schema.PackageName(e_i.ServerOwner), schema.PackageName(e_j.ServerOwner))
	})

	sort.Slice(internal, func(i, j int) bool {
		e_i := internal[i]
		e_j := internal[j]

		if e_i.ServerOwner == e_j.ServerOwner {
			return e_i.GetPort().GetContainerPort() < e_j.Port.GetContainerPort()
		}
		return order(focus, schema.PackageName(e_i.ServerOwner), schema.PackageName(e_j.ServerOwner))
	})

	s.Endpoints = endpoints
	s.InternalEndpoints = internal
	return s
}

func (stack *stackBuilder) changeServer(handler func()) {
	stack.mu.Lock()
	defer stack.mu.Unlock()
	handler()
}

func order(foci schema.PackageList, a, b schema.PackageName) bool {
	if foci.Has(a) {
		if !foci.Has(b) {
			return true
		}
	} else if foci.Has(b) {
		return false
	}

	return strings.Compare(a.String(), b.String()) < 0
}
