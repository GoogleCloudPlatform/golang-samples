// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START logging_quickstart]
// Sample logging_quickstart writes a log entry to Stackdriver Logging.
package main

import (
	"fmt"
	"golang.org/x/net/context"
	"log"

	// Imports the Stackdriver Logging client package
	"cloud.google.com/go/logging"
)

func main() {
	ctx := context.Background()

	// Your Google Cloud Platform project ID
	projectID := "YOUR_PROJECT_ID"

	// Creates a client
	client, err := logging.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// The name of the log to write to
	logName := "my-log"

	// Selects the log to write to
	logger := client.Logger(logName)

	// The data to log
	text := "Hello, world!"

	// Adds an entry to the log buffer
	logger.Log(logging.Entry{Payload: text})

	// Closes the client and flushes the buffer to the Stackdriver Logging service
	err = client.Close()
	if err != nil {
		log.Fatalf("Failed to close client: %v", err)
	}

	fmt.Printf("Logged: %v", text)
}

// [END logging_quickstart]
