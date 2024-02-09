//go:build go1.20
// +build go1.20

// Copyright 2024 Google LLC
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

package spanner

// [START spanner_opentelemetry_usage]
// TODO: This works only from Go 1.20 due to OpenTelemetry incompatibility
import (
	"context"
	"fmt"
	"io"
	"log"

	"cloud.google.com/go/spanner"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"google.golang.org/api/iterator"
)

func enableOpenTelemetryMetricsAndTraces(w io.Writer, db string) error {
	// db = `projects/<project>/instances/<instance-id>/database/<database-id>`
	ctx := context.Background()

	// Create a new resource to uniquely identify the application
	res, err := newResource()
	if err != nil {
		log.Fatal(err)
	}

	// Enable OpenTelemetry traces by setting environment variable GOOGLE_API_GO_EXPERIMENTAL_TELEMETRY_PLATFORM_TRACING=opentelemetry at the start of the application
	// Enable OpenTelemetry metrics before injecting meter provider.
	spanner.EnableOpenTelemetryMetrics()

	// Create a new tracer provider
	tracerProvider, err := getOtlpTracerProvider(ctx, res)
	if err != nil {
		log.Fatal(err)
	}
	// Register tracer provider as global
	otel.SetTracerProvider(tracerProvider)

	// Create a new meter provider
	meterProvider := getOtlpMeterProvider(ctx, res)

	// Inject meter provider locally via ClientConfig when creating a spanner client.
	client, err := spanner.NewClientWithConfig(ctx, db, spanner.ClientConfig{OpenTelemetryMeterProvider: meterProvider})
	if err != nil {
		return err
	}
	defer client.Close()

	stmt := spanner.Statement{SQL: `SELECT SingerId, AlbumId, AlbumTitle FROM Albums`}
	iter := client.Single().Query(ctx, stmt)
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			return nil
		}
		if err != nil {
			return err
		}
		var singerID, albumID int64
		var albumTitle string
		if err := row.Columns(&singerID, &albumID, &albumTitle); err != nil {
			return err
		}
		fmt.Fprintf(w, "%d %d %s\n", singerID, albumID, albumTitle)
	}

	meterProvider.ForceFlush(ctx)
	tracerProvider.ForceFlush(ctx)

	return nil
}

func getOtlpMeterProvider(ctx context.Context, res *resource.Resource) *metric.MeterProvider {
	otlpExporter, err := otlpmetricgrpc.New(ctx)
	if err != nil {
		log.Fatal(err)
	}
	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(otlpExporter)),
	)
	return meterProvider
}

func getOtlpTracerProvider(ctx context.Context, res *resource.Resource) (*sdktrace.TracerProvider, error) {
	traceExporter, err := otlptracegrpc.New(ctx)
	if err != nil {
		return nil, err
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	return tracerProvider, nil
}

func newResource() (*resource.Resource, error) {
	return resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName("otlp-service"),
			semconv.ServiceVersion("0.1.0"),
		))
}

// [END spanner_opentelemetry_usage]
