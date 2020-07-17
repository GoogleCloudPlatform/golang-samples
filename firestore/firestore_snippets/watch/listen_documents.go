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

package watch

// [START firestore_listen_document]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

func listenDocument(w io.Writer, projectID string) ([]*firestore.DocumentSnapshot, error) {
	// projectID := "project-id"
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("firestore.NewClient: %v", err)
	}
	defer client.Close()

	iter := client.Doc("cities/SF").Snapshots(ctx)
	defer iter.Stop()

	var docSnapshots []*firestore.DocumentSnapshot
	for {
		docSnaphot, err := iter.Next()
		if docSnaphot == nil {
			return nil, fmt.Errorf("current data: null")
		}
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Snapshots: listen failed: %v", err)
		}
		docSnapshots = append(docSnapshots, docSnaphot)
		fmt.Fprintf(w, "Received document snapshot: %v", docSnaphot.Ref.ID)
	}
	fmt.Fprintf(w, "Document snapshots were received")
	return docSnapshots, nil
}

// [END firestore_listen_document]
