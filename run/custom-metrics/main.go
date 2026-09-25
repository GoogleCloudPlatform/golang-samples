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

// [START cloudrun_mc_custom_metrics]
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

var counter metric.Int64Counter

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	shutdownMetrics := setupCounter(ctx)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	http.HandleFunc("/", handler)
	server := &http.Server{Addr: ":" + port}
	log.Printf("listening on port %s", port)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown failed: %v", err)
	}

	metricsCtx, metricsCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer metricsCancel()
	if err := shutdownMetrics(metricsCtx); err != nil {
		log.Printf("metrics shutdown failed: %v", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	counter.Add(context.Background(), 100)
	fmt.Fprintln(w, "Incremented sidecar_sample_counter_total metric!")
}

func setupCounter(ctx context.Context) func(context.Context) error {
	serviceName := os.Getenv("K_SERVICE")
	if serviceName == "" {
		serviceName = "sample-cloud-run-app"
	}
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			resource.Default().SchemaURL(),
			semconv.ServiceName(serviceName),
		),
	)
	if err != nil {
		log.Fatalf("Error creating resource: %v", err)
	}

	exporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("Error creating exporter: %s", err)
	}
	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)),
		sdkmetric.WithResource(r),
	)

	meter := provider.Meter("example.com/metrics")
	counter, err = meter.Int64Counter("sidecar-sample-counter")
	if err != nil {
		log.Fatalf("Error creating counter: %s", err)
	}
	return provider.Shutdown
}

// [END cloudrun_mc_custom_metrics]
