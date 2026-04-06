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

func unnestStage(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START unnest_stage]
	snapshot := client.Pipeline().
		Database().
		UnnestWithAlias("arrayField", "unnestedArrayField", firestore.WithUnnestIndexField("index")).
		Execute(ctx)
	// [END unnest_stage]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func unnestStageEmptyOrNonArray(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START unnest_edge_cases]
	// Input
	// { "identifier" : 1, "neighbors": [ "Alice", "Cathy" ] }
	// { "identifier" : 2, "neighbors": []                   }
	// { "identifier" : 3, "neighbors": "Bob"                }

	results, err := client.Pipeline().
		Database().
		UnnestWithAlias("neighbors", "unnestedNeighbors", firestore.WithUnnestIndexField("index")).
		Execute(ctx).Results().GetAll()
	if err != nil {
		return err
	}

	// Output
	// { "identifier": 1, "neighbors": [ "Alice", "Cathy" ],
	//   "unnestedNeighbors": "Alice", "index": 0 }
	// { "identifier": 1, "neighbors": [ "Alice", "Cathy" ],
	//   "unnestedNeighbors": "Cathy", "index": 1 }
	// { "identifier": 3, "neighbors": "Bob", "index": nil}
	// [END unnest_edge_cases]
	fmt.Fprintln(w, results)
	return nil
}

func unnestSyntaxExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START unnest_syntax]
	snapshot := client.Pipeline().Collection("users").
		UnnestWithAlias("scores", "userScore", firestore.WithUnnestIndexField("attempt")).
		Execute(ctx)
	// [END unnest_syntax]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func unnestAliasIndexDataExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START unnest_alias_index_data]
	client.Collection("users").Add(ctx, map[string]any{
		"name":      "foo",
		"scores":    []int{5, 4},
		"userScore": 0,
	})
	client.Collection("users").Add(ctx, map[string]any{
		"name":    "bar",
		"scores":  []int{1, 3},
		"attempt": 5,
	})
	// [END unnest_alias_index_data]
	return nil
}

func unnestAliasIndexExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START unnest_alias_index]
	snapshot := client.Pipeline().Collection("users").
		UnnestWithAlias("scores", "userScore", firestore.WithUnnestIndexField("attempt")).
		Execute(ctx)
	// [END unnest_alias_index]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func unnestNonArrayDataExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START unnest_nonarray_data]
	client.Collection("users").Add(ctx, map[string]any{
		"name":   "foo",
		"scores": 1,
	})
	client.Collection("users").Add(ctx, map[string]any{
		"name":   "bar",
		"scores": nil,
	})
	client.Collection("users").Add(ctx, map[string]any{
		"name": "qux",
		"scores": map[string]any{
			"backupScores": 1,
		},
	})
	// [END unnest_nonarray_data]
	return nil
}

func unnestNonArrayExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START unnest_nonarray]
	snapshot := client.Pipeline().Collection("users").
		UnnestWithAlias("scores", "userScore", firestore.WithUnnestIndexField("attempt")).
		Execute(ctx)
	// [END unnest_nonarray]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func unnestEmptyArrayDataExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START unnest_empty_array_data]
	client.Collection("users").Add(ctx, map[string]any{
		"name":   "foo",
		"scores": []int{5, 4},
	})
	client.Collection("users").Add(ctx, map[string]any{
		"name":   "bar",
		"scores": []int{},
	})
	// [END unnest_empty_array_data]
	return nil
}

func unnestEmptyArrayExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START unnest_empty_array]
	snapshot := client.Pipeline().Collection("users").
		UnnestWithAlias("scores", "userScore", firestore.WithUnnestIndexField("attempt")).
		Execute(ctx)
	// [END unnest_empty_array]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func unnestPreserveEmptyArrayExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START unnest_preserve_empty_array]
	userScore, err := client.Pipeline().
		Collection("users").
		Unnest(firestore.Conditional(
			firestore.FieldOf("scores").Equal([]any{}),
			firestore.Array(firestore.FieldOf("scores")),
			firestore.FieldOf("scores"),
		).As("userScore"), firestore.WithUnnestIndexField("attempt")).
		Execute(ctx).Results().GetAll()
	if err != nil {
		return err
	}
	// [END unnest_preserve_empty_array]
	fmt.Fprintln(w, userScore)
	return nil
}

func unnestNestedDataExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START unnest_nested_data]
	client.Collection("users").Add(ctx, map[string]any{
		"name": "foo",
		"record": []any{
			map[string]any{
				"scores": []int{5, 4},
				"avg":    4.5,
			},
			map[string]any{
				"scores":  []int{1, 3},
				"old_avg": 2,
			},
		},
	})
	// [END unnest_nested_data]
	return nil
}

func unnestNestedExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START unnest_nested]
	snapshot := client.Pipeline().Collection("users").
		UnnestWithAlias("record", "record").
		UnnestWithAlias("record.scores", "userScore", firestore.WithUnnestIndexField("attempt")).
		Execute(ctx)
	// [END unnest_nested]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}
