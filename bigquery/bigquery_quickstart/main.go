// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START bigquery_quickstart]
// Sample bigquery_quickstart creates a Google BigQuery dataset.
package main

import (
	"fmt"
	"golang.org/x/net/context"
	"log"

	// Imports the Google Cloud Datastore client package
	"cloud.google.com/go/bigquery"
)

func main() {
	ctx := context.Background()

	// Your Google Cloud Platform project ID
	projectID := "YOUR_PROJECT_ID"

	// Creates a client
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// The name for the new dataset
	datasetName := "my_new_dataset"

	// Prepares the new dataset
	dataset := client.Dataset(datasetName)

	// Creates the dataset
	if err := dataset.Create(ctx); err != nil {
		log.Fatalf("Failed to create dataset: %v", err)
	}

	fmt.Printf("Dataset %v created", dataset)
}

// [END bigquery_quickstart]
