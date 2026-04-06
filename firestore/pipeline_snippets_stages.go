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

func creatingIndexes(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START query_example]
	snapshot := client.Pipeline().
		Collection("books").
		Where(firestore.FieldOf("published").LessThan(1900)).
		Where(firestore.FieldOf("genre").Equal("Science Fiction")).
		Where(firestore.FieldOf("rating").GreaterThan(4.3)).
		Sort(firestore.Orders(firestore.Descending(firestore.FieldOf("published")))).
		Execute(ctx)
	// [END query_example]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func sparseIndexes(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START sparse_index_example]
	snapshot := client.Pipeline().
		Collection("books").
		Where(firestore.FieldOf("category").Like("%fantasy%")).
		Execute(ctx)
	// [END sparse_index_example]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func sparseIndexes2(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START sparse_index_example_2]
	snapshot := client.Pipeline().
		Collection("books").
		Sort(firestore.Orders(firestore.Ascending(firestore.FieldOf("release_date")))).
		Execute(ctx)
	// [END sparse_index_example_2]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func coveredQuery(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START covered_query]
	snapshot := client.Pipeline().
		Collection("books").
		Where(firestore.FieldOf("category").Like("%fantasy%")).
		Where(firestore.FieldOf("title").FieldExists()).
		Where(firestore.FieldOf("author").FieldExists()).
		Select(firestore.Fields("title", "author")).
		Execute(ctx)
	// [END covered_query]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func pagination(w io.Writer, client *firestore.Client) error {
	// [START pagination_not_supported_preview]
	// Existing pagination via `StartAt()`
	query := client.Collection("cities").OrderBy("population", firestore.Asc).StartAt(1000000)

	pipeline := client.Pipeline().
		Collection("cities").
		Where(firestore.FieldOf("population").GreaterThanOrEqual(1000000)).
		Sort(firestore.Orders(firestore.Descending(firestore.FieldOf("population"))))
	// [END pagination_not_supported_preview]
	fmt.Fprintln(w, query)
	fmt.Fprintln(w, pipeline)
	return nil
}

func collectionStage(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START collection_example]
	snapshot := client.Pipeline().
		Collection("users/bob/games").
		Sort(firestore.Orders(firestore.Ascending(firestore.FieldOf("name")))).
		Execute(ctx)
	// [END collection_example]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func collectionGroupStage(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START collection_group_example]
	snapshot := client.Pipeline().
		CollectionGroup("games").
		Sort(firestore.Orders(firestore.Ascending(firestore.FieldOf("name")))).
		Execute(ctx)
	// [END collection_group_example]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func databaseStage(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START database_example]
	// Count all documents in the database
	snapshot := client.Pipeline().
		Database().
		Aggregate(firestore.Accumulators(firestore.CountAll().As("total"))).
		Execute(ctx)
	// [END database_example]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func documentsStage(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START documents_example]
	snapshot := client.Pipeline().
		Documents([]*firestore.DocumentRef{
			client.Collection("cities").Doc("SF"),
			client.Collection("cities").Doc("DC"),
			client.Collection("cities").Doc("NY"),
		}).
		Execute(ctx)
	// [END documents_example]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func replaceWithStage(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START initial_data]
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
	// [END initial_data]

	// [START full_replace]
	snapshot := client.Pipeline().
		Collection("cities").
		ReplaceWith(firestore.FieldOf("location")).
		Execute(ctx)
	// [END full_replace]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}

	// [START map_merge_overwrite]
	// unsupported in client SDKs for now
	// [END map_merge_overwrite]
	fmt.Fprintln(w, results)
	return nil
}

func sampleStage(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START sample_example]
	// Get a sample of 100 documents in a database
	results1, err := client.Pipeline().Database().Sample(firestore.WithDocLimit(100)).Execute(ctx).Results().GetAll()
	if err != nil {
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
		return err
	}
	// [END sample_example]
	fmt.Fprintln(w, results1)
	fmt.Fprintln(w, results2)
	return nil
}

func samplePercent(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START sample_percent]
	// Get a sample of on average 50% of the documents in the database
	snapshot := client.Pipeline().
		Database().
		Sample(firestore.WithPercentage(0.5)).
		Execute(ctx)
	// [END sample_percent]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func sampleSyntaxExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START sample_syntax]
	sampled1, err := client.Pipeline().Database().Sample(firestore.WithDocLimit(50)).Execute(ctx).Results().GetAll()
	if err != nil {
		return err
	}

	sampled2, err := client.Pipeline().Database().Sample(firestore.WithPercentage(0.5)).Execute(ctx).Results().GetAll()
	if err != nil {
		return err
	}
	// [END sample_syntax]
	fmt.Fprintln(w, sampled1)
	fmt.Fprintln(w, sampled2)
	return nil
}

func sampleDocumentsDataExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START sample_documents_data]
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
	// [END sample_documents_data]
	return nil
}

func sampleDocumentsExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START sample_documents]
	snapshot := client.Pipeline().Collection("cities").Sample(firestore.WithDocLimit(1)).Execute(ctx)
	// [END sample_documents]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func sampleAllDocumentsExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START sample_all_documents]
	snapshot := client.Pipeline().Collection("cities").Sample(firestore.WithDocLimit(5)).Execute(ctx)
	// [END sample_all_documents]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func samplePercentageDataExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START sample_percentage_data]
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
	// [END sample_percentage_data]
	return nil
}

func samplePercentageExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START sample_percentage]
	snapshot := client.Pipeline().Collection("cities").Sample(firestore.WithPercentage(0.5)).Execute(ctx)
	// [END sample_percentage]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func unionStage(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START union_stage]
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
	// [END union_stage]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func unionStageStable(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START union_stage_stable]
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
	// [END union_stage_stable]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func offsetSyntaxExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START offset_syntax]
	snapshot := client.Pipeline().Collection("cities").Offset(10).Execute(ctx)
	// [END offset_syntax]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func collectionInputSyntaxExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START collection_input_syntax]
	snapshot := client.Pipeline().Collection("cities/SF/departments").Execute(ctx)
	// [END collection_input_syntax]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func collectionInputExampleData(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START collection_input_data]
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
	// [END collection_input_data]
	return nil
}

func collectionInputExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START collection_input]
	snapshot := client.Pipeline().Collection("cities").
		Sort(firestore.Orders(firestore.Ascending(firestore.FieldOf("name")))).
		Execute(ctx)
	// [END collection_input]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func subcollectionInputExampleData(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START subcollection_input_data]
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
	// [END subcollection_input_data]
	return nil
}

func subcollectionInputExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START subcollection_input]
	snapshot := client.Pipeline().Collection("cities/NY/departments").
		Sort(firestore.Orders(firestore.Ascending(firestore.FieldOf("employees")))).
		Execute(ctx)
	// [END subcollection_input]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func collectionGroupInputSyntaxExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START collection_group_input_syntax]
	snapshot := client.Pipeline().CollectionGroup("departments").Execute(ctx)
	// [END collection_group_input_syntax]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func collectionGroupInputExampleData(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START collection_group_data]
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
	// [END collection_group_data]
	return nil
}

func collectionGroupInputExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START collection_group_input]
	snapshot := client.Pipeline().CollectionGroup("departments").
		Sort(firestore.Orders(firestore.Ascending(firestore.FieldOf("employees")))).
		Execute(ctx)
	// [END collection_group_input]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func databaseInputSyntaxExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START database_syntax]
	snapshot := client.Pipeline().Database().Execute(ctx)
	// [END database_syntax]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func databaseInputSyntaxExampleData(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START database_input_data]
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
	// [END database_input_data]
	return nil
}

func databaseInputExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START database_input]
	snapshot := client.Pipeline().Database().
		Sort(firestore.Orders(firestore.Ascending(firestore.FieldOf("population")))).
		Execute(ctx)
	// [END database_input]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func documentInputSyntaxExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START document_input_syntax]
	snapshot := client.Pipeline().
		Documents([]*firestore.DocumentRef{
			client.Collection("cities").Doc("SF"),
			client.Collection("cities").Doc("NY"),
		}).
		Execute(ctx)
	// [END document_input_syntax]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func documentInputExampleData(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START document_input_data]
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
	// [END document_input_data]
	return nil
}

func documentInputExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START document_input]
	snapshot := client.Pipeline().
		Documents([]*firestore.DocumentRef{
			client.Collection("cities").Doc("SF"),
			client.Collection("cities").Doc("NYC"),
		}).
		Sort(firestore.Orders(firestore.Ascending(firestore.FieldOf("name")))).
		Execute(ctx)
	// [END document_input]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func limitSyntaxExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START limit_syntax]
	snapshot := client.Pipeline().Collection("cities").Limit(10).Execute(ctx)
	// [END limit_syntax]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func unionSyntaxExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START union_syntax]
	snapshot := client.Pipeline().
		Collection("cities/SF/restaurants").
		Union(client.Pipeline().Collection("cities/NYC/restaurants")).
		Execute(ctx)
	// [END union_syntax]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}
