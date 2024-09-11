// Copyright 2024 Google LLC
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

// [START firestore_query_explain_aggregation]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/firestore"
)

func aggregationQueryExplain(w io.Writer, projectID string) error {
	ctx := context.Background()

	// Create client
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("firestore.NewClient: %w", err)
	}
	defer client.Close()

	// Set the explain options to get back *only* the plan summary
	query := client.Collection("cities").WithRunOptions(firestore.ExplainOptions{})

	aggResult, err := query.NewAggregationQuery().WithCount("count").GetResponse(ctx)
	if err != nil {
		fmt.Fprintf(w, "Error running aggregation query: %v", err)
		return err
	}

	// Get explain metrics
	explainMetrics := aggResult.ExplainMetrics

	// Get plan summary
	planSummary := explainMetrics.PlanSummary
	fmt.Fprintln(w, "----- Indexes Used -----")
	for k, v := range planSummary.IndexesUsed {
		fmt.Fprintf(w, "%+v: %+v\n", k, v)
	}

	return nil
}

// [END firestore_query_explain_aggregation]
