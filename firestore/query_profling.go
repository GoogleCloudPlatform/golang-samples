package firestore

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

func createChainedQuery(client *firestore.Client) {
	cities := client.Collection("cities")
	// [START firestore_query_filter_compound_multi_eq]
	denverQuery := cities.Where("name", "==", "Denver").Where("state", "==", "CO")
	caliQuery := cities.Where("state", "==", "CA").Where("population", "<=", 1000000)
	// [END firestore_query_filter_compound_multi_eq]

	_ = denverQuery
	_ = caliQuery
}

func SnippetQuery_RunQueryWithExplain(client *firestore.Client, w io.Writer) {
	ctx := context.Background()

	// [START firestore_query_explain_entity]
	// Build the query
	query := client.Collection("cities")

	// Set the explain options to get back *only* the plan summary
	it := query.WithRunOptions(firestore.ExplainOptions{}).Documents(ctx)

	// Query results will be empty
	fmt.Fprintln(w, "----- Query Results -----")
	for {
		city, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Fprintf(w, "Error fetching next city: %v", err)
			return
		}
		fmt.Fprintf(w, "City %q\n", city.Name)
	}

	// Get the explain metrics
	explainMetrics := it.ExplainMetrics

	planSummary := explainMetrics.PlanSummary
	fmt.Fprintln(w, "----- Indexes Used -----")
	for k, v := range planSummary.IndexesUsed {
		fmt.Fprintf(w, "%+v: %+v\n", k, v)
	}
	// [END firestore_query_explain_entity]
}

func SnippetQuery_RunQueryWithExplainAnalyze(client *firestore.Client, w io.Writer) {
	ctx := context.Background()

	// [START firestore_query_explain_analyze_entity]
	// Build the query
	query := client.Collection("cities")

	// Set explain options with analzye = true to get back the query stats, plan info, and query
	// results
	it := query.WithRunOptions(firestore.ExplainOptions{Analyze: true}).Documents(ctx)

	// Get the query results
	fmt.Fprintln(w, "----- Query Results -----")
	for {
		city, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Fprintf(w, "Error fetching next city: %v", err)
			return
		}
		fmt.Fprintf(w, "City %q\n", city.Name)
	}

	// Get plan summary
	planSummary := it.ExplainMetrics.PlanSummary
	fmt.Fprintln(w, "----- Indexes Used -----")
	for k, v := range planSummary.IndexesUsed {
		fmt.Fprintf(w, "%+v: %+v\n", k, v)
	}

	// Get the execution stats
	executionStats := it.ExplainMetrics.ExecutionStats
	fmt.Fprintln(w, "----- Execution Stats -----")
	fmt.Fprintf(w, "%+v\n", executionStats)
	fmt.Fprintln(w, "----- Debug Stats -----")
	for k, v := range *executionStats.DebugStats {
		fmt.Fprintf(w, "%+v: %+v\n", k, v)
	}
	// [END firestore_query_explain_analyze_entity]
}

func SnippetQuery_RunAggregationQueryWithExplain(client *firestore.Client, w io.Writer) {
	ctx := context.Background()

	// [START firestore_query_explain_aggregation]
	// Build the query
	query := client.Collection("cities")

	// Set the explain options to get back *only* the plan summary
	ar, err := query.NewAggregationQuery().WithCount("count").Get firestore.ExplainOptions{})

	// Get the explain metrics
	explainMetrics := ar.ExplainMetrics

	planSummary := explainMetrics.PlanSummary
	fmt.Fprintln(w, "----- Indexes Used -----")
	for k, v := range planSummary.IndexesUsed {
		fmt.Fprintf(w, "%+v: %+v\n", k, v)
	}
	// [END firestore_query_explain_aggregation]

	_ = err // Check non-nil errors
}

func SnippetQuery_RunAggregationQueryWithExplainAnalyze(client *firestore.Client, w io.Writer) {
	ctx := context.Background()

	// [START firestore_query_explain_analyze_aggregation]
	// Build the query
	query := client.Collection("cities")

	// Set explain options with analzye = true to get back the query stats, plan info, and query
	// results
	countAlias := "count"
	ar, err := client.RunAggregationQueryWithOptions(ctx,
		query.NewAggregationQuery().WithCount(countAlias), firestore.ExplainOptions{Analyze: true})

	// Get the query results
	fmt.Fprintln(w, "----- Query Results -----")
	result := ar.Result[countAlias]
	fmt.Fprintf(w, "Count %v\n", result)

	// Get plan summary
	planSummary := ar.ExplainMetrics.PlanSummary
	fmt.Fprintln(w, "----- Indexes Used -----")
	for k, v := range planSummary.IndexesUsed {
		fmt.Fprintf(w, "%+v: %+v\n", k, v)
	}

	// Get the execution stats
	executionStats := ar.ExplainMetrics.ExecutionStats
	fmt.Fprintln(w, "----- Execution Stats -----")
	fmt.Fprintf(w, "%+v\n", executionStats)
	fmt.Fprintln(w, "----- Debug Stats -----")
	for k, v := range *executionStats.DebugStats {
		fmt.Fprintf(w, "%+v: %+v\n", k, v)
	}
	// [END firestore_query_explain_analyze_aggregation]

	_ = err // Check non-nil errors
}
