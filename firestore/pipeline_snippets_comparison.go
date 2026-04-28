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

func equalFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_equal_function]
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.Equal(firestore.FieldOf("rating"), 5).As("hasPerfectRating"),
		)).
		Execute(ctx)
	// [END firestore_equal_function]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func greaterThanFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_greater_than]
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.GreaterThan(firestore.FieldOf("rating"), 4).As("hasHighRating"),
		)).
		Execute(ctx)
	// [END firestore_greater_than]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func greaterThanOrEqualToFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_greater_or_equal]
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.GreaterThanOrEqual(firestore.FieldOf("published"), 1900).As("publishedIn20thCentury"),
		)).
		Execute(ctx)
	// [END firestore_greater_or_equal]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func lessThanFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_less_than]
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.LessThan(firestore.FieldOf("published"), 1923).As("isPublicDomainProbably"),
		)).
		Execute(ctx)
	// [END firestore_less_than]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func lessThanOrEqualToFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_less_or_equal]
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.LessThanOrEqual(firestore.FieldOf("rating"), 2).As("hasBadRating"),
		)).
		Execute(ctx)
	// [END firestore_less_or_equal]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func notEqualFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_not_equal]
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.NotEqual(firestore.FieldOf("title"), "1984").As("not1984"),
		)).
		Execute(ctx)
	// [END firestore_not_equal]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}
