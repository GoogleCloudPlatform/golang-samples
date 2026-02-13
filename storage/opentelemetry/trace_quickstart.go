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

// [START storage_enable_otel_tracing]
import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"

	"cloud.google.com/go/storage"
	traceExporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/contrib/detectors/gcp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

func run_quickstart(w io.Writer, projectID, bucket, object string) error {
	// projectID := "project-id"
	// bucket := "bucket-name"
	// object := "object-name"

	// Configure the sample rate to control trace ingestion volume.
	// Tracing sample rate must be in the range [0.0,1.0].
	sampleRate := 1.0
	ctx := context.Background()

	// Create a new cloud trace exporter.
	exporter, err := traceExporter.New(traceExporter.WithProjectID(projectID))
	if err != nil {
		return fmt.Errorf("Error creating new cloud trace exporter: %w", err)
	}

	// Identify your application using resource detection.
	res, err := resource.New(ctx,
		// Use the GCP resource detector to detect information about the GCP platform.
		resource.WithDetectors(gcp.NewDetector()),
		// Keep the default detectors.
		resource.WithTelemetrySDK(),
		// Add your own custom attributes to identify your application.
		resource.WithAttributes(
			semconv.ServiceName("My App"),
		),
	)
	if err != nil {
		return fmt.Errorf("Error creating new resource: %w", err)
	}

	// Create trace provider with the exporter.
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(sampleRate)),
	)
	defer tp.ForceFlush(ctx) // Flushes any pending spans.
	otel.SetTracerProvider(tp)

	tracer := otel.GetTracerProvider().Tracer("My App")
	ctx, span := tracer.Start(ctx, "trace-quickstart")

	// Instantiate a storage client and perform a write and read workload.
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	b := []byte("Hello world.")
	_, _ = client.Bucket(bucket).Attrs(ctx)
	wc := client.Bucket(bucket).Object(object).NewWriter(ctx)
	if _, err = io.Copy(wc, bytes.NewBuffer(b)); err != nil {
		return fmt.Errorf("io.Copy: %w", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %w", err)
	}
	fmt.Fprintf(w, "Uploaded blob %v to %v.\n", object, bucket)

	rc, err := client.Bucket(bucket).Object(object).NewReader(ctx)
	if err != nil {
		return fmt.Errorf("Object(%q).NewReader: %w", object, err)
	}
	defer rc.Close()

	_, err = ioutil.ReadAll(rc)
	if err != nil {
		return fmt.Errorf("ioutil.ReadAll: %w", err)
	}
	span.End()

	fmt.Fprintf(w, "Downloaded blob %v.\n", object)
	return nil
}

// [END storage_enable_otel_tracing]
