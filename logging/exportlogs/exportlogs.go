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

// Sample exportlogs lists, creates, updates, and deletes log sinks.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/api/iterator"

	"cloud.google.com/go/logging/logadmin"
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

	// [START create_logging_client]
	ctx := context.Background()
	client, err := logadmin.NewClient(ctx, projID)
	if err != nil {
		log.Fatalf("logadmin.NewClient: %v", err)
	}
	// [END create_logging_client]

	switch command {
	case "list":
		log.Print("Listing log sinks.")
		sinks, err := listSinks(client)
		if err != nil {
			log.Fatalf("Could not list log sinks: %v", err)
		}
		for _, sink := range sinks {
			fmt.Printf("Sink: %v\n", sink)
		}
	case "create":
		if err := createSink(client); err != nil {
			log.Fatalf("Could not create sink: %v", err)
		}
	case "update":
		if err := updateSink(client); err != nil {
			log.Fatalf("Could not update sink: %v", err)
		}
	case "delete":
		if err := deleteSink(client); err != nil {
			log.Fatalf("Could not delete sink: %v", err)
		}
	default:
		usage("Unknown command.")
	}
}

func listSinks(client *logadmin.Client) ([]string, error) {
	// [START logging_list_sinks]
	ctx := context.Background()

	var sinks []string
	it := client.Sinks(ctx)
	for {
		sink, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		sinks = append(sinks, sink.ID)
	}
	// [END logging_list_sinks]
	return sinks, nil
}

func createSink(client *logadmin.Client) error {
	// [START logging_create_sink]
	ctx := context.Background()
	_, err := client.CreateSink(ctx, &logadmin.Sink{
		ID:          "severe-errors-to-gcs",
		Destination: "storage.googleapis.com/logsinks-bucket",
		Filter:      "severity >= ERROR",
	})
	// [END logging_create_sink]
	return err
}

func updateSink(client *logadmin.Client) error {
	// [START logging_update_sink]
	ctx := context.Background()
	_, err := client.UpdateSink(ctx, &logadmin.Sink{
		ID:          "severe-errors-to-gcs",
		Destination: "storage.googleapis.com/logsinks-new-bucket",
		Filter:      "severity >= INFO",
	})
	// [END logging_update_sink]
	return err
}

func deleteSink(client *logadmin.Client) error {
	// [START logging_delete_sink]
	ctx := context.Background()
	if err := client.DeleteSink(ctx, "severe-errors-to-gcs"); err != nil {
		return err
	}
	// [END logging_delete_sink]
	return nil
}

func usage(msg string) {
	if msg != "" {
		fmt.Fprintln(os.Stderr, msg)
	}
	fmt.Fprintln(os.Stderr, "usage: exportlogs <project-id> [list|create|update|delete]")
	os.Exit(2)
}
