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
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// listenChanges listens to a query, returning the list of document changes.
func listenChanges(ctx context.Context, w io.Writer, projectID string) ([]firestore.DocumentChange, error) {
	// projectID := "project-id"
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("firestore.NewClient: %v", err)
	}
	defer client.Close()

	modified := false
	changes := make([]firestore.DocumentChange, 0)
	it := client.Collection("cities").Where("state", "==", "CA").Snapshots(ctx)
	for {
		snap, err := it.Next()
		if err != nil {
			// DeadlineExceeded will be returned when ctx is cancelled.
			if status.Code(err) == codes.DeadlineExceeded && modified {
				return changes, nil
			} else if status.Code(err) == codes.DeadlineExceeded && !modified {
				return changes, fmt.Errorf("expected document wasn't modified")
			}
			return nil, fmt.Errorf("Snapshots.Next: %v", err)
		}
		if snap != nil {
			for _, change := range snap.Changes {
				changes = append(changes, change)
				switch change.Kind {
				case firestore.DocumentAdded:
					fmt.Fprintf(w, "New city: %v\n", change.Doc.Data())
				case firestore.DocumentModified:
					modified = true
					fmt.Fprintf(w, "Modified city: %v\n", change.Doc.Data())
				case firestore.DocumentRemoved:
					fmt.Fprintf(w, "Removed city: %v\n", change.Doc.Data())
				}
			}
		}
	}
}

// [END fs_listen_changes]
