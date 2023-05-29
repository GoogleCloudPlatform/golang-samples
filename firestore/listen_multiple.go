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

package firestore

// [START firestore_listen_query_snapshots]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// listenMultiple listens to a query, returning the names of all cities
// for a state.
func listenMultiple(ctx context.Context, w io.Writer, projectID, collection string) error {
	// projectID := "project-id"
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("firestore.NewClient: %w", err)
	}
	defer client.Close()

	it := client.Collection(collection).Where("state", "==", "CA").Snapshots(ctx)
	for {
		snap, err := it.Next()
		// DeadlineExceeded will be returned when ctx is cancelled.
		if status.Code(err) == codes.DeadlineExceeded {
			return nil
		}
		if err != nil {
			return fmt.Errorf("Snapshots.Next: %w", err)
		}
		if snap != nil {
			for {
				doc, err := snap.Documents.Next()
				if err == iterator.Done {
					break
				}
				if err != nil {
					return fmt.Errorf("Documents.Next: %w", err)
				}
				fmt.Fprintf(w, "Current cities in California: %v\n", doc.Ref.ID)
			}
		}
	}
}

// [END firestore_listen_query_snapshots]
