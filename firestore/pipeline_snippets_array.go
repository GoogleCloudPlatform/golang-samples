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
