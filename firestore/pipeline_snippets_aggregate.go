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

func aggregateGroups(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START aggregate_groups]
	snapshot := client.Pipeline().
		Collection("books").
		Aggregate(
			firestore.Accumulators(firestore.Average("rating").As("avg_rating")),
			firestore.WithAggregateGroups("genre"),
		).
		Execute(ctx)
	// [END aggregate_groups]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func aggregateDistinct(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START aggregate_distinct]
	snapshot := client.Pipeline().
		Collection("books").
		Distinct(firestore.Fields(
			firestore.ToUpper(firestore.FieldOf("author")).As("author"),
			firestore.FieldOf("genre"),
		)).
		Execute(ctx)
	// [END aggregate_distinct]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func aggregateSyntaxExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START aggregate_syntax]
	snapshot := client.Pipeline().Collection("cities").
		Aggregate(firestore.Accumulators(
			firestore.CountAll().As("total"),
			firestore.Average("population").As("averagePopulation"),
		)).
		Execute(ctx)
	// [END aggregate_syntax]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func aggregateGroupSyntax(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START aggregate_group_syntax]
	snapshot := client.Pipeline().CollectionGroup("cities").
		Aggregate(
			firestore.Accumulators(
				firestore.CountAll().As("cities"),
				firestore.Sum("population").As("totalPopulation"),
			),
			firestore.WithAggregateGroups(firestore.FieldOf("location.state").As("state")),
		).
		Execute(ctx)
	// [END aggregate_group_syntax]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func aggregateExampleData(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START aggregate_data]
	client.Collection("cities").Doc("SF").Set(ctx, map[string]any{
		"name":       "San Francisco",
		"state":      "CA",
		"country":    "USA",
		"population": 870000,
	})
	client.Collection("cities").Doc("LA").Set(ctx, map[string]any{
		"name":       "Los Angeles",
		"state":      "CA",
		"country":    "USA",
		"population": 3970000,
	})
	client.Collection("cities").Doc("NY").Set(ctx, map[string]any{
		"name":       "New York",
		"state":      "NY",
		"country":    "USA",
		"population": 8530000,
	})
	client.Collection("cities").Doc("TOR").Set(ctx, map[string]any{
		"name":       "Toronto",
		"state":      nil,
		"country":    "Canada",
		"population": 2930000,
	})
	client.Collection("cities").Doc("MEX").Set(ctx, map[string]any{
		"name":       "Mexico City",
		"state":      nil,
		"country":    "Mexico",
		"population": 9200000,
	})
	// [END aggregate_data]
	return nil
}

func aggregateWithoutGroupExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START aggregate_without_group]
	snapshot := client.Pipeline().Collection("cities").
		Aggregate(firestore.Accumulators(
			firestore.CountAll().As("total"),
			firestore.Average("population").As("averagePopulation"),
		)).
		Execute(ctx)
	// [END aggregate_without_group]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func aggregateGroupExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START aggregate_group_example]
	snapshot := client.Pipeline().Collection("cities").
		Aggregate(
			firestore.Accumulators(
				firestore.CountAll().As("numberOfCities"),
				firestore.Maximum("population").As("maxPopulation"),
			),
			firestore.WithAggregateGroups("country", "state"),
		).
		Execute(ctx)
	// [END aggregate_group_example]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func aggregateGroupComplexExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START aggregate_group_complex]
	snapshot := client.Pipeline().Collection("cities").
		Aggregate(
			firestore.Accumulators(firestore.Sum("population").As("totalPopulation")),
			firestore.WithAggregateGroups(firestore.FieldOf("state").Equal(nil).As("stateIsNull")),
		).
		Execute(ctx)
	// [END aggregate_group_complex]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func distinctSyntaxExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START distinct_syntax]
	cities1, err := client.Pipeline().Collection("cities").Distinct(firestore.Fields("country")).Execute(ctx).Results().GetAll()
	if err != nil {
		return err
	}

	cities2, err := client.Pipeline().Collection("cities").
		Distinct(firestore.Fields(
			firestore.ToLower(firestore.FieldOf("state")).As("normalizedState"),
			firestore.FieldOf("country"),
		)).
		Execute(ctx).Results().GetAll()
	if err != nil {
		return err
	}
	// [END distinct_syntax]
	fmt.Fprintln(w, cities1)
	fmt.Fprintln(w, cities2)
	return nil
}

func distinctExampleData(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START distinct_data]
	client.Collection("cities").Doc("SF").Set(ctx, map[string]any{
		"name":    "San Francisco",
		"state":   "CA",
		"country": "USA",
	})
	client.Collection("cities").Doc("LA").Set(ctx, map[string]any{
		"name":    "Los Angeles",
		"state":   "CA",
		"country": "USA",
	})
	client.Collection("cities").Doc("NY").Set(ctx, map[string]any{
		"name":    "New York",
		"state":   "NY",
		"country": "USA",
	})
	client.Collection("cities").Doc("TOR").Set(ctx, map[string]any{
		"name":    "Toronto",
		"state":   nil,
		"country": "Canada",
	})
	client.Collection("cities").Doc("MEX").Set(ctx, map[string]any{
		"name":    "Mexico City",
		"state":   nil,
		"country": "Mexico",
	})
	// [END distinct_data]
	return nil
}

func distinctExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START distinct_example]
	snapshot := client.Pipeline().Collection("cities").Distinct(firestore.Fields("country")).Execute(ctx)
	// [END distinct_example]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func distinctExpressionsExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START distinct_expressions]
	results, err := client.Pipeline().Collection("cities").
		Distinct(firestore.Fields(
			firestore.ToLower(firestore.FieldOf("state")).As("normalizedState"),
			firestore.FieldOf("country"),
		)).
		Execute(ctx).Results().GetAll()
	if err != nil {
		return err
	}
	// [END distinct_expressions]
	fmt.Fprintln(w, results)
	return nil
}

func countFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START count_function]
	// Total number of books in the collection
	countAll, err := client.Pipeline().Collection("books").
		Aggregate(firestore.Accumulators(firestore.CountAll().As("count"))).
		Execute(ctx).Results().GetAll()
	if err != nil {
		return err
	}

	// Number of books with nonnull `ratings` field
	countField, err := client.Pipeline().
		Collection("books").
		Aggregate(firestore.Accumulators(firestore.Count("ratings").As("count"))).
		Execute(ctx).Results().GetAll()
	if err != nil {
		return err
	}
	// [END count_function]
	fmt.Fprintln(w, countAll, countField)
	return nil
}

func countIfFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START count_if]
	snapshot := client.Pipeline().
		Collection("books").
		Aggregate(firestore.Accumulators(
			firestore.CountIf(firestore.FieldOf("rating").GreaterThan(4)).As("filteredCount"),
		)).
		Execute(ctx)
	// [END count_if]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func countDistinctFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START count_distinct]
	snapshot := client.Pipeline().
		Collection("books").
		Aggregate(firestore.Accumulators(
			firestore.CountDistinct("author").As("unique_authors"),
		)).
		Execute(ctx)
	// [END count_distinct]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func sumFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START sum_function]
	snapshot := client.Pipeline().
		Collection("cities").
		Aggregate(firestore.Accumulators(
			firestore.Sum("population").As("totalPopulation"),
		)).
		Execute(ctx)
	// [END sum_function]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func avgFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START avg_function]
	snapshot := client.Pipeline().
		Collection("cities").
		Aggregate(firestore.Accumulators(
			firestore.Average("population").As("averagePopulation"),
		)).
		Execute(ctx)
	// [END avg_function]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func minFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START min_function]
	snapshot := client.Pipeline().
		Collection("books").
		Aggregate(firestore.Accumulators(
			firestore.Minimum("price").As("minimumPrice"),
		)).
		Execute(ctx)
	// [END min_function]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func maxFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START max_function]
	snapshot := client.Pipeline().
		Collection("books").
		Aggregate(firestore.Accumulators(
			firestore.Maximum("price").As("maximumPrice"),
		)).
		Execute(ctx)
	// [END max_function]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}
