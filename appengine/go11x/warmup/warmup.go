// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START gae_go111_warmup]

// Sample warmup demonstrates usage of the /_ah/warmup handler.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/storage"
)

var startupTime time.Time
var client *storage.Client

func main() {
	// Perform required setup steps for the application to function.
	// This assumes any returned error requires a new instance to be created.
	if err := setup(context.Background()); err != nil {
		log.Fatalf("setup: %v", err)
	}

	// Log when an appengine warmup request is used to create the new instance.
	// Warmup steps are taken in setup for consistency with "cold start" instances.
	http.HandleFunc("/_ah/warmup", func(w http.ResponseWriter, r *http.Request) {
		log.Println("warmup done")
	})
	http.HandleFunc("/", indexHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

// setup executes per-instance one-time warmup and initialization actions.
func setup(ctx context.Context) error {
	// Store the startup time of the server.
	startupTime = time.Now()

	// Initialize a Google Cloud Storage client.
	var err error
	if client, err = storage.NewClient(ctx); err != nil {
		return err
	}

	return nil
}

// indexHandler responds to requests with our greeting.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	uptime := time.Since(startupTime).Seconds()
	fmt.Fprintf(w, "Hello, World! Uptime: %.2fs\n", uptime)
}

// [END gae_go111_warmup]
