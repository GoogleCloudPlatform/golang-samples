// Copyright 2024 Google LLC
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

// [START firestore_vector_search_distance_result_field]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/firestore"
)

func vectorSearchDistanceResultField(w io.Writer, projectID string) error {
	ctx := context.Background()

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("firestore.NewClient: %w", err)
	}
	defer client.Close()

	collection := client.Collection("coffee-beans")

	// Requires a vector index
	// https://firebase.google.com/docs/firestore/vector-search#create_and_manage_vector_indexes
	vectorQuery := collection.FindNearest("embedding_field",
		[]float32{3.0, 1.0, 2.0},
		10,
		firestore.DistanceMeasureEuclidean,
		&firestore.FindNearestOptions{
			DistanceResultField: "vector_distance",
		})

	docs, err := vectorQuery.Documents(ctx).GetAll()
	if err != nil {
		fmt.Fprintf(w, "failed to get vector query results: %v", err)
		return err
	}

	for _, doc := range docs {
		fmt.Fprintf(w, "%v, Distance: %v\n", doc.Data()["name"], doc.Data()["vector_distance"])
	}
	return nil
}

// [END firestore_vector_search_distance_result_field]
