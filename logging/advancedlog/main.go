// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// [START logging_write_log_entry_advanced]

// Writes an advanced log entry to Cloud Logging.
package main

import (
	"context"
	"log"
	"net/http"

	"cloud.google.com/go/logging"
)

func main() {
	ctx := context.Background()

	// Creates a client.
	client, err := logging.NewClient(ctx, "YOUR_PROJECT_ID")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Sets the name of the log to write to.
	logger := client.Logger("my-log")

	// Logs a basic entry
	logger.Log(logging.Entry{Payload: "hello world"})

	// Logs an HTTPRequest type entry
	httpEntry := logging.Entry{
		Payload: "optional message",
		HTTPRequest: &logging.HTTPRequest{
			// TODO: pass in request
			Request: &http.Request{},
		},
	}
	logger.Log(httpEntry)
}

// [END logging_write_log_entry_advanced]
