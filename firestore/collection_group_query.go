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

package firestore

// [START firestore_query_collection_group_filter_eq]
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
		return fmt.Errorf("firestore.NewClient: %w", err)
	}
	defer client.Close()

	it := client.CollectionGroup("landmarks").Where("type", "==", "museum").Documents(ctx)
	for {
		doc, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("documents iterator: %w", err)
		}
		fmt.Fprintf(w, "%s: %s", doc.Ref.ID, doc.Data()["name"])
	}

	return nil
}

// [END firestore_query_collection_group_filter_eq]

func partitionQuery(ctx context.Context, client *firestore.Client) error {
	// [START firestore_partition_query]
	cities := client.CollectionGroup("cities")

	// Get a partioned query for the cities collection group, with a maximum
	// partition count of 10
	partitionedQueries, err := cities.GetPartitionedQueries(ctx, 10)
	if err != nil {
		return fmt.Errorf("GetPartitionedQueries: %w", err)
	}

	fmt.Printf("Collection Group query partitioned to %d queries\n", len(partitionedQueries))

	// Retrieve the first query and iterate over the documents contained.
	query := partitionedQueries[0]
	iter := query.Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("documents iterator: %w", err)
		}
		fmt.Println(doc.Data())
	}

	// [END firestore_partition_query]
	return nil
}

func serializePartitionQuery(ctx context.Context, client *firestore.Client) error {
	// [START firestore_partition_query_serialization]
	cities := client.CollectionGroup("cities")

	// Get a partioned query for the cities collection group, with a maximum
	// partition count of 10
	partitionedQueries, err := cities.GetPartitionedQueries(ctx, 10)
	if err != nil {
		return err
	}

	fmt.Printf("Collection Group query partitioned to %d queries\n", len(partitionedQueries))

	query := partitionedQueries[0]

	// Serialize a query created by GetPartitionedQueries
	bytes, err := query.Serialize()
	if err != nil {
		return fmt.Errorf("Serialize: %w", err)
	}

	// Deserialize a query created by Query.Serialize
	deserializedQuery, err := client.CollectionGroup("").Deserialize(bytes)
	if err != nil {
		return fmt.Errorf("Deserialize: %w", err)
	}

	// [END firestore_partition_query_serialization]
	_ = deserializedQuery
	return nil
}
