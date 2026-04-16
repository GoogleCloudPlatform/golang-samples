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

func searchBasic(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_search_basic]
	snapshot := client.Pipeline().
		Collection("restaurants").
		Search(firestore.WithSearchQuery(firestore.DocumentMatches("waffles"))).
		Execute(ctx)
	// [END firestore_search_basic]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func searchExact(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_search_exact]
	snapshot := client.Pipeline().
		Collection("restaurants").
		Search(firestore.WithSearchQuery(firestore.DocumentMatches("\"belgian waffles\""))).
		Execute(ctx)
	// [END firestore_search_exact]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func searchTwoTerms(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_search_two_terms]
	snapshot := client.Pipeline().
		Collection("restaurants").
		Search(firestore.WithSearchQuery(firestore.DocumentMatches("waffles eggs"))).
		Execute(ctx)
	// [END firestore_search_two_terms]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func searchExcludeTerm(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_search_exclude_term]
	snapshot := client.Pipeline().
		Collection("restaurants").
		Search(firestore.WithSearchQuery(firestore.DocumentMatches("-waffles"))).
		Execute(ctx)
	// [END firestore_search_exclude_term]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func searchSpecialFields(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_search_special_fields]
	snapshot := client.Pipeline().
		Collection("restaurants").
		Search(
			firestore.WithSearchQuery(firestore.FieldOf("menu").RegexMatch("waffles")),
			firestore.WithSearchAddFields(firestore.Score().As("score")),
		).
		Execute(ctx)
	// [END firestore_search_special_fields]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}
