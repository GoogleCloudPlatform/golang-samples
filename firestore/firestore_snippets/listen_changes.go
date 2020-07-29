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

// [START fs_listen_changes]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/firestore"
)

// listenChanges listens to a query, returning the list of DocumentChange
// in the first snapshot.
func listenChanges(w io.Writer, projectID string) error {
	// projectID := "project-id"
	ctx := context.Background()

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("firestore.NewClient: %v", err)
	}
	defer client.Close()

	snap, err := client.Collection("cities").Where("state", "==", "CA").Snapshots(ctx).Next()
	if err != nil {
		return fmt.Errorf("Snapshots: listen failed: %v", err)
	}
	if snap.Size == 0 {
		return fmt.Errorf("current data: null")
	}
	for _, change := range snap.Changes {
		switch change.Kind {
		case firestore.DocumentAdded:
			fmt.Fprintf(w, "New city: %v\n", change.Doc.Data())
		case firestore.DocumentModified:
			fmt.Fprintf(w, "Modified city: %v\n", change.Doc.Data())
		case firestore.DocumentRemoved:
			fmt.Fprintf(w, "Removed city: %v\n", change.Doc.Data())
		}
	}
	return nil
}

// [END fs_listen_changes]
