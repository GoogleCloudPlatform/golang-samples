// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Sample force_trace_http demonstrates tracing 100% of requests using Cloud Trace.
package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/cloud/trace"
)

var traceClient *trace.Client

func main() {
	ctx := context.Background()

	var err error
	traceClient, err = trace.NewClient(ctx, mustGetenv("GCLOUD_PROJECT"))
	if err != nil {
		log.Fatalf("Could not start Cloud Trace: %v", err)
	}

	http.HandleFunc("/", slowHandler)
	// Force tracing of every request.
	http.ListenAndServe(":8080", http.HandlerFunc(forceTraceHandler))
}

func forceTraceHandler(w http.ResponseWriter, r *http.Request) {
	if h := "X-Cloud-Trace-Context"; r.Header.Get(h) == "" {
		// Generate a trace header.
		// https://cloud.google.com/trace/docs/faq
		r.Header.Set(h, generateTraceID()+"/0;o=1")
	}
	http.DefaultServeMux.ServeHTTP(w, r)
}

func generateTraceID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func slowHandler(w http.ResponseWriter, r *http.Request) {
	span := traceClient.SpanFromRequest(r)
	defer span.Finish()

	time.Sleep(50 * time.Millisecond)
	{
		span := span.NewChild("more_slowness")
		time.Sleep(100 * time.Millisecond)
		span.Finish()
	}
	time.Sleep(50 * time.Millisecond)

	fmt.Fprintf(w, "done")
}

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("%s environment variable not set.", k)
	}
	return v
}
