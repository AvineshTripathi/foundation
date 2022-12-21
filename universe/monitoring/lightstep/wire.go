// Copyright 2022 Namespace Labs Inc; All rights reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.

package lightstep

import (
	"context"
	"os"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
)

func Prepare(ctx context.Context, deps ExtensionDeps) error {
	accessToken := os.Getenv("MONITORING_LIGHTSTEP_ACCESS_TOKEN")
	if accessToken == "" {
		// No secret specified.
		return nil
	}

	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint("ingest.lightstep.com:443"),
		otlptracegrpc.WithHeaders(map[string]string{
			"lightstep-access-token": accessToken,
		}),
	}

	client := otlptracegrpc.NewClient(opts...)
	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		return err
	}

	return deps.OpenTelemetry.Register(exporter)
}
