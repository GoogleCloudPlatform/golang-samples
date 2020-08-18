// Copyright 2020 Google LLC
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

package main

// [START fs_listen_multiple]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// listenMultiple listens to a query, returning the names of all cities
// for a state.
func listenMultiple(w io.Writer, projectID string) error {
	// projectID := "project-id"
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("firestore.NewClient: %v", err)
	}
	defer client.Close()

	it := client.Collection("cities").Where("state", "==", "CA").Snapshots(ctx)
	for {
		snap, err := it.Next()
		if err != nil {
			return fmt.Errorf("Snapshots: listen failed: %v", err)
		}
		if snap.Size == 0 {
			return fmt.Errorf("current data: null")
		}
		for {
			doc, err := snap.Documents.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return fmt.Errorf("listen failed: %v", err)
			}
			fmt.Fprintf(w, "Current cities in California: %v\n", doc.Ref.ID)
		}
	}
}

// [END fs_listen_multiple]
