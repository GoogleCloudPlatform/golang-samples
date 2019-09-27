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

// Sample firestore_quickstart demonstrates how to connect to Firestore, and add and list documents.
package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/api/iterator"

	"cloud.google.com/go/firestore"
)

func main() {

	// [START fs_initialize]
	// Sets your Google Cloud Platform project ID.
	projectID := "YOUR_PROJECT_ID"

	// Get a Firestore client.
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Close client when done.
	defer client.Close()
	// [END fs_initialize]

	// [START fs_add_data_1]
	_, _, err = client.Collection("users").Add(ctx, map[string]interface{}{
		"first": "Ada",
		"last":  "Lovelace",
		"born":  1815,
	})
	if err != nil {
		log.Fatalf("Failed adding alovelace: %v", err)
	}
	// [END fs_add_data_1]

	// [START fs_add_data_2]
	_, _, err = client.Collection("users").Add(ctx, map[string]interface{}{
		"first":  "Alan",
		"middle": "Mathison",
		"last":   "Turing",
		"born":   1912,
	})
	if err != nil {
		log.Fatalf("Failed adding aturing: %v", err)
	}
	// [END fs_add_data_2]

	// [START fs_get_all_users]
	iter := client.Collection("users").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		fmt.Println(doc.Data())
	}
	// [END fs_get_all_users]
}
