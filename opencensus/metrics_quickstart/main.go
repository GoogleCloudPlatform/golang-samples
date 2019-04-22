// Copyright 2019 Google LLC
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

	// [START monitoring_opencensus_setup_exporter]
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

	if err := exporter.StartMetricsExporter(); err != nil {
		log.Fatalf("Error starting metric exporter: %v", err)
	}
	defer exporter.StopMetricsExporter()
	// [END monitoring_opencensus_setup_exporter]

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
