// Copyright 2026 Google LLC
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

package observability

// [START go_observability_tracing]
import (
	"context"
	"fmt"
	"log"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// enableTracing demonstrates how to configure OpenTelemetry tracing for
// Google Cloud Go client libraries. Set GOOGLE_SDK_GO_TRACING=true in the
// environment before starting the application to enable client span emission.
func enableTracing() error {
	ctx := context.Background()

	exporter, err := otlptracegrpc.New(ctx)
	if err != nil {
		return fmt.Errorf("otlptracegrpc.New: %w", err)
	}

	tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter))
	defer func() {
		if err := tp.Shutdown(context.WithoutCancel(ctx)); err != nil {
			log.Printf("tp.Shutdown: %v", err)
		}
	}()
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("secretmanager.NewClient: %w", err)
	}
	defer client.Close()

	// Use the client to make requests...

	return nil
}

// [END go_observability_tracing]
