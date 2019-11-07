// Copyright 2019 Google LLC. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Sample Manual Logging shows how to leverage Cloud Run structured logging without a client library.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/compute/metadata"
)

var projectID string

func main() {
	projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		projectID, _ = metadata.ProjectID()
	}
	if projectID == "" {
		log.Println("Could not determine Google Cloud Project. Running without log correlation. For local use set the $GOOGLE_CLOUD_PROJECT environment variable.")
	}

	http.HandleFunc("/", indexHandler)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		fmt.Println("Defaulting to port", port)
	}

	// Start HTTP server.
	fmt.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// [START run_manual_logging_object]

// Entry defines a log entry.
type Entry struct {
	Message  string `json:"message"`
	Severity string `json:"severity,omitempty"`
	Trace    string `json:"logging.googleapis.com/trace,omitempty"`

	// Log viewer accesses 'component' as 'jsonPayload.component'.
	Component string `json:"component,omitempty"`
}

// String renders an enty structure to the JSON format expected by Stackdriver.
func (e Entry) String() string {
	if e.Severity == "" {
		e.Severity = "INFO"
	}
	out, _ := json.Marshal(e)
	return string(out)
}

// [END run_manual_logging_object]
// [START run_manual_logging]

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Uncomment and populate this variable in your code:
	// projectID = "The project ID of your Cloud Run service"

	// Derive the traceID associated with the current request.
	var trace string
	if projectID != "" {
		traceHeader := r.Header.Get("X-Cloud-Trace-Context")
		traceParts := strings.Split(traceHeader, "/")
		if len(traceParts) > 0 && len(traceParts[0]) > 0 {
			trace = fmt.Sprintf("projects/%s/traces/%s", projectID, traceParts[0])
		}
	}

	log.Println(Entry{
		Severity:  "NOTICE",
		Message:   "This is the default display field.",
		Component: "arbitrary-property",
		Trace:     trace,
	})

	fmt.Fprintln(w, "Hello Logger!")
}

// [END run_manual_logging]
