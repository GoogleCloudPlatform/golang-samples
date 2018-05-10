// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// +build go1.8

// Sample opencensus_spanner_quickstart contains a sample application that
// uses Google Spanner Go client, and reports metrics
// and traces for the outgoing requests.
package main

import (
	"log"

	"cloud.google.com/go/spanner"
	"contrib.go.opencensus.io/exporter/stackdriver"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"golang.org/x/net/context"
)

func main() {
	ctx := context.Background()

	// Enable OpenCensus exporters to export traces and metrics
	// to Stackdriver Monitoring and Tracing.
	// Exporters use Application Default Credentials to authenticate.
	// See https://developers.google.com/identity/protocols/application-default-credentials
	// for more details.
	exporter, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID: "your-project-id",
	})
	if err != nil {
		log.Fatal(err)
	}
	view.RegisterExporter(exporter)
	trace.RegisterExporter(exporter)

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

	// Make sure data is uploaded before program finishes.
	exporter.Flush()
}
