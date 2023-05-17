// Copyright 2023 Google LLC
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

// [START firestore_data_delete_collection]

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

func deleteCollection(w io.Writer, projectID, collectionName string,
	batchSize int) error {

	// Instantiate a client
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return err
	}

	col := client.Collection(collectionName)
	bulkwriter := client.BulkWriter(ctx)

	for {
		// Get a batch of documents
		iter := col.Limit(batchSize).Documents(ctx)
		numDeleted := 0

		// Iterate through the documents, adding
		// a delete operation for each one to the BulkWriter.
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return err
			}

			bulkwriter.Delete(doc.Ref)
			numDeleted++
		}

		// If there are no documents to delete,
		// the process is over.
		if numDeleted == 0 {
			bulkwriter.End()
			break
		}

		bulkwriter.Flush()
	}
	fmt.Fprintf(w, "Deleted collection \"%s\"", collectionName)
	return nil
}

// [END firestore_data_delete_collection]
