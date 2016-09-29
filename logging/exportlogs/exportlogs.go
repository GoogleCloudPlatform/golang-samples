// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Sample exportlogs lists, creates, updates, and deletes log sinks.
package main

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/net/context"
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
	// [START list_log_sinks]
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
	// [END list_log_sinks]
	return sinks, nil
}

func createSink(client *logadmin.Client) error {
	// [START create_log_sink]
	ctx := context.Background()
	_, err := client.CreateSink(ctx, &logadmin.Sink{
		ID:          "severe-errors-to-gcs",
		Destination: "storage.googleapis.com/logsinks-bucket",
		Filter:      "severity >= ERROR",
	})
	// [END create_log_sink]
	return err
}

func updateSink(client *logadmin.Client) error {
	// [START update_log_sink]
	ctx := context.Background()
	_, err := client.UpdateSink(ctx, &logadmin.Sink{
		ID:          "severe-errors-to-gcs",
		Destination: "storage.googleapis.com/logsinks-new-bucket",
		Filter:      "severity >= INFO",
	})
	// [END update_log_sink]
	return err
}

func deleteSink(client *logadmin.Client) error {
	// [START delete_log_sink]
	ctx := context.Background()
	if err := client.DeleteSink(ctx, "severe-errors-to-gcs"); err != nil {
		return err
	}
	// [END delete_log_sink]
	return nil
}

func usage(msg string) {
	if msg != "" {
		fmt.Fprintln(os.Stderr, msg)
	}
	fmt.Fprintln(os.Stderr, "usage: exportlogs <project-id> [list|create|update|delete]")
	os.Exit(2)
}
