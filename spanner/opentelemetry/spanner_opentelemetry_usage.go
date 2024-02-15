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

package opentelemetry

// [START spanner_opentelemetry_usage]
// Ensure that your Go runtime version is supported by the OpenTelemetry-Go compatibility policy before enabling OpenTelemetry instrumentation.
// Refer to compatibility here https://github.com/googleapis/google-cloud-go/blob/main/debug.md#opentelemetry

import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"cloud.google.com/go/spanner"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
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

	// Enable OpenTelemetry traces by setting environment variable GOOGLE_API_GO_EXPERIMENTAL_TELEMETRY_PLATFORM_TRACING to the case-insensitive value "opentelemetry" before loading the client library.
	// Enable OpenTelemetry metrics before injecting meter provider.
	spanner.EnableOpenTelemetryMetrics()

	// Create a new tracer provider
	tracerProvider, err := getOtlpTracerProvider(ctx, res)
	defer tracerProvider.ForceFlush(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// Register tracer provider as global
	otel.SetTracerProvider(tracerProvider)

	// Create a new meter provider
	meterProvider := getOtlpMeterProvider(ctx, res)
	defer meterProvider.ForceFlush(ctx)

	// Inject meter provider locally via ClientConfig when creating a spanner client or set globally via setMeterProvider.
	client, err := spanner.NewClientWithConfig(ctx, db, spanner.ClientConfig{OpenTelemetryMeterProvider: meterProvider})
	if err != nil {
		return err
	}
	defer client.Close()
	return nil
}

func getOtlpMeterProvider(ctx context.Context, res *resource.Resource) *sdkmetric.MeterProvider {
	otlpExporter, err := otlpmetricgrpc.New(ctx)
	if err != nil {
		log.Fatal(err)
	}
	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(otlpExporter)),
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

// [START spanner_opentelemetry_gfe_metric]
// GFE_Latency and other Spanner metrics are automatically collected
// when OpenTelemetry metrics are enabled.
func captureGFELatencyMetric(ctx context.Context, client spanner.Client) error {
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
	}
}

// [END spanner_opentelemetry_gfe_metric]

// [START spanner_opentelemetry_capture_query_stats_metric]
func captureQueryStatsMetric(ctx context.Context, mp metric.MeterProvider, client spanner.Client) error {
	meter := mp.Meter(spanner.OtInstrumentationScope)
	// Register query stats metric with OpenTelemetry to record the data.
	// This should be done once before start recording the data.
	queryStats, err := meter.Float64Histogram(
		"spanner/query_stats_elapsed",
		metric.WithDescription("The execution of the query"),
		metric.WithUnit("ms"),
		metric.WithExplicitBucketBoundaries(0.0, 0.01, 0.05, 0.1, 0.3, 0.6, 0.8, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 8.0, 10.0, 13.0,
			16.0, 20.0, 25.0, 30.0, 40.0, 50.0, 65.0, 80.0, 100.0, 130.0, 160.0, 200.0, 250.0,
			300.0, 400.0, 500.0, 650.0, 800.0, 1000.0, 2000.0, 5000.0, 10000.0, 20000.0, 50000.0,
			100000.0),
	)
	if err != nil {
		fmt.Print(err)
	}

	stmt := spanner.Statement{SQL: `SELECT SingerId, AlbumId, AlbumTitle FROM Albums`}
	iter := client.Single().QueryWithStats(ctx, stmt)
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			// Record query execution time with OpenTelemetry.
			elapasedTime := iter.QueryStats["elapsed_time"].(string)
			elapasedTimeMs, err := strconv.ParseFloat(strings.TrimSuffix(elapasedTime, " msecs"), 64)
			if err != nil {
				return err
			}
			queryStats.Record(ctx, elapasedTimeMs)
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
	}
}

// [END spanner_opentelemetry_capture_query_stats_metric]
