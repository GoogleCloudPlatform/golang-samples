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

// Sample simplelog writes some entries, lists them, then deletes the log.
package main

// [START logging_delete_log]
// [START logging_list_log_entries]
// [START logging_write_log_entry]
import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/logging"
	"cloud.google.com/go/logging/logadmin"
	"google.golang.org/api/iterator"
)

// [END logging_delete_log]
// [END logging_list_log_entries]
// [END logging_write_log_entry]

func main() {
	if len(os.Args) == 2 {
		usage("Missing command.")
	}
	if len(os.Args) != 3 {
		usage("")
	}

	projID := os.Args[1]
	command := os.Args[2]

	switch command {
	case "write":
		log.Print("Writing log entries.")
		structuredWrite(projID)

	case "read":
		log.Print("Fetching and printing log entries.")
		entries, err := getEntries(projID)
		if err != nil {
			log.Fatalf("Could not get entries: %v", err)
		}
		log.Printf("Found %d entries.", len(entries))
		for _, entry := range entries {
			fmt.Printf("Entry: %6s @%s: %v\n",
				entry.Severity,
				entry.Timestamp.Format(time.RFC3339),
				entry.Payload)
		}

	case "delete":
		log.Print("Deleting log.")
		if err := deleteLog(projID); err != nil {
			log.Fatalf("Could not delete log: %v", err)
		}

	default:
		usage("Unknown command.")
	}
}

// [START logging_write_log_entry]
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

// [START logging_delete_log]

func deleteLog(projectID string) error {
	ctx := context.Background()
	adminClient, err := logadmin.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create logadmin client: %v", err)
	}
	defer adminClient.Close()

	const name = "log-example"
	if err := adminClient.DeleteLog(ctx, name); err != nil {
		return err
	}
	return nil
}

// [END logging_delete_log]

// [START logging_list_log_entries]
func getEntries(projectID string) ([]*logging.Entry, error) {
	ctx := context.Background()
	adminClient, err := logadmin.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create logadmin client: %v", err)
	}
	defer adminClient.Close()

	var entries []*logging.Entry
	const name = "log-example"
	lastHour := time.Now().Add(-1 * time.Hour).Format(time.RFC3339)

	iter := adminClient.Entries(ctx,
		// Only get entries from the "log-example" log within the last hour.
		logadmin.Filter(fmt.Sprintf(`logName = "projects/%s/logs/%s" AND timestamp > "%s"`, projectID, name, lastHour)),
		// Get most recent entries first.
		logadmin.NewestFirst(),
	)

	// Fetch the most recent 20 entries.
	for len(entries) < 20 {
		entry, err := iter.Next()
		if err == iterator.Done {
			return entries, nil
		}
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

// [END logging_list_log_entries]

func usage(msg string) {
	if msg != "" {
		fmt.Fprintln(os.Stderr, msg)
	}
	fmt.Fprintln(os.Stderr, "usage: simplelog <project-id> [write|read|delete]")
	os.Exit(2)
}
