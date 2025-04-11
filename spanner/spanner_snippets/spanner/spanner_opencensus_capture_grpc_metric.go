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

// [START spanner_opencensus_capture_grpc_metric]

import (
	"context"
	"fmt"
	"io"
	"regexp"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats/view"
)

var validDatabasePattern = regexp.MustCompile("^projects/(?P<project>[^/]+)/instances/(?P<instance>[^/]+)/databases/(?P<database>[^/]+)$")

func queryWithGRPCMetric(w io.Writer, db string) error {
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
	if err := view.Register(ocgrpc.DefaultClientViews...); err != nil {
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
}

func parseDatabaseName(databaseUri string) (project, instance, database string, err error) {
	matches := validDatabasePattern.FindStringSubmatch(databaseUri)
	if len(matches) == 0 {
		return "", "", "", fmt.Errorf("failed to parse database name from %q according to pattern %q",
			databaseUri, validDatabasePattern.String())
	}
	return matches[1], matches[2], matches[3], nil
}

// [END spanner_opencensus_capture_grpc_metric]
