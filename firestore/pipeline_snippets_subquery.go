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

func defineExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_define_example]
	snapshot := client.Pipeline().
		Collection("authors").
		Define([]*firestore.AliasedExpression{
			firestore.FieldOf("id").As("currentAuthorId"),
		}).
		// [END firestore_define_example]
		AddFields(firestore.Selectables(
			client.Pipeline().
				Collection("books").
				Where(firestore.FieldOf("author_id").Equal(firestore.Variable("currentAuthorId"))).
				Aggregate(firestore.Accumulators(firestore.Average("rating").As("avgRating"))).
				ToScalarExpression().As("averageBookRating"),
		)).
		Execute(ctx)
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func toArrayExpression(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_to_array_expression]
	snapshot := client.Pipeline().
		Collection("projects").
		Define([]*firestore.AliasedExpression{
			firestore.FieldOf("id").As("parentId"),
		}).
		AddFields(firestore.Selectables(
			client.Pipeline().
				Collection("tasks").
				Where(firestore.FieldOf("project_id").Equal(firestore.Variable("parentId"))).
				Select(firestore.Fields("title")).
				ToArrayExpression().As("taskTitles"),
		)).
		Execute(ctx)
	// [END firestore_to_array_expression]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func toScalarExpression(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_to_scalar_expression]
	snapshot := client.Pipeline().
		Collection("authors").
		Define([]*firestore.AliasedExpression{
			firestore.FieldOf("id").As("currentAuthorId"),
		}).
		AddFields(firestore.Selectables(
			client.Pipeline().
				Collection("books").
				Where(firestore.FieldOf("author_id").Equal(firestore.Variable("currentAuthorId"))).
				Aggregate(firestore.Accumulators(firestore.Average("rating").As("avgRating"))).
				ToScalarExpression().As("averageBookRating"),
		)).
		Execute(ctx)
	// [END firestore_to_scalar_expression]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}
