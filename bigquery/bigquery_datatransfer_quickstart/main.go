// Copyright 2017 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START bigquerydatatransfer_quickstart]

// Sample bigquery-quickstart creates a Google BigQuery dataset.
package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/api/iterator"

	// Imports the BigQuery Data Transfer client package.
	datatransfer "cloud.google.com/go/bigquery/datatransfer/apiv1"
	datatransferpb "google.golang.org/genproto/googleapis/cloud/bigquery/datatransfer/v1"
)

func main() {
	ctx := context.Background()

	// Sets your Google Cloud Platform project ID.
	projectID := "YOUR_PROJECT_ID"

	// Creates a client.
	client, err := datatransfer.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	req := &datatransferpb.ListDataSourcesRequest{
		Parent: fmt.Sprintf("projects/%s", projectID),
	}
	it := client.ListDataSources(ctx, req)
	fmt.Println("Supported Data Sources:")
	for {
		ds, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to list sources: %v", err)
		}
		fmt.Println(ds.DisplayName)
		fmt.Println("\tID: ", ds.DataSourceId)
		fmt.Println("\tFull path: ", ds.Name)
		fmt.Println("\tDescription: ", ds.Description)
	}
}

// [END bigquerydatatransfer_quickstart]
