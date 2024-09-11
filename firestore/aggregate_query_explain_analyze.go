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

// [START firestore_query_explain_analyze_aggregation]

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/firestore"
)

func aggregationQueryExplainAnalyze(w io.Writer, projectID string) error {
	ctx := context.Background()

	// Create client
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("firestore.NewClient: %w", err)
	}
	defer client.Close()

	// Set explain options with analzye = true to get back the
	// query stats, plan info, and query results
	query := client.Collection("cities").WithRunOptions(firestore.ExplainOptions{Analyze: true})

	countAlias := "count"
	ar, err := query.NewAggregationQuery().WithCount(countAlias).GetResponse(ctx)
	if err != nil {
		fmt.Fprintf(w, "Error running aggregation query: %v", err)
		return err
	}

	// Get query results
	fmt.Fprintln(w, "----- Query Results -----")
	result := ar.Result[countAlias]
	fmt.Fprintf(w, "Count %v\n", result)

	// Get plan summary
	planSummary := ar.ExplainMetrics.PlanSummary
	fmt.Fprintln(w, "----- Indexes Used -----")
	for k, v := range planSummary.IndexesUsed {
		fmt.Fprintf(w, "%+v: %+v\n", k, v)
	}

	// Get execution stats
	executionStats := ar.ExplainMetrics.ExecutionStats
	fmt.Fprintln(w, "----- Execution Stats -----")
	fmt.Fprintf(w, "%+v\n", executionStats)
	fmt.Fprintln(w, "----- Debug Stats -----")
	for k, v := range *executionStats.DebugStats {
		fmt.Fprintf(w, "%+v: %+v\n", k, v)
	}
	return err
}

// [END firestore_query_explain_analyze_aggregation]
