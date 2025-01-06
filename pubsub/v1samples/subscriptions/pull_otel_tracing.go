// Copyright 2025 Google LLC
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

package subscriptions

// [START pubsub_old_version_subscribe_otel_tracing]
import (
	"context"
	"fmt"
	"io"
	"sync/atomic"
	"time"

	"cloud.google.com/go/pubsub"
	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"google.golang.org/api/option"
)

func subscribeOpenTelemetryTracing(w io.Writer, projectID, subID string, sampleRate float64) error {
	// projectID := "my-project-id"
	// subID := "my-sub"
	// sampleRate := "1.0"
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
		semconv.ServiceNameKey.String("subscriber"),
	)

	// Instantiate a tracer provider with the following settings
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resources),
		sdktrace.WithSampler(
			sdktrace.ParentBased(sdktrace.TraceIDRatioBased(sampleRate)),
		),
	)

	defer tp.ForceFlush(ctx) // flushes any pending spans
	otel.SetTracerProvider(tp)

	// Create a new client with tracing enabled.
	client, err := pubsub.NewClientWithConfig(ctx, projectID, &pubsub.ClientConfig{
		EnableOpenTelemetryTracing: true,
	})
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %w", err)
	}
	defer client.Close()

	sub := client.Subscription(subID)

	// Receive messages for 10 seconds, which simplifies testing.
	// Comment this out in production, since `Receive` should
	// be used as a long running operation.
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var received int32
	err = sub.Receive(ctx, func(_ context.Context, msg *pubsub.Message) {
		fmt.Fprintf(w, "Got message: %q\n", string(msg.Data))
		atomic.AddInt32(&received, 1)
		msg.Ack()
	})
	if err != nil {
		return fmt.Errorf("sub.Receive: %w", err)
	}
	fmt.Fprintf(w, "Received %d messages\n", received)

	return nil
}

// [END pubsub_old_version_subscribe_otel_tracing]
