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

func sumFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_sum_function]
	snapshot := client.Pipeline().
		Collection("cities").
		Aggregate(firestore.Accumulators(
			firestore.Sum("population").As("totalPopulation"),
		)).
		Execute(ctx)
	// [END firestore_sum_function]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func minFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_min_function]
	snapshot := client.Pipeline().
		Collection("books").
		Aggregate(firestore.Accumulators(
			firestore.Minimum("price").As("minimumPrice"),
		)).
		Execute(ctx)
	// [END firestore_min_function]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func maxFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_max_function]
	snapshot := client.Pipeline().
		Collection("books").
		Aggregate(firestore.Accumulators(
			firestore.Maximum("price").As("maximumPrice"),
		)).
		Execute(ctx)
	// [END firestore_max_function]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func arrayConcatFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_array_concat]
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.ArrayConcat(firestore.FieldOf("genre"), firestore.FieldOf("subGenre")).As("allGenres"),
		)).
		Execute(ctx)
	// [END firestore_array_concat]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func arrayContainsFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_array_contains]
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.ArrayContains(firestore.FieldOf("genre"), "mystery").As("isMystery"),
		)).
		Execute(ctx)
	// [END firestore_array_contains]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func arrayContainsAllFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_array_contains_all]
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.ArrayContainsAll(firestore.FieldOf("genre"), []string{"fantasy", "adventure"}).As("isFantasyAdventure"),
		)).
		Execute(ctx)
	// [END firestore_array_contains_all]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func arrayContainsAnyFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_array_contains_any]
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.ArrayContainsAny(firestore.FieldOf("genre"), []string{"fantasy", "nonfiction"}).As("isMysteryOrFantasy"),
		)).
		Execute(ctx)
	// [END firestore_array_contains_any]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func arrayLengthFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_array_length]
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.ArrayLength(firestore.FieldOf("genre")).As("genreCount"),
		)).
		Execute(ctx)
	// [END firestore_array_length]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func arrayReverseFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_array_reverse]
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.ArrayReverse(firestore.FieldOf("genre")).As("reversedGenres"),
		)).
		Execute(ctx)
	// [END firestore_array_reverse]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func equalAnyFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_eq_any]
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.EqualAny(firestore.FieldOf("genre"), []string{"Science Fiction", "Psychological Thriller"}).As("matchesGenreFilters"),
		)).
		Execute(ctx)
	// [END firestore_eq_any]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func notEqualAnyFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_not_eq_any]
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.NotEqualAny(firestore.FieldOf("author"), []string{"George Orwell", "F. Scott Fitzgerald"}).As("byExcludedAuthors"),
		)).
		Execute(ctx)
	// [END firestore_not_eq_any]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func functionsExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_functions_example]
	// Type 1: Scalar (for use in non-aggregation stages)
	// Example: Return the min store price for each book.
	results1, err := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.LogicalMinimum(firestore.FieldOf("current"), firestore.FieldOf("updated")).As("price_min"),
		)).
		Execute(ctx).Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "GetAll failed: %v", err)
		return err
	}

	// Type 2: Aggregation (for use in aggregate stages)
	// Example: Return the min price of all books.
	results2, err := client.Pipeline().
		Collection("books").
		Aggregate(firestore.Accumulators(
			firestore.Minimum("price").As("min_price"),
		)).
		Execute(ctx).Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "GetAll failed: %v", err)
		return err
	}
	// [END firestore_functions_example]
	fmt.Fprintln(w, results1)
	fmt.Fprintln(w, results2)
	return nil
}
