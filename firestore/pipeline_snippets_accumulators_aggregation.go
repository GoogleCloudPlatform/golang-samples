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

func countFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_count_function]
	// Total number of books in the collection
	countAll, err := client.Pipeline().Collection("books").
		Aggregate(firestore.Accumulators(firestore.CountAll().As("count"))).
		Execute(ctx).Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "GetAll failed: %v", err)
		return err
	}

	// Number of books with nonnull `ratings` field
	countField, err := client.Pipeline().
		Collection("books").
		Aggregate(firestore.Accumulators(firestore.Count("ratings").As("count"))).
		Execute(ctx).Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "GetAll failed: %v", err)
		return err
	}
	// [END firestore_count_function]
	fmt.Fprintln(w, countAll, countField)
	return nil
}

func countIfFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_count_if]
	snapshot := client.Pipeline().
		Collection("books").
		Aggregate(firestore.Accumulators(
			firestore.CountIf(firestore.FieldOf("rating").GreaterThan(4)).As("filteredCount"),
		)).
		Execute(ctx)
	// [END firestore_count_if]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func countDistinctFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_count_distinct]
	snapshot := client.Pipeline().
		Collection("books").
		Aggregate(firestore.Accumulators(
			firestore.CountDistinct("author").As("unique_authors"),
		)).
		Execute(ctx)
	// [END firestore_count_distinct]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func avgFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_avg_function]
	snapshot := client.Pipeline().
		Collection("cities").
		Aggregate(firestore.Accumulators(
			firestore.Average("population").As("averagePopulation"),
		)).
		Execute(ctx)
	// [END firestore_avg_function]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}
