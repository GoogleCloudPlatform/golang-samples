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

func mapGetFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START map_get]
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.MapGet(firestore.FieldOf("awards"), "pulitzer").As("hasPulitzerAward"),
		)).
		Execute(ctx)
	// [END map_get]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func mapSetFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START map_set]
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.MapSet(firestore.FieldOf("awards"), "pulitzer", true).As("awards"),
		)).
		Execute(ctx)
	// [END map_set]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func mapKeysFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START map_keys]
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.MapKeys(firestore.FieldOf("awards")).As("award_categories"),
		)).
		Execute(ctx)
	// [END map_keys]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func mapValuesFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START map_values]
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.MapValues(firestore.FieldOf("awards")).As("award_details"),
		)).
		Execute(ctx)
	// [END map_values]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func mapEntriesFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START map_entries]
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.MapEntries(firestore.FieldOf("awards")).As("awards_list"),
		)).
		Execute(ctx)
	// [END map_entries]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}
