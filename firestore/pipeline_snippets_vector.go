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
	// [START cosine_distance]
	sampleVector := []float64{0.0, 1.0, 2.0, 3.0, 4.0, 5.0}
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.CosineDistance(firestore.FieldOf("embedding"), sampleVector).As("cosineDistance"),
		)).
		Execute(ctx)
	// [END cosine_distance]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func dotProductFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START dot_product]
	sampleVector := []float64{0.0, 1.0, 2.0, 3.0, 4.0, 5.0}
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.DotProduct(firestore.FieldOf("embedding"), sampleVector).As("dotProduct"),
		)).
		Execute(ctx)
	// [END dot_product]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func euclideanDistanceFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START euclidean_distance]
	sampleVector := []float64{0.0, 1.0, 2.0, 3.0, 4.0, 5.0}
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.EuclideanDistance(firestore.FieldOf("embedding"), sampleVector).As("euclideanDistance"),
		)).
		Execute(ctx)
	// [END euclidean_distance]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func vectorLengthFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START vector_length]
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.VectorLength(firestore.FieldOf("embedding")).As("vectorLength"),
		)).
		Execute(ctx)
	// [END vector_length]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func findNearestSyntaxExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START find_nearest_syntax]
	snapshot := client.Pipeline().Collection("cities").
		FindNearest("embedding", []float64{1.5, 2.345}, firestore.PipelineDistanceMeasureEuclidean).
		Execute(ctx)
	// [END find_nearest_syntax]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func findNearestLimitExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START find_nearest_limit]
	snapshot := client.Pipeline().Collection("cities").
		FindNearest("embedding", []float64{1.5, 2.345}, firestore.PipelineDistanceMeasureEuclidean, firestore.WithFindNearestLimit(10)).
		Execute(ctx)
	// [END find_nearest_limit]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func findNearestDistanceDataExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START find_nearest_distance_data]
	client.Collection("cities").Doc("SF").Set(ctx, map[string]any{
		"name":      "San Francisco",
		"embedding": []float64{1.0, -1.0},
	})
	client.Collection("cities").Doc("TO").Set(ctx, map[string]any{
		"name":      "Toronto",
		"embedding": []float64{5.0, -10.0},
	})
	client.Collection("cities").Doc("AT").Set(ctx, map[string]any{
		"name":      "Atlantis",
		"embedding": []float64{2.0, -4.0},
	})
	// [END find_nearest_distance_data]
	return nil
}

func findNearestDistanceExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START find_nearest_distance]
	snapshot := client.Pipeline().Collection("cities").
		FindNearest("embedding", []float64{1.3, 2.345}, firestore.PipelineDistanceMeasureEuclidean, firestore.WithFindNearestDistanceField("computedDistance")).
		Execute(ctx)
	// [END find_nearest_distance]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}
