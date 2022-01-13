// Copyright 2021 Google LLC
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

package spanner

// [START spanner_opencensus_capture_query_stats_metric]

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

// OpenCensus Tag, Measure and View.
var (
	QueryStatsElapsed = stats.Float64("cloud.google.com/go/spanner/query_stats_elapsed",
		"The execution of the query", "ms")
	QueryStatsLatencyView = view.View{
		Name:        "cloud.google.com/go/spanner/query_stats_elapsed",
		Measure:     QueryStatsElapsed,
		Description: "The execution of the query",
		Aggregation: view.Distribution(0.0, 0.01, 0.05, 0.1, 0.3, 0.6, 0.8, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 8.0, 10.0, 13.0,
			16.0, 20.0, 25.0, 30.0, 40.0, 50.0, 65.0, 80.0, 100.0, 130.0, 160.0, 200.0, 250.0,
			300.0, 400.0, 500.0, 650.0, 800.0, 1000.0, 2000.0, 5000.0, 10000.0, 20000.0, 50000.0,
			100000.0),
		TagKeys: []tag.Key{}}
)

func queryWithQueryStats(w io.Writer, db string) error {
	projectID, _, _, err := parseDatabaseName(db)
	if err != nil {
		return err
	}

	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	// Register OpenCensus views.
	err = view.Register(&QueryStatsLatencyView)
	if err != nil {
		return err
	}

	// Create OpenCensus Stackdriver exporter.
	sd, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID: projectID,
	})
	if err != nil {
		return err
	}
	// It is imperative to invoke flush before your main function exits
	defer sd.Flush()

	// Start the metrics exporter
	sd.StartMetricsExporter()
	defer sd.StopMetricsExporter()

	// Execute a SQL query and get the query stats.
	stmt := spanner.Statement{SQL: `SELECT SingerId, AlbumId, AlbumTitle FROM Albums`}
	iter := client.Single().QueryWithStats(ctx, stmt)
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			// Record query execution time with OpenCensus.
			elapasedTime := iter.QueryStats["elapsed_time"].(string)
			elapasedTimeMs, err := strconv.ParseFloat(strings.TrimSuffix(elapasedTime, " msecs"), 64)
			if err != nil {
				return err
			}
			stats.Record(ctx, QueryStatsElapsed.M(elapasedTimeMs))
			return nil
		}
		if err != nil {
			return err
		}
		var singerID, albumID int64
		var albumTitle string
		if err := row.Columns(&singerID, &albumID, &albumTitle); err != nil {
			return err
		}
		fmt.Fprintf(w, "%d %d %s\n", singerID, albumID, albumTitle)
	}
}

// [END spanner_opencensus_capture_query_stats_metric]
