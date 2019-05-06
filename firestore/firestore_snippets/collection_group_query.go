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

package main

// [START fs_collection_group_query]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// collectionGroupQuery runs a collection group query over the data created by
// collectionGroupSetup.
func collectionGroupQuery(w io.Writer, projectID string) error {
	ctx := context.Background()

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("firestore.NewClient: %v", err)
	}
	defer client.Close()

	it := client.CollectionGroup("landmarks").Where("type", "==", "museum").Documents(ctx)
	for {
		doc, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "%s: %s", doc.Ref.ID, doc.Data()["name"])
	}

	return nil
}

// [END fs_collection_group_query]
