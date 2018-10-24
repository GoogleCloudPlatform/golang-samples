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
	if err := setup(context.Background()); err != nil {
		log.Fatalf("setup: %v", err)
	}

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

// setup executes per-instance one-time setup actions.
func setup(ctx context.Context) error {
	var err error

	// Store the startup time of the server.
	startupTime = time.Now()

	// Initialize a Google Cloud Storage client.
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
