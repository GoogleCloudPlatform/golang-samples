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

// [START firestore_listen_handle_error]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// listenErrors demonstrates how to handle listening errors.
func listenErrors(ctx context.Context, w io.Writer, projectID, collection string) error {
	// projectID := "project-id"
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("firestore.NewClient: %w", err)
	}
	defer client.Close()

	it := client.Collection(collection).Snapshots(ctx)
	for {
		snap, err := it.Next()
		// Canceled will be returned when ctx is cancelled and DeadlineExceeded will
		// be returned when ctx reaches its deadline.
		if e := status.Code(err); e == codes.Canceled || e == codes.DeadlineExceeded {
			return nil
		}
		if err != nil {
			return fmt.Errorf("Snapshots.Next: %w", err)
		}
		if snap != nil {
			for _, change := range snap.Changes {
				if change.Kind == firestore.DocumentAdded {
					fmt.Fprintf(w, "New city: %v\n", change.Doc.Data())
				}
			}
		}
	}
}

// [END firestore_listen_handle_error]
