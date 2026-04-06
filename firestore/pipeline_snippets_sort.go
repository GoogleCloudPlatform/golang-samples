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

func sort(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START sort]
	snapshot := client.Pipeline().
		Collection("books").
		Sort(firestore.Orders(
			firestore.Descending(firestore.FieldOf("release_date")),
			firestore.Ascending(firestore.FieldOf("author")),
		)).
		Execute(ctx)
	// [END sort]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func sortComparison(w io.Writer, client *firestore.Client) error {
	// [START sort_comparison]
	query := client.Collection("cities").
		OrderBy("state", firestore.Asc).
		OrderBy("population", firestore.Desc)

	pipeline := client.Pipeline().
		Collection("books").
		Sort(firestore.Orders(
			firestore.Descending(firestore.FieldOf("release_date")),
			firestore.Ascending(firestore.FieldOf("author")),
		))
	// [END sort_comparison]
	fmt.Fprintln(w, query)
	fmt.Fprintln(w, pipeline)
	return nil
}

func sortSyntaxExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START sort_syntax]
	snapshot := client.Pipeline().Collection("cities").
		Sort(firestore.Orders(firestore.Ascending(firestore.FieldOf("population")))).
		Execute(ctx)
	// [END sort_syntax]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func sortSyntaxExample2(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START sort_syntax_2]
	snapshot := client.Pipeline().Collection("cities").
		Sort(firestore.Orders(firestore.Ascending(firestore.CharLength(firestore.FieldOf("name"))))).
		Execute(ctx)
	// [END sort_syntax_2]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func sortDocumentIDExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START sort_document_id]
	snapshot := client.Pipeline().Collection("cities").
		Sort(firestore.Orders(
			firestore.Ascending(firestore.FieldOf("country")),
			firestore.Ascending(firestore.FieldOf("__name__")),
		)).
		Execute(ctx)
	// [END sort_document_id]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}
