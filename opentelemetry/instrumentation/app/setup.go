// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"go.opentelemetry.io/contrib/exporters/autoexport"
	"go.opentelemetry.io/contrib/propagators/autoprop"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
)

func withTelemetry(ctx context.Context, fn func() error) error {
	// Configure Context Propagation to use the W3C traceparent
	otel.SetTextMapPropagator(autoprop.NewTextMapPropagator())

	// Configure Trace Export to send spans as OTLP
	texporter, err := autoexport.NewSpanExporter(ctx)
	if err != nil {
		return err
	}
	tp := trace.NewTracerProvider(trace.WithBatcher(texporter))
	defer tp.Shutdown(ctx)
	otel.SetTracerProvider(tp)

	// Configure Metric Export to send metrics as OTLP
	// TODO(https://github.com/open-telemetry/opentelemetry-go-contrib/issues/4131) - use env-based configuration instead of hardcoding
	mexporter, err := otlpmetrichttp.New(ctx)
	if err != nil {
		return err
	}
	mp := metric.NewMeterProvider(
		metric.WithReader(
			metric.NewPeriodicReader(
				mexporter,
				metric.WithInterval(5*time.Second),
			),
		),
	)
	defer mp.Shutdown(ctx)
	otel.SetMeterProvider(mp)

	// Configure Log Export to write JSON logs to stdout
	// Add trace attributes from context
	slog.SetDefault(slog.New(handerWithTraceContext(slog.NewJSONHandler(os.Stdout, nil))))

	return fn()
}
