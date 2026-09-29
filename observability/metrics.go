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

// [START go_observability_metrics]
import (
	"context"
	"fmt"
	"log"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

// enableMetrics demonstrates how to configure OpenTelemetry metrics for
// Google Cloud Go client libraries. Set GOOGLE_SDK_GO_METRICS=true in the
// environment before starting the application to enable metric collection.
func enableMetrics() error {
	ctx := context.Background()

	exporter, err := otlpmetricgrpc.New(ctx)
	if err != nil {
		return fmt.Errorf("otlpmetricgrpc.New: %w", err)
	}

	mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)))
	defer func() {
		if err := mp.Shutdown(context.WithoutCancel(ctx)); err != nil {
			log.Printf("mp.Shutdown: %v", err)
		}
	}()
	otel.SetMeterProvider(mp)

	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("secretmanager.NewClient: %w", err)
	}
	defer client.Close()

	// Use the client to make requests...

	return nil
}

// [END go_observability_metrics]
