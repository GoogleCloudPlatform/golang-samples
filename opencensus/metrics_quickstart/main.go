// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// metrics_quickstart is an example of exporting a custom metric from
// OpenCensus to Stackdriver.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"golang.org/x/exp/rand"
)

var (
	// The latency in milliseconds.
	latencyMs = stats.Float64("demo/latency", "The latency in milliseconds", "ms")

	latencyView = &view.View{
		Name:        "demo/latency",
		Measure:     latencyMs,
		Description: "The distribution of the latencies",

		// Latency in buckets:
		// [>=0ms, >=100ms, >=200ms, >=400ms, >=1s, >=2s, >=4s]
		Aggregation: view.Distribution(0, 100, 200, 400, 1000, 2000, 4000),
	}
)

var projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")

func main() {
	if projectID == "" {
		log.Fatal("The GOOGLE_CLOUD_PROJECT environment variable must be set.")
	}

	ctx := context.Background()

	// Register the view. It is imperative that this step exists,
	// otherwise recorded metrics will be dropped and never exported.
	if err := view.Register(latencyView); err != nil {
		log.Fatalf("Failed to register the view: %v", err)
	}

	// Enable OpenCensus exporters to export metrics
	// to Stackdriver Monitoring.
	// Exporters use Application Default Credentials to authenticate.
	// See https://developers.google.com/identity/protocols/application-default-credentials
	// for more details.
	exporter, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID:    projectID,
		MetricPrefix: "my-metrics",
	})
	if err != nil {
		log.Fatal(err)
	}
	// Flush must be called before main() exits to ensure metrics are recorded.
	defer exporter.Flush()

	view.RegisterExporter(exporter)

	// The minimum reporting period for Stackdriver is 1 minute.
	view.SetReportingPeriod(60 * time.Second)

	// Record 100 fake latency values between 0 and 5 seconds.
	for i := 0; i < 100; i++ {
		ms := float64(5*time.Second/time.Millisecond) * rand.Float64()
		fmt.Printf("Latency %d: %f\n", i, ms)
		stats.Record(ctx, latencyMs.M(ms))
		time.Sleep(1 * time.Second)
	}

	fmt.Println("Done recording metrics")
}
