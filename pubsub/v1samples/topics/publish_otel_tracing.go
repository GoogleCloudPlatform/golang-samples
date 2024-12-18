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

package topics

// [START pubsub_old_version_publish_otel_tracing]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/pubsub"
	"go.opentelemetry.io/otel"
	"google.golang.org/api/option"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// publishOpenTelemetryTracing publishes a single message with OpenTelemetry tracing
// enabled, exporting to Cloud Trace.
func publishOpenTelemetryTracing(w io.Writer, projectID, topicID string, sampling float64) error {
	// projectID := "my-project-id"
	// topicID := "my-topic"
	ctx := context.Background()

	exporter, err := texporter.New(texporter.WithProjectID(projectID),
		// Disable spans created by the exporter.
		texporter.WithTraceClientOptions(
			[]option.ClientOption{option.WithTelemetryDisabled()},
		),
	)
	if err != nil {
		return fmt.Errorf("error instantiating exporter: %w", err)
	}

	resources := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("publisher"),
	)

	// Instantiate a tracer provider with the following settings
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resources),
		sdktrace.WithSampler(
			sdktrace.ParentBased(sdktrace.TraceIDRatioBased(sampling)),
		),
	)

	defer tp.ForceFlush(ctx) // flushes any pending spans
	otel.SetTracerProvider(tp)

	// Create a new client with tracing enabled.
	client, err := pubsub.NewClientWithConfig(ctx, projectID, &pubsub.ClientConfig{
		EnableOpenTelemetryTracing: true,
	})
	if err != nil {
		return fmt.Errorf("pubsub: NewClient: %w", err)
	}
	defer client.Close()

	t := client.Topic(topicID)
	result := t.Publish(ctx, &pubsub.Message{
		Data: []byte("Publishing message with tracing"),
	})
	if _, err := result.Get(ctx); err != nil {
		return fmt.Errorf("pubsub: result.Get: %w", err)
	}
	fmt.Fprintln(w, "Published a traced message")
	return nil
}

// [END pubsub_old_version_publish_otel_tracing]
