// Copyright 2020 Google LLC
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
//
// The trace command is an example of setting up OpenTelemetry to export traces to Google
// Cloud Trace.
package main

// [START opentelemetry_trace_import]
import (
	"context"
	"log"
	"os"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// [END opentelemetry_trace_import]
// [START opentelemetry_trace_main_function]
func main() {
	// Create exporter.
	ctx := context.Background()
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	exporter, err := texporter.NewExporter(texporter.WithProjectID(projectID))
	if err != nil {
		log.Fatalf("texporter.NewExporter: %v", err)
	}
	defer exporter.Shutdown(ctx) // flushes any pending spans

	// Create trace provider with the exporter.
	//
	// By default it uses AlwaysSample() which samples all traces.
	// In a production environment or high QPS setup please use
	// probabilistic sampling.
	// Example:
	//   tp := sdktrace.NewTracerProvider(sdktrace.WithSampler(sdktrace.TraceIDRatioBased(0.0001)), ...)
	tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter))
	otel.SetTracerProvider(tp)

	// [START opentelemetry_trace_custom_span]
	// Create custom span.
	tracer := otel.GetTracerProvider().Tracer("example.com/trace")
	err = func(ctx context.Context) error {
		ctx, span := tracer.Start(ctx, "foo")
		defer span.End()

		// Do some work.

		return nil
	}(ctx)
	// [END opentelemetry_trace_custom_span]
}

// [END opentelemetry_trace_main_function]
