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

import (
	"fmt"
	"log"
	"os"
	"time"

	// [START imports]
	"context"

	"cloud.google.com/go/logging"
	"cloud.google.com/go/logging/logadmin"

	"google.golang.org/api/iterator"
	// [END imports]
)

func main() {
	if len(os.Args) == 2 {
		usage("Missing command.")
	}
	if len(os.Args) != 3 {
		usage("")
	}

	projID := os.Args[1]
	command := os.Args[2]

	// [START setup]
	ctx := context.Background()
	client, err := logging.NewClient(ctx, projID)
	if err != nil {
		log.Fatalf("Failed to create logging client: %v", err)
	}

	adminClient, err := logadmin.NewClient(ctx, projID)
	if err != nil {
		log.Fatalf("Failed to create logadmin client: %v", err)
	}

	client.OnError = func(err error) {
		// Print an error to the local log.
		// For example, if Flush() failed.
		log.Printf("client.OnError: %v", err)
	}
	// [END setup]

	switch command {
	case "write":
		log.Print("Writing some log entries.")
		writeEntry(client)
		structuredWrite(client)

	case "read":
		log.Print("Fetching and printing log entries.")
		entries, err := getEntries(adminClient, projID)
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
		if err := deleteLog(adminClient); err != nil {
			log.Fatalf("Could not delete log: %v", err)
		}

	default:
		usage("Unknown command.")
	}
}

func writeEntry(client *logging.Client) {
	// [START logging_write_log_entry]
	const name = "log-example"
	logger := client.Logger(name)
	defer logger.Flush() // Ensure the entry is written.

	infolog := logger.StandardLogger(logging.Info)
	infolog.Printf("infolog is a standard Go log.Logger with INFO severity.")
	// [END logging_write_log_entry]
}

func structuredWrite(client *logging.Client) {
	// [START write_structured_log_entry]
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
	// [END logging_write_log_entry]
}

func deleteLog(adminClient *logadmin.Client) error {
	ctx := context.Background()

	// [START logging_delete_log]
	const name = "log-example"
	if err := adminClient.DeleteLog(ctx, name); err != nil {
		return err
	}
	// [END logging_delete_log]
	return nil
}

func getEntries(adminClient *logadmin.Client, projID string) ([]*logging.Entry, error) {
	ctx := context.Background()

	// [START logging_list_log_entries]
	var entries []*logging.Entry
	const name = "log-example"
	iter := adminClient.Entries(ctx,
		// Only get entries from the log-example log.
		logadmin.Filter(fmt.Sprintf(`logName = "projects/%s/logs/%s"`, projID, name)),
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
	// [END logging_list_log_entries]
}

func usage(msg string) {
	if msg != "" {
		fmt.Fprintln(os.Stderr, msg)
	}
	fmt.Fprintln(os.Stderr, "usage: simplelog <project-id> [write|read|delete]")
	os.Exit(2)
}
