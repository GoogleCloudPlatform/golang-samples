// Copyright 2017 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START logging_stdlogging]

// Sample stdlogging writes log.Logger logs to the Stackdriver Logging.
package main

import (
	"log"

	// Imports the Stackdriver Logging client package.
	"cloud.google.com/go/logging"
	"golang.org/x/net/context"
)

// logger is a Go standard library logger.
var logger *log.Logger
var logClient *logging.Client

func init() {
	ctx := context.Background()

	// Sets your Google Cloud Platform project ID.
	projectID := "YOUR_PROJECT_ID"

	// Creates a client.
	client, err := logging.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create logging client: %v", err)
	}

	// Sets the name of the log to write to.
	logName := "my-log"

	logger = client.Logger(logName).StandardLogger(logging.Info)
	logClient = client
}

func main() {
	// Close flushes any pending messages and closes the client.
	defer logClient.Close()

	// Logs "hello world", log entry is visible at
	// Stackdriver Logs.
	logger.Println("hello world")
}

// [END logging_stdlogging]
