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
	defer client.Close()

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
