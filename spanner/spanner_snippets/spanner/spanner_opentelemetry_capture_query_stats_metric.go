//go:build go1.20
// +build go1.20

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

package spanner

// [START spanner_opentelemetry_capture_query_stats_metric]

import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"cloud.google.com/go/spanner"
	"go.opentelemetry.io/otel/metric"
	"google.golang.org/api/iterator"
)

func queryWithQueryStatsMetric(w io.Writer, db string) error {
	// db = `projects/<project>/instances/<instance-id>/database/<database-id>`
	ctx := context.Background()

	// Create a new resource to uniquely identify the application
	res, err := newResource()
	if err != nil {
		log.Fatal(err)
	}

	// Enable OpenTelemetry metrics before injecting meter provider.
	spanner.EnableOpenTelemetryMetrics()

	// Create a new meter provider
	meterProvider := getOtlpMeterProvider(ctx, res)

	queryStats := registerQueryStatsMetric(meterProvider)

	// Inject meter provider locally via ClientConfig when creating a spanner client.
	client, err := spanner.NewClientWithConfig(ctx, db, spanner.ClientConfig{OpenTelemetryMeterProvider: meterProvider})
	if err != nil {
		return err
	}
	defer client.Close()

	stmt := spanner.Statement{SQL: `SELECT SingerId, AlbumId, AlbumTitle FROM Albums`}
	iter := client.Single().QueryWithStats(ctx, stmt)
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			// Record query execution time with OpenTelemetry.
			elapasedTime := iter.QueryStats["elapsed_time"].(string)
			elapasedTimeMs, err := strconv.ParseFloat(strings.TrimSuffix(elapasedTime, " msecs"), 64)
			if err != nil {
				return err
			}
			queryStats.Record(ctx, elapasedTimeMs)
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

	meterProvider.ForceFlush(ctx)
	return nil
}

func registerQueryStatsMetric(mp metric.MeterProvider) metric.Float64Histogram {
	meter := mp.Meter(spanner.OtInstrumentationScope)
	queryStatsLatencyInstrument, err := meter.Float64Histogram(
		"spanner/query_stats_elapsed",
		metric.WithDescription("The execution of the query"),
		metric.WithUnit("ms"),
		metric.WithExplicitBucketBoundaries(0.0, 0.01, 0.05, 0.1, 0.3, 0.6, 0.8, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 8.0, 10.0, 13.0,
			16.0, 20.0, 25.0, 30.0, 40.0, 50.0, 65.0, 80.0, 100.0, 130.0, 160.0, 200.0, 250.0,
			300.0, 400.0, 500.0, 650.0, 800.0, 1000.0, 2000.0, 5000.0, 10000.0, 20000.0, 50000.0,
			100000.0),
	)
	if err != nil {
		fmt.Print(err)
	}
	return queryStatsLatencyInstrument
}

// [END spanner_opentelemetry_capture_query_stats_metric]
