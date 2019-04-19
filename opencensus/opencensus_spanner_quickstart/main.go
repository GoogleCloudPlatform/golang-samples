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

// +build go1.8

// Sample opencensus_spanner_quickstart contains a sample application that
// uses Google Spanner Go client, and reports metrics
// and traces for the outgoing requests.
package main

import (
	"context"
	"log"

	"cloud.google.com/go/spanner"
	"contrib.go.opencensus.io/exporter/stackdriver"
	"go.opencensus.io/trace"
)

func main() {
	ctx := context.Background()

	// Enable OpenCensus exporters to export traces and metrics
	// to Stackdriver Monitoring and Tracing.
	// Exporters use Application Default Credentials to authenticate.
	// See https://developers.google.com/identity/protocols/application-default-credentials
	// for more details.
	exporter, err := stackdriver.NewExporter(stackdriver.Options{})
	if err != nil {
		log.Fatal(err)
	}
	// Flush must be called before main() exits to ensure metrics are recorded.
	defer exporter.Flush()

	trace.RegisterExporter(exporter)

	if err := exporter.StartMetricsExporter(); err != nil {
		log.Fatalf("Error starting metric exporter: %v", err)
	}
	defer exporter.StopMetricsExporter()

	// Use trace.AlwaysSample() to always record traces. The
	// default sampler skips some traces to conserve resources,
	// but can make it hard to debug test traffic. So, remove
	// the following line before pushing to production.
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	// This database must exist.
	databaseName := "projects/your-project-id/instances/your-instance-id/databases/your-database-id"

	client, err := spanner.NewClient(ctx, databaseName)
	if err != nil {
		log.Fatalf("Failed to create client %v", err)
	}
	defer client.Close()

	_, err = client.Apply(ctx, []*spanner.Mutation{
		spanner.Insert("Users",
			[]string{"name", "email"},
			[]interface{}{"alice", "a@example.com"})})
	if err != nil {
		log.Printf("Failed to insert: %v", err)
	}
}
