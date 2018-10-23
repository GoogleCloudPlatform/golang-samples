// Copyright 2017 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START trace_setup_go_quickstart]

// Sample trace_quickstart traces incoming and outgoing requests.
package main

import (
	"log"
	"net/http"
	"os"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"contrib.go.opencensus.io/exporter/stackdriver/propagation"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
)

func main() {
	// Create and register a OpenCensus Stackdriver Trace exporter.
	exporter, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID: os.Getenv("GOOGLE_CLOUD_PROJECT"),
	})
	if err != nil {
		log.Fatal(err)
	}
	trace.RegisterExporter(exporter)

	client := &http.Client{
		Transport: &ochttp.Transport{
			// Use Google Cloud propagation format.
			Propagation: &propagation.HTTPFormat{},
		},
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, _ := http.NewRequest("GET", "https://metadata/users", nil)

		// The trace ID from the incoming request will be
		// propagated to the outgoing request.
		req = req.WithContext(r.Context())

		// The outgoing request will be traced with r's trace ID.
		if _, err := client.Do(req); err != nil {
			log.Fatal(err)
		}
	})
	http.Handle("/foo", handler)
	log.Fatal(http.ListenAndServe(":6060", &ochttp.Handler{}))
}

// [END trace_setup_go_quickstart]
