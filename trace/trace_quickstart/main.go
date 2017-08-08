// Copyright 2017 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START trace_quickstart]

// Sample trace_quickstart creates traces incoming and outgoing requests.
package main

import (
	"log"
	"net/http"

	// Imports the Google Cloud Trace client package.
	"cloud.google.com/go/trace"
	"golang.org/x/net/context"
)

func main() {
	ctx := context.Background()

	// Sets your Google Cloud Platform project ID.
	projectID := "YOUR_PROJECT_ID"

	// Creates a client.
	traceClient, err := trace.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	httpClient := &http.Client{
		Transport: &trace.Transport{},
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, _ := http.NewRequest("GET", "https://metadata/users", nil)

		// The trace ID from the incoming request will be
		// propagated to the outgoing request.
		req = req.WithContext(r.Context())

		// The outgoing request will be traced with r's trace ID.
		if _, err := httpClient.Do(req); err != nil {
			log.Fatal(err)
		}
	})
	http.Handle("/foo", traceClient.HTTPHandler(handler))
	log.Fatal(http.ListenAndServe(":6060", nil))
}

// [END trace_quickstart]
