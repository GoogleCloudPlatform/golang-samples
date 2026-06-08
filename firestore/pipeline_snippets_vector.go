// Copyright 2026 Google LLC
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

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/firestore"
)

func cosineDistanceFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_cosine_distance]
	sampleVector := []float64{0.0, 1.0, 2.0, 3.0, 4.0, 5.0}
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.CosineDistance(firestore.FieldOf("embedding"), sampleVector).As("cosineDistance"),
		)).
		Execute(ctx)
	// [END firestore_cosine_distance]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func dotProductFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_dot_product]
	sampleVector := []float64{0.0, 1.0, 2.0, 3.0, 4.0, 5.0}
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.DotProduct(firestore.FieldOf("embedding"), sampleVector).As("dotProduct"),
		)).
		Execute(ctx)
	// [END firestore_dot_product]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func euclideanDistanceFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_euclidean_distance]
	sampleVector := []float64{0.0, 1.0, 2.0, 3.0, 4.0, 5.0}
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.EuclideanDistance(firestore.FieldOf("embedding"), sampleVector).As("euclideanDistance"),
		)).
		Execute(ctx)
	// [END firestore_euclidean_distance]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func vectorLengthFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_vector_length]
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.VectorLength(firestore.FieldOf("embedding")).As("vectorLength"),
		)).
		Execute(ctx)
	// [END firestore_vector_length]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}
