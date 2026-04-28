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

func existsFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_exists_function]
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.FieldExists(firestore.FieldOf("rating")).As("hasRating"),
		)).
		Execute(ctx)
	// [END firestore_exists_function]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func andFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_and_function]
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.And(
				firestore.GreaterThan(firestore.FieldOf("rating"), 4),
				firestore.LessThan(firestore.FieldOf("price"), 10),
			).As("under10Recommendation"),
		)).
		Execute(ctx)
	// [END firestore_and_function]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func orFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_or_function]
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.Or(
				firestore.Equal(firestore.FieldOf("genre"), "Fantasy"),
				firestore.ArrayContains(firestore.FieldOf("tags"), "adventure"),
			).As("matchesSearchFilters"),
		)).
		Execute(ctx)
	// [END firestore_or_function]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func xorFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_xor_function]
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.Xor(
				firestore.ArrayContains(firestore.FieldOf("tags"), "magic"),
				firestore.ArrayContains(firestore.FieldOf("tags"), "nonfiction"),
			).As("matchesSearchFilters"),
		)).
		Execute(ctx)
	// [END firestore_xor_function]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func notFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_not_function]
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.Not(firestore.ArrayContains(firestore.FieldOf("tags"), "nonfiction")).As("isFiction"),
		)).
		Execute(ctx)
	// [END firestore_not_function]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func condFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_cond_function]
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.ArrayConcat(
				firestore.FieldOf("tags"),
				firestore.Conditional(
					firestore.GreaterThan(firestore.FieldOf("pages"), 100),
					firestore.ConstantOf("longRead"),
					firestore.ConstantOf("shortRead"),
				),
			).As("extendedTags"),
		)).
		Execute(ctx)
	// [END firestore_cond_function]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func maxLogicalFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_max_logical_function]
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.LogicalMaximum(firestore.FieldOf("rating"), 1).As("flooredRating"),
		)).
		Execute(ctx)
	// [END firestore_max_logical_function]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func minLogicalFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_min_logical_function]
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.LogicalMinimum(firestore.FieldOf("rating"), 5).As("cappedRating"),
		)).
		Execute(ctx)
	// [END firestore_min_logical_function]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}
