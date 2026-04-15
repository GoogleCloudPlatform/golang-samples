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

func wherePipeline(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_pipeline_where]
	results1, err := client.Pipeline().
		Collection("books").
		Where(firestore.FieldOf("rating").Equal(5)).
		Where(firestore.FieldOf("published").LessThan(1900)).
		Execute(ctx).Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "GetAll failed: %v", err)
		return err
	}

	results2, err := client.Pipeline().
		Collection("books").
		Where(firestore.And(
			firestore.FieldOf("rating").Equal(5),
			firestore.FieldOf("published").LessThan(1900),
		)).
		Execute(ctx).Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "GetAll failed: %v", err)
		return err
	}
	// [END firestore_pipeline_where]
	fmt.Fprintln(w, results1)
	fmt.Fprintln(w, results2)
	return nil
}

func whereEqualityExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_where_equality_example]
	snapshot := client.Pipeline().Collection("cities").
		Where(firestore.FieldOf("state").Equal("CA")).
		Execute(ctx)
	// [END firestore_where_equality_example]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func whereMultipleStagesExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_where_multiple_stages]
	snapshot := client.Pipeline().Collection("cities").
		Where(firestore.FieldOf("location.country").Equal("USA")).
		Where(firestore.FieldOf("population").GreaterThan(500000)).
		Execute(ctx)
	// [END firestore_where_multiple_stages]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func whereComplexExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_where_complex]
	snapshot := client.Pipeline().Collection("cities").
		Where(
			firestore.Or(
				firestore.Like(firestore.FieldOf("name"), "San%"),
				firestore.And(
					firestore.FieldOf("location.state").CharLength().GreaterThan(7),
					firestore.FieldOf("location.country").Equal("USA"),
				),
			),
		).
		Execute(ctx)
	// [END firestore_where_complex]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func whereStageOrderExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_where_stage_order]
	snapshot := client.Pipeline().Collection("cities").
		Limit(10).
		Where(firestore.FieldOf("location.country").Equal("USA")).
		Execute(ctx)
	// [END firestore_where_stage_order]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func createWhereData(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_create_where_data]
	client.Collection("cities").Doc("SF").Set(ctx, map[string]any{
		"name":       "San Francisco",
		"state":      "CA",
		"country":    "USA",
		"population": 870000,
	})
	client.Collection("cities").Doc("LA").Set(ctx, map[string]any{
		"name":       "Los Angeles",
		"state":      "CA",
		"country":    "USA",
		"population": 3970000,
	})
	client.Collection("cities").Doc("NY").Set(ctx, map[string]any{
		"name":       "New York",
		"state":      "NY",
		"country":    "USA",
		"population": 8530000,
	})
	client.Collection("cities").Doc("TOR").Set(ctx, map[string]any{
		"name":       "Toronto",
		"state":      nil,
		"country":    "Canada",
		"population": 2930000,
	})
	client.Collection("cities").Doc("MEX").Set(ctx, map[string]any{
		"name":       "Mexico City",
		"state":      nil,
		"country":    "Mexico",
		"population": 9200000,
	})
	// [END firestore_create_where_data]
	return nil
}
