// Copyright 2023 Google LLC
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

// [START firestore_count_query]
package firestore

import (
	"context"
	"errors"
	"fmt"
	"io"

	"cloud.google.com/go/firestore"
	firestorepb "cloud.google.com/go/firestore/apiv1/firestorepb"
)

func createCountQuery(w io.Writer, projectID string) error {

	// Instantiate the client
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return err
	}
	defer client.Close()

	collection := client.Collection("users")
	query := collection.Where("born", ">", 1850)

	// `alias` argument--"all"--provides a key for accessing the aggregate query
	// results. The alias value must be unique across all aggregation aliases in
	// an aggregation query and must conform to allowed Document field names.
	//
	// See https://cloud.google.com/firestore/docs/reference/rpc/google.firestore.v1#document for details.
	aggregationQuery := query.NewAggregationQuery().WithCount("all")
	results, err := aggregationQuery.Get(ctx)
	if err != nil {
		return err
	}

	count, ok := results["all"]
	if !ok {
		return errors.New("firestore: couldn't get alias for COUNT from results")
	}

	countValue := count.(*firestorepb.Value)
	fmt.Fprintf(w, "Number of results from query: %d\n", countValue.GetIntegerValue())
	return nil
}

// [END firestore_count_query]
