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

func pipelineConcepts(w io.Writer, client *firestore.Client) error {
	// [START pipeline_concepts]
	pipeline := client.Pipeline().
		Collection("cities").
		Where(firestore.FieldOf("population").GreaterThan(100000)).
		Sort(firestore.Orders(firestore.Ascending(firestore.FieldOf("name")))).
		Limit(10)
	// [END pipeline_concepts]
	fmt.Fprintln(w, pipeline)
	return nil
}

func basicRead(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START basic_read]
	pipeline := client.Pipeline().Collection("users")
	snapshot := pipeline.Execute(ctx)
	results, err := snapshot.Results().GetAll()
	if err != nil {
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
			return err
		}
		fmt.Fprintf(w, "%s => %v\n", result.Ref().ID, result.Data())
	}
	// [END basic_read]
	return nil
}

func pipelineInitialization(w io.Writer, projectID string) error {
	ctx := context.Background()
	// [START pipeline_initialization]
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return err
	}
	pipeline := client.Pipeline().Collection("books")
	// [END pipeline_initialization]
	fmt.Fprintln(w, pipeline)
	return nil
}

func fieldVsConstants(w io.Writer, client *firestore.Client) error {
	// [START field_or_constant]
	pipeline := client.Pipeline().Collection("cities").
		Where(firestore.FieldOf("name").Equal(firestore.ConstantOf("Toronto")))
	// [END field_or_constant]
	fmt.Fprintln(w, pipeline)
	return nil
}

func inputStages(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START input_stages]
	// Return all restaurants in San Francisco
	results1, err := client.Pipeline().Collection("cities/sf/restaurants").Execute(ctx).Results().GetAll()
	if err != nil {
		return err
	}

	// Return all restaurants
	results2, err := client.Pipeline().CollectionGroup("restaurants").Execute(ctx).Results().GetAll()
	if err != nil {
		return err
	}

	// Return all documents across all collections in the database (the entire database)
	results3, err := client.Pipeline().Database().Execute(ctx).Results().GetAll()
	if err != nil {
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
		return err
	}
	// [END input_stages]
	fmt.Fprintln(w, results1)
	fmt.Fprintln(w, results2)
	fmt.Fprintln(w, results3)
	fmt.Fprintln(w, results4)
	return nil
}
