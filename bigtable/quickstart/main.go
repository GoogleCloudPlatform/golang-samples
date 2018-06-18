// Copyright 2018 Google LLC
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START bigtable_quickstart]
// Quickstart is a sample program demonstrating use of the Cloud Bigtable client
// library to read a row from an existing table.
package main

import (
	"flag"
	"log"
	"context"
	"cloud.google.com/go/bigtable"
)

// User-provided constants.
const (
	rowKey = "r1"
	columnFamilyName = "cf1"
)

func main() {
	projectID := flag.String("project", "", "The Google Cloud Platform project ID. Required.")
	instanceID := flag.String("instance", "", "The Google Cloud Bigtable instance ID. Required.")
	tableID := flag.String("table", "", "The Google Cloud Bigtable table ID. Required.")
	flag.Parse()

	for _, f := range []string{"project", "instance", "table"} {
		if flag.Lookup(f).Value.String() == "" {
			log.Fatalf("The %s flag is required.", f)
		}
	}

	ctx := context.Background()

	// Set up Bigtable data operations client.
	// NewClient uses Application Default Credentials to authenticate.
	client, err := bigtable.NewClient(ctx, *projectID, *instanceID)
	if err != nil {
		log.Fatalf("Could not create data operations client: %v", err)
	}

	tbl := client.Open(*tableID)

	// Read data in a row using a row key
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