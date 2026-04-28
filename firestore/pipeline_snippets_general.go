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
	"google.golang.org/api/iterator"
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

func sort(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_sort]
	snapshot := client.Pipeline().
		Collection("books").
		Sort(firestore.Orders(
			firestore.Descending(firestore.FieldOf("release_date")),
			firestore.Ascending(firestore.FieldOf("author")),
		)).
		Execute(ctx)
	// [END firestore_sort]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func sortComparison(w io.Writer, client *firestore.Client) error {
	// [START firestore_sort_comparison]
	query := client.Collection("cities").
		OrderBy("state", firestore.Asc).
		OrderBy("population", firestore.Desc)

	pipeline := client.Pipeline().
		Collection("books").
		Sort(firestore.Orders(
			firestore.Descending(firestore.FieldOf("release_date")),
			firestore.Ascending(firestore.FieldOf("author")),
		))
	// [END firestore_sort_comparison]
	fmt.Fprintln(w, query)
	fmt.Fprintln(w, pipeline)
	return nil
}

func sortSyntaxExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_sort_syntax]
	snapshot := client.Pipeline().Collection("cities").
		Sort(firestore.Orders(firestore.Ascending(firestore.FieldOf("population")))).
		Execute(ctx)
	// [END firestore_sort_syntax]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func sortSyntaxExample2(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_sort_syntax_2]
	snapshot := client.Pipeline().Collection("cities").
		Sort(firestore.Orders(firestore.Ascending(firestore.CharLength(firestore.FieldOf("name"))))).
		Execute(ctx)
	// [END firestore_sort_syntax_2]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func sortDocumentIDExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_sort_document_id]
	snapshot := client.Pipeline().Collection("cities").
		Sort(firestore.Orders(
			firestore.Ascending(firestore.FieldOf("country")),
			firestore.Ascending(firestore.FieldOf("__name__")),
		)).
		Execute(ctx)
	// [END firestore_sort_document_id]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func aggregateGroups(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_aggregate_groups]
	snapshot := client.Pipeline().
		Collection("books").
		Aggregate(
			firestore.Accumulators(firestore.Average("rating").As("avg_rating")),
			firestore.WithAggregateGroups("genre"),
		).
		Execute(ctx)
	// [END firestore_aggregate_groups]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func aggregateDistinct(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_aggregate_distinct]
	snapshot := client.Pipeline().
		Collection("books").
		Distinct(firestore.Fields(
			firestore.ToUpper(firestore.FieldOf("author")).As("author"),
			firestore.FieldOf("genre"),
		)).
		Execute(ctx)
	// [END firestore_aggregate_distinct]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func aggregateSyntaxExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_aggregate_syntax]
	snapshot := client.Pipeline().Collection("cities").
		Aggregate(firestore.Accumulators(
			firestore.CountAll().As("total"),
			firestore.Average("population").As("averagePopulation"),
		)).
		Execute(ctx)
	// [END firestore_aggregate_syntax]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func aggregateGroupSyntax(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_aggregate_group_syntax]
	snapshot := client.Pipeline().CollectionGroup("cities").
		Aggregate(
			firestore.Accumulators(
				firestore.CountAll().As("cities"),
				firestore.Sum("population").As("totalPopulation"),
			),
			firestore.WithAggregateGroups(firestore.FieldOf("location.state").As("state")),
		).
		Execute(ctx)
	// [END firestore_aggregate_group_syntax]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func aggregateExampleData(_ io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_aggregate_data]
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
	// [END firestore_aggregate_data]
	return nil
}

func aggregateWithoutGroupExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_aggregate_without_group]
	snapshot := client.Pipeline().Collection("cities").
		Aggregate(firestore.Accumulators(
			firestore.CountAll().As("total"),
			firestore.Average("population").As("averagePopulation"),
		)).
		Execute(ctx)
	// [END firestore_aggregate_without_group]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func aggregateGroupExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_aggregate_group_example]
	snapshot := client.Pipeline().Collection("cities").
		Aggregate(
			firestore.Accumulators(
				firestore.CountAll().As("numberOfCities"),
				firestore.Maximum("population").As("maxPopulation"),
			),
			firestore.WithAggregateGroups("country", "state"),
		).
		Execute(ctx)
	// [END firestore_aggregate_group_example]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func aggregateGroupComplexExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_aggregate_group_complex]
	snapshot := client.Pipeline().Collection("cities").
		Aggregate(
			firestore.Accumulators(firestore.Sum("population").As("totalPopulation")),
			firestore.WithAggregateGroups(firestore.FieldOf("state").Equal(nil).As("stateIsNull")),
		).
		Execute(ctx)
	// [END firestore_aggregate_group_complex]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func distinctSyntaxExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_distinct_syntax]
	cities1, err := client.Pipeline().Collection("cities").Distinct(firestore.Fields("country")).Execute(ctx).Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "GetAll failed: %v", err)
		return err
	}

	cities2, err := client.Pipeline().Collection("cities").
		Distinct(firestore.Fields(
			firestore.ToLower(firestore.FieldOf("state")).As("normalizedState"),
			firestore.FieldOf("country"),
		)).
		Execute(ctx).Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "GetAll failed: %v", err)
		return err
	}
	// [END firestore_distinct_syntax]
	fmt.Fprintln(w, cities1)
	fmt.Fprintln(w, cities2)
	return nil
}

func distinctExampleData(_ io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_distinct_data]
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
	// [END firestore_distinct_data]
	return nil
}

func distinctExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_distinct_example]
	snapshot := client.Pipeline().Collection("cities").Distinct(firestore.Fields("country")).Execute(ctx)
	// [END firestore_distinct_example]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func distinctExpressionsExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_distinct_expressions]
	results, err := client.Pipeline().Collection("cities").
		Distinct(firestore.Fields(
			firestore.ToLower(firestore.FieldOf("state")).As("normalizedState"),
			firestore.FieldOf("country"),
		)).
		Execute(ctx).Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "GetAll failed: %v", err)
		return err
	}
	// [END firestore_distinct_expressions]
	fmt.Fprintln(w, results)
	return nil
}

func selectSyntaxExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_select_syntax]
	snapshot := client.Pipeline().Collection("cities").
		Select(firestore.Fields(
			firestore.StringConcat(firestore.FieldOf("name"), ", ", firestore.FieldOf("location.country")).As("name"),
			firestore.FieldOf("population"),
		)).
		Execute(ctx)
	// [END firestore_select_syntax]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func selectPositionDataExample(_ io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_select_position_data]
	client.Collection("cities").Doc("SF").Set(ctx, map[string]any{
		"name":       "San Francisco",
		"population": 800000,
		"location": map[string]any{
			"country": "USA",
			"state":   "California",
		},
	})
	client.Collection("cities").Doc("TO").Set(ctx, map[string]any{
		"name":       "Toronto",
		"population": 3000000,
		"location": map[string]any{
			"country":  "Canada",
			"province": "Ontario",
		},
	})
	// [END firestore_select_position_data]
	return nil
}

func selectPositionExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_select_position]
	snapshot := client.Pipeline().Collection("cities").
		Where(firestore.FieldOf("location.country").Equal("Canada")).
		Select(firestore.Fields(
			firestore.StringConcat(firestore.FieldOf("name"), ", ", firestore.FieldOf("location.country")).As("name"),
			firestore.FieldOf("population"),
		)).
		Execute(ctx)
	// [END firestore_select_position]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func selectBadPositionExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_select_bad_position]
	snapshot := client.Pipeline().Collection("cities").
		Select(firestore.Fields(
			firestore.StringConcat(firestore.FieldOf("name"), ", ", firestore.FieldOf("location.country")).As("name"),
			firestore.FieldOf("population"),
		)).
		Where(firestore.FieldOf("location.country").Equal("Canada")).
		Execute(ctx)
	// [END firestore_select_bad_position]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func selectNestedDataExample(_ io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_select_nested_data]
	client.Collection("cities").Doc("SF").Set(ctx, map[string]any{
		"name":       "San Francisco",
		"population": 800000,
		"location": map[string]any{
			"country": "USA",
			"state":   "California",
		},
		"landmarks": []string{"Golden Gate Bridge", "Alcatraz"},
	})
	client.Collection("cities").Doc("TO").Set(ctx, map[string]any{
		"name":       "Toronto",
		"population": 3000000,
		"province":   "ON",
		"location": map[string]any{
			"country":  "Canada",
			"province": "Ontario",
		},
		"landmarks": []string{"CN Tower", "Casa Loma"},
	})
	client.Collection("cities").Doc("AT").Set(ctx, map[string]any{
		"name":       "Atlantis",
		"population": nil,
	})
	// [END firestore_select_nested_data]
	return nil
}

func selectNestedExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_select_nested]
	snapshot := client.Pipeline().Collection("cities").
		Select(firestore.Fields(
			firestore.FieldOf("name").As("city"),
			firestore.FieldOf("location.country").As("country"),
			firestore.ArrayGet(firestore.FieldOf("landmarks"), 0).As("topLandmark"),
		)).
		Execute(ctx)
	// [END firestore_select_nested]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func addFieldsSyntaxExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_add_fields_syntax]
	snapshot := client.Pipeline().Collection("users").
		AddFields(firestore.Selectables(
			firestore.StringConcat(firestore.FieldOf("firstName"), " ", firestore.FieldOf("lastName")).As("fullName"),
		)).
		Execute(ctx)
	// [END firestore_add_fields_syntax]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func addFieldsOverlapExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_add_fields_overlap]
	snapshot := client.Pipeline().Collection("users").
		AddFields(firestore.Selectables(firestore.Abs(firestore.FieldOf("age")).As("age"))).
		AddFields(firestore.Selectables(firestore.Add(firestore.FieldOf("age"), 10).As("age"))).
		Execute(ctx)
	// [END firestore_add_fields_overlap]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func addFieldsNestingExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_add_fields_nesting]
	snapshot := client.Pipeline().Collection("users").
		AddFields(firestore.Selectables(
			firestore.ToLower(firestore.FieldOf("address.city")).As("address.city"),
		)).
		Execute(ctx)
	// [END firestore_add_fields_nesting]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func removeFieldsSyntaxExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_remove_fields_syntax]
	snapshot := client.Pipeline().Collection("cities").
		RemoveFields(firestore.Fields("population", "location.state")).
		Execute(ctx)
	// [END firestore_remove_fields_syntax]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func removeFieldsNestedDataExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_remove_fields_nested_data]
	client.Collection("cities").Doc("SF").Set(ctx, map[string]any{
		"name": "San Francisco",
		"location": map[string]any{
			"country": "USA",
			"state":   "California",
		},
	})
	client.Collection("cities").Doc("TO").Set(ctx, map[string]any{
		"name": "Toronto",
		"location": map[string]any{
			"country":  "Canada",
			"province": "Ontario",
		},
	})
	// [END firestore_remove_fields_nested_data]
	return nil
}

func removeFieldsNestedExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_remove_fields_nested]
	snapshot := client.Pipeline().Collection("cities").
		RemoveFields(firestore.Fields("location.state")).
		Execute(ctx)
	// [END firestore_remove_fields_nested]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

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

func creatingIndexes(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_query_example]
	snapshot := client.Pipeline().
		Collection("books").
		Where(firestore.FieldOf("published").LessThan(1900)).
		Where(firestore.FieldOf("genre").Equal("Science Fiction")).
		Where(firestore.FieldOf("rating").GreaterThan(4.3)).
		Sort(firestore.Orders(firestore.Descending(firestore.FieldOf("published")))).
		Execute(ctx)
	// [END firestore_query_example]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func sparseIndexes(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_sparse_index_example]
	snapshot := client.Pipeline().
		Collection("books").
		Where(firestore.FieldOf("category").Like("%fantasy%")).
		Execute(ctx)
	// [END firestore_sparse_index_example]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func sparseIndexes2(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_sparse_index_example_2]
	snapshot := client.Pipeline().
		Collection("books").
		Sort(firestore.Orders(firestore.Ascending(firestore.FieldOf("release_date")))).
		Execute(ctx)
	// [END firestore_sparse_index_example_2]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func coveredQuery(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_covered_query]
	snapshot := client.Pipeline().
		Collection("books").
		Where(firestore.FieldOf("category").Like("%fantasy%")).
		Where(firestore.FieldOf("title").FieldExists()).
		Where(firestore.FieldOf("author").FieldExists()).
		Select(firestore.Fields("title", "author")).
		Execute(ctx)
	// [END firestore_covered_query]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func pagination(w io.Writer, client *firestore.Client) error {
	// [START firestore_pagination_not_supported_preview]
	// Existing pagination via `StartAt()`
	query := client.Collection("cities").OrderBy("population", firestore.Asc).StartAt(1000000)

	pipeline := client.Pipeline().
		Collection("cities").
		Where(firestore.FieldOf("population").GreaterThanOrEqual(1000000)).
		Sort(firestore.Orders(firestore.Descending(firestore.FieldOf("population"))))
	// [END firestore_pagination_not_supported_preview]
	fmt.Fprintln(w, query)
	fmt.Fprintln(w, pipeline)
	return nil
}

func collectionStage(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_collection_example]
	snapshot := client.Pipeline().
		Collection("users/bob/games").
		Sort(firestore.Orders(firestore.Ascending(firestore.FieldOf("name")))).
		Execute(ctx)
	// [END firestore_collection_example]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func collectionGroupStage(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_collection_group_example]
	snapshot := client.Pipeline().
		CollectionGroup("games").
		Sort(firestore.Orders(firestore.Ascending(firestore.FieldOf("name")))).
		Execute(ctx)
	// [END firestore_collection_group_example]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func databaseStage(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_database_example]
	// Count all documents in the database
	snapshot := client.Pipeline().
		Database().
		Aggregate(firestore.Accumulators(firestore.CountAll().As("total"))).
		Execute(ctx)
	// [END firestore_database_example]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func documentsStage(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_documents_example]
	snapshot := client.Pipeline().
		Documents([]*firestore.DocumentRef{
			client.Collection("cities").Doc("SF"),
			client.Collection("cities").Doc("DC"),
			client.Collection("cities").Doc("NY"),
		}).
		Execute(ctx)
	// [END firestore_documents_example]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func replaceWithStage(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_initial_data]
	client.Collection("cities").Doc("SF").Set(ctx, map[string]any{
		"name":       "San Francisco",
		"population": 800000,
		"location": map[string]any{
			"country": "USA",
			"state":   "California",
		},
	})
	client.Collection("cities").Doc("TO").Set(ctx, map[string]any{
		"name":       "Toronto",
		"population": 3000000,
		"province":   "ON",
		"location": map[string]any{
			"country":  "Canada",
			"province": "Ontario",
		},
	})
	client.Collection("cities").Doc("NY").Set(ctx, map[string]any{
		"name":       "New York",
		"population": 8500000,
		"location": map[string]any{
			"country": "USA",
			"state":   "New York",
		},
	})
	client.Collection("cities").Doc("AT").Set(ctx, map[string]any{
		"name": "Atlantis",
	})
	// [END firestore_initial_data]

	// [START firestore_full_replace]
	snapshot := client.Pipeline().
		Collection("cities").
		ReplaceWith(firestore.FieldOf("location")).
		Execute(ctx)
	// [END firestore_full_replace]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}

	// [START firestore_map_merge_overwrite]
	// unsupported in client SDKs for now
	// [END firestore_map_merge_overwrite]
	fmt.Fprintln(w, results)
	return nil
}

func sampleStage(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_sample_example]
	// Get a sample of 100 documents in a database
	results1, err := client.Pipeline().Database().Sample(firestore.WithDocLimit(100)).Execute(ctx).Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "GetAll failed: %v", err)
		return err
	}

	// Randomly shuffle a list of 3 documents
	results2, err := client.Pipeline().
		Documents([]*firestore.DocumentRef{
			client.Collection("cities").Doc("SF"),
			client.Collection("cities").Doc("NY"),
			client.Collection("cities").Doc("DC"),
		}).
		Sample(firestore.WithDocLimit(3)).
		Execute(ctx).Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "GetAll failed: %v", err)
		return err
	}
	// [END firestore_sample_example]
	fmt.Fprintln(w, results1)
	fmt.Fprintln(w, results2)
	return nil
}

func samplePercent(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_sample_percent]
	// Get a sample of on average 50% of the documents in the database
	snapshot := client.Pipeline().
		Database().
		Sample(firestore.WithPercentage(0.5)).
		Execute(ctx)
	// [END firestore_sample_percent]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func sampleSyntaxExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_sample_syntax]
	sampled1, err := client.Pipeline().Database().Sample(firestore.WithDocLimit(50)).Execute(ctx).Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "GetAll failed: %v", err)
		return err
	}

	sampled2, err := client.Pipeline().Database().Sample(firestore.WithPercentage(0.5)).Execute(ctx).Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "GetAll failed: %v", err)
		return err
	}
	// [END firestore_sample_syntax]
	fmt.Fprintln(w, sampled1)
	fmt.Fprintln(w, sampled2)
	return nil
}

func sampleDocumentsDataExample(_ io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_sample_documents_data]
	client.Collection("cities").Doc("SF").Set(ctx, map[string]any{
		"name":  "San Francisco",
		"state": "California",
	})
	client.Collection("cities").Doc("NYC").Set(ctx, map[string]any{
		"name":  "New York City",
		"state": "New York",
	})
	client.Collection("cities").Doc("CHI").Set(ctx, map[string]any{
		"name":  "Chicago",
		"state": "Illinois",
	})
	// [END firestore_sample_documents_data]
	return nil
}

func sampleDocumentsExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_sample_documents]
	snapshot := client.Pipeline().Collection("cities").Sample(firestore.WithDocLimit(1)).Execute(ctx)
	// [END firestore_sample_documents]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func sampleAllDocumentsExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_sample_all_documents]
	snapshot := client.Pipeline().Collection("cities").Sample(firestore.WithDocLimit(5)).Execute(ctx)
	// [END firestore_sample_all_documents]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func samplePercentageDataExample(_ io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_sample_percentage_data]
	client.Collection("cities").Doc("SF").Set(ctx, map[string]any{
		"name":  "San Francisco",
		"state": "California",
	})
	client.Collection("cities").Doc("NYC").Set(ctx, map[string]any{
		"name":  "New York City",
		"state": "New York",
	})
	client.Collection("cities").Doc("CHI").Set(ctx, map[string]any{
		"name":  "Chicago",
		"state": "Illinois",
	})
	client.Collection("cities").Doc("ATL").Set(ctx, map[string]any{
		"name":  "Atlanta",
		"state": "Georgia",
	})
	// [END firestore_sample_percentage_data]
	return nil
}

func samplePercentageExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_sample_percentage]
	snapshot := client.Pipeline().Collection("cities").Sample(firestore.WithPercentage(0.5)).Execute(ctx)
	// [END firestore_sample_percentage]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func unionStage(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_union_stage]
	snapshot := client.Pipeline().
		Collection("cities/SF/restaurants").
		Where(firestore.FieldOf("type").Equal("Chinese")).
		Union(
			client.Pipeline().
				Collection("cities/NY/restaurants").
				Where(firestore.FieldOf("type").Equal("Italian")),
		).
		Where(firestore.FieldOf("rating").GreaterThanOrEqual(4.5)).
		Sort(firestore.Orders(firestore.Descending(firestore.FieldOf("__name__")))).
		Execute(ctx)
	// [END firestore_union_stage]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func unionStageStable(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_union_stage_stable]
	snapshot := client.Pipeline().
		Collection("cities/SF/restaurants").
		Where(firestore.FieldOf("type").Equal("Chinese")).
		Union(
			client.Pipeline().
				Collection("cities/NY/restaurants").
				Where(firestore.FieldOf("type").Equal("Italian")),
		).
		Where(firestore.FieldOf("rating").GreaterThanOrEqual(4.5)).
		Sort(firestore.Orders(firestore.Descending(firestore.FieldOf("__name__")))).
		Execute(ctx)
	// [END firestore_union_stage_stable]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func offsetSyntaxExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_offset_syntax]
	snapshot := client.Pipeline().Collection("cities").Offset(10).Execute(ctx)
	// [END firestore_offset_syntax]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func collectionInputSyntaxExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_collection_input_syntax]
	snapshot := client.Pipeline().Collection("cities/SF/departments").Execute(ctx)
	// [END firestore_collection_input_syntax]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func collectionInputExampleData(_ io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_collection_input_data]
	client.Collection("cities").Doc("SF").Set(ctx, map[string]any{
		"name":  "San Francisco",
		"state": "California",
	})
	client.Collection("cities").Doc("NYC").Set(ctx, map[string]any{
		"name":  "New York City",
		"state": "New York",
	})
	client.Collection("cities").Doc("CHI").Set(ctx, map[string]any{
		"name":  "Chicago",
		"state": "Illinois",
	})
	client.Collection("states").Doc("CA").Set(ctx, map[string]any{
		"name": "California",
	})
	// [END firestore_collection_input_data]
	return nil
}

func collectionInputExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_collection_input]
	snapshot := client.Pipeline().Collection("cities").
		Sort(firestore.Orders(firestore.Ascending(firestore.FieldOf("name")))).
		Execute(ctx)
	// [END firestore_collection_input]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func subcollectionInputExampleData(_ io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_subcollection_input_data]
	client.Collection("cities/SF/departments").Doc("building").Set(ctx, map[string]any{
		"name":      "SF Building Department",
		"employees": 750,
	})
	client.Collection("cities/NY/departments").Doc("building").Set(ctx, map[string]any{
		"name":      "NY Building Department",
		"employees": 1000,
	})
	client.Collection("cities/CHI/departments").Doc("building").Set(ctx, map[string]any{
		"name":      "CHI Building Department",
		"employees": 900,
	})
	client.Collection("cities/NY/departments").Doc("finance").Set(ctx, map[string]any{
		"name":      "NY Finance Department",
		"employees": 1200,
	})
	// [END firestore_subcollection_input_data]
	return nil
}

func subcollectionInputExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_subcollection_input]
	snapshot := client.Pipeline().Collection("cities/NY/departments").
		Sort(firestore.Orders(firestore.Ascending(firestore.FieldOf("employees")))).
		Execute(ctx)
	// [END firestore_subcollection_input]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func collectionGroupInputSyntaxExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_collection_group_input_syntax]
	snapshot := client.Pipeline().CollectionGroup("departments").Execute(ctx)
	// [END firestore_collection_group_input_syntax]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func collectionGroupInputExampleData(_ io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_collection_group_data]
	client.Collection("cities/SF/departments").Doc("building").Set(ctx, map[string]any{
		"name":      "SF Building Department",
		"employees": 750,
	})
	client.Collection("cities/NY/departments").Doc("building").Set(ctx, map[string]any{
		"name":      "NY Building Department",
		"employees": 1000,
	})
	client.Collection("cities/CHI/departments").Doc("building").Set(ctx, map[string]any{
		"name":      "CHI Building Department",
		"employees": 900,
	})
	client.Collection("cities/NY/departments").Doc("finance").Set(ctx, map[string]any{
		"name":      "NY Finance Department",
		"employees": 1200,
	})
	// [END firestore_collection_group_data]
	return nil
}

func collectionGroupInputExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_collection_group_input]
	snapshot := client.Pipeline().CollectionGroup("departments").
		Sort(firestore.Orders(firestore.Ascending(firestore.FieldOf("employees")))).
		Execute(ctx)
	// [END firestore_collection_group_input]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func databaseInputSyntaxExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_database_syntax]
	snapshot := client.Pipeline().Database().Execute(ctx)
	// [END firestore_database_syntax]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func databaseInputSyntaxExampleData(_ io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_database_input_data]
	client.Collection("cities").Doc("SF").Set(ctx, map[string]any{
		"name":       "San Francisco",
		"state":      "California",
		"population": 800000,
	})
	client.Collection("states").Doc("CA").Set(ctx, map[string]any{
		"name":       "California",
		"population": 39000000,
	})
	client.Collection("countries").Doc("USA").Set(ctx, map[string]any{
		"name":       "United States of America",
		"population": 340000000,
	})
	// [END firestore_database_input_data]
	return nil
}

func databaseInputExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_database_input]
	snapshot := client.Pipeline().Database().
		Sort(firestore.Orders(firestore.Ascending(firestore.FieldOf("population")))).
		Execute(ctx)
	// [END firestore_database_input]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func documentInputSyntaxExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_document_input_syntax]
	snapshot := client.Pipeline().
		Documents([]*firestore.DocumentRef{
			client.Collection("cities").Doc("SF"),
			client.Collection("cities").Doc("NY"),
		}).
		Execute(ctx)
	// [END firestore_document_input_syntax]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func documentInputExampleData(_ io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_document_input_data]
	client.Collection("cities").Doc("SF").Set(ctx, map[string]any{
		"name":  "San Francisco",
		"state": "California",
	})
	client.Collection("cities").Doc("NYC").Set(ctx, map[string]any{
		"name":  "New York City",
		"state": "New York",
	})
	client.Collection("cities").Doc("CHI").Set(ctx, map[string]any{
		"name":  "Chicago",
		"state": "Illinois",
	})
	// [END firestore_document_input_data]
	return nil
}

func documentInputExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_document_input]
	snapshot := client.Pipeline().
		Documents([]*firestore.DocumentRef{
			client.Collection("cities").Doc("SF"),
			client.Collection("cities").Doc("NYC"),
		}).
		Sort(firestore.Orders(firestore.Ascending(firestore.FieldOf("name")))).
		Execute(ctx)
	// [END firestore_document_input]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func limitSyntaxExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_limit_syntax]
	snapshot := client.Pipeline().Collection("cities").Limit(10).Execute(ctx)
	// [END firestore_limit_syntax]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func unionSyntaxExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_union_syntax]
	snapshot := client.Pipeline().
		Collection("cities/SF/restaurants").
		Union(client.Pipeline().Collection("cities/NYC/restaurants")).
		Execute(ctx)
	// [END firestore_union_syntax]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func forceIndexExamples(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_force_index]
	// Force Planner to use Index ID CICAgOi36pgK
	snapshot1 := client.Pipeline().
		CollectionGroup("customers", firestore.WithForceIndex("CICAgOi36pgK")).
		Limit(100).
		Execute(ctx)
	// [END firestore_force_index]
	results1, err := snapshot1.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results1)

	// [START firestore_force_scan]
	// Force Planner to only do a collection scan
	snapshot2 := client.Pipeline().
		CollectionGroup("customers", firestore.WithForceIndex("primary")).
		Limit(100).
		Execute(ctx)
	// [END firestore_force_scan]
	results2, err := snapshot2.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results2)
	return nil
}

func stagesExpressionsExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_stages_expressions_example]
	nowMillis := firestore.ConstantOf(1712404800000) // Example timestamp
	trailing30Days := nowMillis.UnixMillisToTimestamp().TimestampSubtract("day", 30)

	snapshot := client.Pipeline().
		Collection("productViews").
		Where(firestore.FieldOf("viewedAt").GreaterThan(trailing30Days)).
		Aggregate(firestore.Accumulators(firestore.CountDistinct("productId").As("uniqueProductViews"))).
		Execute(ctx)
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	// [END firestore_stages_expressions_example]
	fmt.Fprintln(w, results)
	return nil
}

func wherePipeline(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_pipeline_where]
	results1, err := client.Pipeline().
		Collection("books").
		Where(firestore.FieldOf("rating").Equal(5)).
		Where(firestore.FieldOf("published").LessThan(1900)).
		Execute(ctx).Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "GetAll failed: %v", err)
		return err
	}

	results2, err := client.Pipeline().
		Collection("books").
		Where(firestore.And(
			firestore.FieldOf("rating").Equal(5),
			firestore.FieldOf("published").LessThan(1900),
		)).
		Execute(ctx).Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "GetAll failed: %v", err)
		return err
	}
	// [END firestore_pipeline_where]
	fmt.Fprintln(w, results1)
	fmt.Fprintln(w, results2)
	return nil
}

func whereEqualityExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_where_equality_example]
	snapshot := client.Pipeline().Collection("cities").
		Where(firestore.FieldOf("state").Equal("CA")).
		Execute(ctx)
	// [END firestore_where_equality_example]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func whereMultipleStagesExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_where_multiple_stages]
	snapshot := client.Pipeline().Collection("cities").
		Where(firestore.FieldOf("location.country").Equal("USA")).
		Where(firestore.FieldOf("population").GreaterThan(500000)).
		Execute(ctx)
	// [END firestore_where_multiple_stages]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func whereComplexExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_where_complex]
	snapshot := client.Pipeline().Collection("cities").
		Where(
			firestore.Or(
				firestore.Like(firestore.FieldOf("name"), "San%"),
				firestore.And(
					firestore.FieldOf("location.state").CharLength().GreaterThan(7),
					firestore.FieldOf("location.country").Equal("USA"),
				),
			),
		).
		Execute(ctx)
	// [END firestore_where_complex]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func whereStageOrderExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_where_stage_order]
	snapshot := client.Pipeline().Collection("cities").
		Limit(10).
		Where(firestore.FieldOf("location.country").Equal("USA")).
		Execute(ctx)
	// [END firestore_where_stage_order]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func createWhereData(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_create_where_data]
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
	// [END firestore_create_where_data]
	return nil
}

func whereHavingExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_where_having_example]
	snapshot := client.Pipeline().
		Collection("cities").
		Aggregate(
			firestore.Accumulators(firestore.Sum("population").As("totalPopulation")),
			firestore.WithAggregateGroups(firestore.FieldOf("state")),
		).
		Where(firestore.FieldOf("totalPopulation").GreaterThan(10000000)).
		Execute(ctx)
	// [END firestore_where_having_example]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func pipelineConcepts(w io.Writer, client *firestore.Client) error {
	// [START firestore_pipeline_concepts]
	pipeline := client.Pipeline().
		Collection("cities").
		Where(firestore.FieldOf("population").GreaterThan(100000)).
		Sort(firestore.Orders(firestore.Ascending(firestore.FieldOf("name")))).
		Limit(10)
	// [END firestore_pipeline_concepts]
	fmt.Fprintln(w, pipeline)
	return nil
}

func basicRead(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_basic_read]
	pipeline := client.Pipeline().Collection("users")
	snapshot := pipeline.Execute(ctx)
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	for _, result := range results {
		fmt.Fprintf(w, "%s => %v\n", result.Ref().ID, result.Data())
	}
	// or, one at a time
	it := pipeline.Execute(ctx).Results()
	for {
		result, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Fprintf(w, "it.Next failed: %v", err)
			return err
		}
		fmt.Fprintf(w, "%s => %v\n", result.Ref().ID, result.Data())
	}
	// [END firestore_basic_read]
	return nil
}

func pipelineInitialization(w io.Writer, projectID string) error {
	ctx := context.Background()
	// [START firestore_pipeline_initialization]
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		fmt.Fprintf(w, "firestore.NewClient failed: %v", err)
		return err
	}
	pipeline := client.Pipeline().Collection("books")
	// [END firestore_pipeline_initialization]
	fmt.Fprintln(w, pipeline)
	defer client.Close()
	return nil
}

func fieldVsConstants(w io.Writer, client *firestore.Client) error {
	// [START firestore_field_or_constant]
	pipeline := client.Pipeline().Collection("cities").
		Where(firestore.FieldOf("name").Equal(firestore.ConstantOf("Toronto")))
	// [END firestore_field_or_constant]
	fmt.Fprintln(w, pipeline)
	return nil
}

func inputStages(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_input_stages]
	// Return all restaurants in San Francisco
	results1, err := client.Pipeline().Collection("cities/sf/restaurants").Execute(ctx).Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "GetAll failed: %v", err)
		return err
	}

	// Return all restaurants
	results2, err := client.Pipeline().CollectionGroup("restaurants").Execute(ctx).Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "GetAll failed: %v", err)
		return err
	}

	// Return all documents across all collections in the database (the entire database)
	results3, err := client.Pipeline().Database().Execute(ctx).Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "GetAll failed: %v", err)
		return err
	}

	// Batch read of 3 documents
	results4, err := client.Pipeline().
		Documents([]*firestore.DocumentRef{
			client.Collection("cities").Doc("SF"),
			client.Collection("cities").Doc("DC"),
			client.Collection("cities").Doc("NY"),
		}).
		Execute(ctx).Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "GetAll failed: %v", err)
		return err
	}
	// [END firestore_input_stages]
	fmt.Fprintln(w, results1)
	fmt.Fprintln(w, results2)
	fmt.Fprintln(w, results3)
	fmt.Fprintln(w, results4)
	return nil
}

func unnestStage(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_unnest_stage]
	snapshot := client.Pipeline().
		Database().
		UnnestWithAlias("arrayField", "unnestedArrayField", firestore.WithUnnestIndexField("index")).
		Execute(ctx)
	// [END firestore_unnest_stage]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func unnestStageEmptyOrNonArray(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_unnest_edge_cases]
	// Input
	// { "identifier" : 1, "neighbors": [ "Alice", "Cathy" ] }
	// { "identifier" : 2, "neighbors": []                   }
	// { "identifier" : 3, "neighbors": "Bob"                }

	results, err := client.Pipeline().
		Database().
		UnnestWithAlias("neighbors", "unnestedNeighbors", firestore.WithUnnestIndexField("index")).
		Execute(ctx).Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "GetAll failed: %v", err)
		return err
	}

	// Output
	// { "identifier": 1, "neighbors": [ "Alice", "Cathy" ],
	//   "unnestedNeighbors": "Alice", "index": 0 }
	// { "identifier": 1, "neighbors": [ "Alice", "Cathy" ],
	//   "unnestedNeighbors": "Cathy", "index": 1 }
	// { "identifier": 3, "neighbors": "Bob", "index": nil}
	// [END firestore_unnest_edge_cases]
	fmt.Fprintln(w, results)
	return nil
}

func unnestSyntaxExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_unnest_syntax]
	snapshot := client.Pipeline().Collection("users").
		UnnestWithAlias("scores", "userScore", firestore.WithUnnestIndexField("attempt")).
		Execute(ctx)
	// [END firestore_unnest_syntax]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func unnestAliasIndexDataExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_unnest_alias_index_data]
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
	// [END firestore_unnest_alias_index_data]
	return nil
}

func unnestAliasIndexExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_unnest_alias_index]
	snapshot := client.Pipeline().Collection("users").
		UnnestWithAlias("scores", "userScore", firestore.WithUnnestIndexField("attempt")).
		Execute(ctx)
	// [END firestore_unnest_alias_index]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func unnestNonArrayDataExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_unnest_nonarray_data]
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
	// [END firestore_unnest_nonarray_data]
	return nil
}

func unnestNonArrayExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_unnest_nonarray]
	snapshot := client.Pipeline().Collection("users").
		UnnestWithAlias("scores", "userScore", firestore.WithUnnestIndexField("attempt")).
		Execute(ctx)
	// [END firestore_unnest_nonarray]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func unnestEmptyArrayDataExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_unnest_empty_array_data]
	client.Collection("users").Add(ctx, map[string]any{
		"name":   "foo",
		"scores": []int{5, 4},
	})
	client.Collection("users").Add(ctx, map[string]any{
		"name":   "bar",
		"scores": []int{},
	})
	// [END firestore_unnest_empty_array_data]
	return nil
}

func unnestEmptyArrayExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_unnest_empty_array]
	snapshot := client.Pipeline().Collection("users").
		UnnestWithAlias("scores", "userScore", firestore.WithUnnestIndexField("attempt")).
		Execute(ctx)
	// [END firestore_unnest_empty_array]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func unnestPreserveEmptyArrayExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_unnest_preserve_empty_array]
	userScore, err := client.Pipeline().
		Collection("users").
		Unnest(firestore.Conditional(
			firestore.FieldOf("scores").Equal([]any{}),
			firestore.Array(firestore.FieldOf("scores")),
			firestore.FieldOf("scores"),
		).As("userScore"), firestore.WithUnnestIndexField("attempt")).
		Execute(ctx).Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "GetAll failed: %v", err)
		return err
	}
	// [END firestore_unnest_preserve_empty_array]
	fmt.Fprintln(w, userScore)
	return nil
}

func unnestNestedDataExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_unnest_nested_data]
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
	// [END firestore_unnest_nested_data]
	return nil
}

func unnestNestedExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_unnest_nested]
	snapshot := client.Pipeline().Collection("users").
		UnnestWithAlias("record", "record").
		UnnestWithAlias("record.scores", "userScore", firestore.WithUnnestIndexField("attempt")).
		Execute(ctx)
	// [END firestore_unnest_nested]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func findNearestSyntaxExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_find_nearest_syntax]
	snapshot := client.Pipeline().Collection("cities").
		FindNearest("embedding", []float64{1.5, 2.345}, firestore.PipelineDistanceMeasureEuclidean).
		Execute(ctx)
	// [END firestore_find_nearest_syntax]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func findNearestLimitExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_find_nearest_limit]
	snapshot := client.Pipeline().Collection("cities").
		FindNearest("embedding", []float64{1.5, 2.345}, firestore.PipelineDistanceMeasureEuclidean, firestore.WithFindNearestLimit(10)).
		Execute(ctx)
	// [END firestore_find_nearest_limit]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func findNearestDistanceDataExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_find_nearest_distance_data]
	client.Collection("cities").Doc("SF").Set(ctx, map[string]any{
		"name":      "San Francisco",
		"embedding": []float64{1.0, -1.0},
	})
	client.Collection("cities").Doc("TO").Set(ctx, map[string]any{
		"name":      "Toronto",
		"embedding": []float64{5.0, -10.0},
	})
	client.Collection("cities").Doc("AT").Set(ctx, map[string]any{
		"name":      "Atlantis",
		"embedding": []float64{2.0, -4.0},
	})
	// [END firestore_find_nearest_distance_data]
	return nil
}

func findNearestDistanceExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_find_nearest_distance]
	snapshot := client.Pipeline().Collection("cities").
		FindNearest("embedding", []float64{1.3, 2.345}, firestore.PipelineDistanceMeasureEuclidean, firestore.WithFindNearestDistanceField("computedDistance")).
		Execute(ctx)
	// [END firestore_find_nearest_distance]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}
