// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START gae_go111_warmup]

// Sample warmup demonstrates usage of the /_ah/warmup handler.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/_ah/warmup", warmupHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func warmupHandler(w http.ResponseWriter, r *http.Request) {
	// Perform warmup tasks, including ones that require a context,
	// such as retrieving data from Datastore.

	log.Println("warmup done")
}

// [END gae_go111_warmup]