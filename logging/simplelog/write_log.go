// Copyright 2023 Google LLC
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

package main

// [START logging_write_log_entry]

import (
	"context"
	"log"

	"cloud.google.com/go/logging"
)

func structuredWrite(projectID string) {
	ctx := context.Background()
	client, err := logging.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create logging client: %v", err)
	}
	defer client.Close()
	const name = "log-example"
	logger := client.Logger(name)
	defer logger.Flush() // Ensure the entry is written.

	logger.Log(logging.Entry{
		// Log anything that can be marshaled to JSON.
		Payload: struct{ Anything string }{
			Anything: "The payload can be any type!",
		},
		Severity: logging.Debug,
	})
}

// [END logging_write_log_entry]
