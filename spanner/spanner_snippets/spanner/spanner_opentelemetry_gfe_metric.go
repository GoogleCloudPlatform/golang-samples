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

// [START spanner_opentelemetry_gfe_metric]

import (
	"context"
	"fmt"
	"io"
	"log"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
)

func queryWithGFELatencyMetric(w io.Writer, db string) error {
	ctx := context.Background()

	// Enable OpenTelemetry metrics for Spanner GFE metrics.
	spanner.EnableOpenTelemetryMetrics()

	// Create a new resource to uniquely identifies the application
	res, err := newResource()
	if err != nil {
		log.Fatal(err)
	}

	// Create a new meter provider
	meterProvider := getOtlpMeterProvider(ctx, res)

	client, err := spanner.NewClientWithConfig(ctx, db, spanner.ClientConfig{OpenTelemetryMeterProvider: meterProvider})
	if err != nil {
		return err
	}
	defer client.Close()

	stmt := spanner.Statement{SQL: `SELECT SingerId, AlbumId, AlbumTitle FROM Albums`}
	iter := client.Single().Query(ctx, stmt)
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done {
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

// [END spanner_opentelemetry_gfe_metric]
