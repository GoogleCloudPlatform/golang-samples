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

package main

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/spanner"
	traceExporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"google.golang.org/api/iterator"
)

// Ensure that your Go runtime version is supported by the OpenTelemetry-Go compatibility policy before enabling OpenTelemetry instrumentation.
// Refer to compatibility here https://github.com/googleapis/google-cloud-go/blob/main/debug.md#opentelemetry

// TODO(developer): Replace below variables before running the sample.
var projectID = "projectID"
var instanceID = "instanceID"
var databaseID = "databaseID"
var useCloudTraceExporter = true

// Sample to export traces to cloudtrace(default) or OTLP in the Go spanner client library
func main() {
	db := fmt.Sprintf("projects/%s/instances/%s/databases/%s", projectID, instanceID, databaseID)
	ctx := context.Background()
	var tracerProvider *sdktrace.TracerProvider

	if useCloudTraceExporter {
		tracerProvider = enableTracingWithCloudTraceExporter()
	} else {
		tracerProvider = enableTracingWithOtlpExporter()
	}
	defer tracerProvider.ForceFlush(ctx)

	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Execute a single read or query against Cloud Spanner.
	iter := client.Single().Read(ctx, "Singers", spanner.AllKeys(), []string{"SingerId", "FirstName", "LastName"})
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		var singerID int64
		var firstName, lastName string
		if err := row.Columns(&singerID, &firstName, &lastName); err != nil {
			log.Fatal(err)
		}
		log.Printf("%d %s %s\n", singerID, firstName, lastName)
	}
}

func enableTracingWithOtlpExporter() *sdktrace.TracerProvider {
	// [START spanner_opentelemetry_traces_otlp_usage]

	// Ensure that your Go runtime version is supported by the OpenTelemetry-Go
	// compatibility policy before enabling OpenTelemetry instrumentation.

	// Enable OpenTelemetry traces by setting environment variable GOOGLE_API_GO_EXPERIMENTAL_TELEMETRY_PLATFORM_TRACING
	// to the case-insensitive value "opentelemetry" before loading the client library.

	ctx := context.Background()

	// Create a new resource to uniquely identify the application
	res, err := resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName("My App"),
			semconv.ServiceVersion("0.1.0"),
		))
	if err != nil {
		log.Fatal(err)
	}

	// Create a new OTLP exporter.
	defaultOtlpEndpoint := "http://localhost:4317" // Replace with the endpoint on which OTLP collector is running
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithEndpoint(defaultOtlpEndpoint))
	if err != nil {
		log.Fatal(err)
	}

	// Create a new tracer provider
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(0.1)),
	)

	// Register tracer provider as global
	otel.SetTracerProvider(tracerProvider)
	// [END spanner_opentelemetry_traces_otlp_usage]

	return tracerProvider
}

func enableTracingWithCloudTraceExporter() *sdktrace.TracerProvider {
	// [START spanner_opentelemetry_traces_cloudtrace_usage]

	// Ensure that your Go runtime version is supported by the OpenTelemetry-Go
	// compatibility policy before enabling OpenTelemetry instrumentation.

	// Enable OpenTelemetry traces by setting environment variable GOOGLE_API_GO_EXPERIMENTAL_TELEMETRY_PLATFORM_TRACING
	// to the case-insensitive value "opentelemetry" before loading the client library.

	// Create a new resource to uniquely identify the application
	res, err := resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName("My App"),
			semconv.ServiceVersion("0.1.0"),
		))
	if err != nil {
		log.Fatal(err)
	}

	// Create a new cloud trace exporter
	exporter, err := traceExporter.New(traceExporter.WithProjectID(projectID))
	if err != nil {
		log.Fatal(err)
	}

	// Create a new tracer provider
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(0.1)),
	)

	// Register tracer provider as global
	otel.SetTracerProvider(tracerProvider)
	// [END spanner_opentelemetry_traces_cloudtrace_usage]

	return tracerProvider
}
