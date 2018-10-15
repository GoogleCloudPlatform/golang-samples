// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START monitoring_opencensus_metrics_quickstart]

// metrics_quickstart is an example of exporting a custom metric from
// OpenCensus to Stackdriver.
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"golang.org/x/exp/rand"
)

var (
	// The task latency in milliseconds.
	latencyMs = stats.Float64("task_latency", "The task latency in milliseconds", "ms")
)

func main() {
	ctx := context.Background()

	// Register the view. It is imperative that this step exists,
	// otherwise recorded metrics will be dropped and never exported.
	v := &view.View{
		Name:        "task_latency_distribution",
		Measure:     latencyMs,
		Description: "The distribution of the task latencies",

		// Latency in buckets:
		// [>=0ms, >=100ms, >=200ms, >=400ms, >=1s, >=2s, >=4s]
		Aggregation: view.Distribution(0, 100, 200, 400, 1000, 2000, 4000),
	}
	if err := view.Register(v); err != nil {
		log.Fatalf("Failed to register the view: %v", err)
	}

	// Enable OpenCensus exporters to export metrics
	// to Stackdriver Monitoring.
	// Exporters use Application Default Credentials to authenticate.
	// See https://developers.google.com/identity/protocols/application-default-credentials
	// for more details.
	exporter, err := stackdriver.NewExporter(stackdriver.Options{})
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

// [END monitoring_opencensus_metrics_quickstart]
