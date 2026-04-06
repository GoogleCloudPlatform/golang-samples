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

func selectSyntaxExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_select_syntax]
	snapshot := client.Pipeline().Collection("cities").
		Select(firestore.Fields(
			firestore.StringConcat(firestore.FieldOf("name"), ", ", firestore.FieldOf("location.country")).As("name"),
			firestore.FieldOf("population"),
		)).
		Execute(ctx)
	// [END firestore_select_syntax]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func selectPositionDataExample(_ io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_select_position_data]
	client.Collection("cities").Doc("SF").Set(ctx, map[string]any{
		"name":       "San Francisco",
		"population": 800000,
		"location": map[string]any{
			"country": "USA",
			"state":   "California",
		},
	})
	client.Collection("cities").Doc("TO").Set(ctx, map[string]any{
		"name":       "Toronto",
		"population": 3000000,
		"location": map[string]any{
			"country":  "Canada",
			"province": "Ontario",
		},
	})
	// [END firestore_select_position_data]
	return nil
}

func selectPositionExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_select_position]
	snapshot := client.Pipeline().Collection("cities").
		Where(firestore.FieldOf("location.country").Equal("Canada")).
		Select(firestore.Fields(
			firestore.StringConcat(firestore.FieldOf("name"), ", ", firestore.FieldOf("location.country")).As("name"),
			firestore.FieldOf("population"),
		)).
		Execute(ctx)
	// [END firestore_select_position]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func selectBadPositionExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_select_bad_position]
	snapshot := client.Pipeline().Collection("cities").
		Select(firestore.Fields(
			firestore.StringConcat(firestore.FieldOf("name"), ", ", firestore.FieldOf("location.country")).As("name"),
			firestore.FieldOf("population"),
		)).
		Where(firestore.FieldOf("location.country").Equal("Canada")).
		Execute(ctx)
	// [END firestore_select_bad_position]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func selectNestedDataExample(_ io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_select_nested_data]
	client.Collection("cities").Doc("SF").Set(ctx, map[string]any{
		"name":       "San Francisco",
		"population": 800000,
		"location": map[string]any{
			"country": "USA",
			"state":   "California",
		},
		"landmarks": []string{"Golden Gate Bridge", "Alcatraz"},
	})
	client.Collection("cities").Doc("TO").Set(ctx, map[string]any{
		"name":       "Toronto",
		"population": 3000000,
		"province":   "ON",
		"location": map[string]any{
			"country":  "Canada",
			"province": "Ontario",
		},
		"landmarks": []string{"CN Tower", "Casa Loma"},
	})
	client.Collection("cities").Doc("AT").Set(ctx, map[string]any{
		"name":       "Atlantis",
		"population": nil,
	})
	// [END firestore_select_nested_data]
	return nil
}

func selectNestedExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_select_nested]
	snapshot := client.Pipeline().Collection("cities").
		Select(firestore.Fields(
			firestore.FieldOf("name").As("city"),
			firestore.FieldOf("location.country").As("country"),
			firestore.ArrayGet(firestore.FieldOf("landmarks"), 0).As("topLandmark"),
		)).
		Execute(ctx)
	// [END firestore_select_nested]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func addFieldsSyntaxExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_add_fields_syntax]
	snapshot := client.Pipeline().Collection("users").
		AddFields(firestore.Selectables(
			firestore.StringConcat(firestore.FieldOf("firstName"), " ", firestore.FieldOf("lastName")).As("fullName"),
		)).
		Execute(ctx)
	// [END firestore_add_fields_syntax]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func addFieldsOverlapExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_add_fields_overlap]
	snapshot := client.Pipeline().Collection("users").
		AddFields(firestore.Selectables(firestore.Abs(firestore.FieldOf("age")).As("age"))).
		AddFields(firestore.Selectables(firestore.Add(firestore.FieldOf("age"), 10).As("age"))).
		Execute(ctx)
	// [END firestore_add_fields_overlap]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func addFieldsNestingExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_add_fields_nesting]
	snapshot := client.Pipeline().Collection("users").
		AddFields(firestore.Selectables(
			firestore.ToLower(firestore.FieldOf("address.city")).As("address.city"),
		)).
		Execute(ctx)
	// [END firestore_add_fields_nesting]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func removeFieldsSyntaxExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_remove_fields_syntax]
	snapshot := client.Pipeline().Collection("cities").
		RemoveFields(firestore.Fields("population", "location.state")).
		Execute(ctx)
	// [END firestore_remove_fields_syntax]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func removeFieldsNestedDataExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_remove_fields_nested_data]
	client.Collection("cities").Doc("SF").Set(ctx, map[string]any{
		"name": "San Francisco",
		"location": map[string]any{
			"country": "USA",
			"state":   "California",
		},
	})
	client.Collection("cities").Doc("TO").Set(ctx, map[string]any{
		"name": "Toronto",
		"location": map[string]any{
			"country":  "Canada",
			"province": "Ontario",
		},
	})
	// [END firestore_remove_fields_nested_data]
	return nil
}

func removeFieldsNestedExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_remove_fields_nested]
	snapshot := client.Pipeline().Collection("cities").
		RemoveFields(firestore.Fields("location.state")).
		Execute(ctx)
	// [END firestore_remove_fields_nested]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}
