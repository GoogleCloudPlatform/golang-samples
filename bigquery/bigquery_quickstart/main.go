// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START bigquery_quickstart]
// Sample bigquery-quickstart creates a Google BigQuery dataset.
package main

import (
	"fmt"
	"log"

	// Imports the Google Cloud BigQuery client package.
	"cloud.google.com/go/bigquery"
	"golang.org/x/net/context"
)

func main() {
	ctx := context.Background()

	// Sets your Google Cloud Platform project ID.
	projectID := "YOUR_PROJECT_ID"

	// Creates a client.
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Sets the name for the new dataset.
	datasetName := "my_new_dataset"

	// Creates a Dataset instance.
	dataset := client.Dataset(datasetName)

	// Creates the new BigQuery dataset.
	if err := dataset.Create(ctx); err != nil {
		log.Fatalf("Failed to create dataset: %v", err)
	}

	fmt.Printf("Dataset created\n")
}

// [END bigquery_quickstart]
