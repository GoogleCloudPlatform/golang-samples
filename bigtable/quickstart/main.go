// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START bigtable_quickstart]

// Quickstart is a sample program demonstrating use of the Cloud Bigtable client
// library to read a row from an existing table.
package main

import (
	"flag"
	"log"

	"cloud.google.com/go/bigtable"
	"golang.org/x/net/context"
)

func main() {
	projectID := "my-project-id"   // The Google Cloud Platform project ID
	instanceID := "my-instance-id" // The Google Cloud Bigtable instance ID
	tableID := "my-table"          // The Google Cloud Bigtable table

	// [END bigtable_quickstart]
	// Override with -project, -instance, -table flags
	flag.StringVar(&projectID, "project", projectID, "The Google Cloud Platform project ID.")
	flag.StringVar(&instanceID, "instance", instanceID, "The Google Cloud Bigtable instance ID.")
	flag.StringVar(&tableID, "table", tableID, "The Google Cloud Bigtable table ID.")
	flag.Parse()

	// [START bigtable_quickstart]
	ctx := context.Background()

	// Set up Bigtable data operations client.
	client, err := bigtable.NewClient(ctx, projectID, instanceID)
	if err != nil {
		log.Fatalf("Could not create data operations client: %v", err)
	}

	tbl := client.Open(tableID)

	// Read data in a row using a row key
	rowKey := "r1"
	columnFamilyName := "cf1"

	log.Printf("Getting a single row by row key:")
	row, err := tbl.ReadRow(ctx, rowKey)
	if err != nil {
		log.Fatalf("Could not read row with key %s: %v", rowKey, err)
	}
	log.Printf("Row key: %s\n", rowKey)
	log.Printf("Data: %s\n", string(row[columnFamilyName][0].Value))

	if err = client.Close(); err != nil {
		log.Fatalf("Could not close data operations client: %v", err)
	}
}

// [END bigtable_quickstart]
