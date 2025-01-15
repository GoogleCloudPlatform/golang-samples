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

// [START spanner_opencensus_capture_gfe_metric]

// We are in the process of adding support in the Cloud Spanner Go Client Library
// to capture the gfe_latency metric.

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	spanner "cloud.google.com/go/spanner/apiv1"
	sppb "cloud.google.com/go/spanner/apiv1/spannerpb"
	gax "github.com/googleapis/gax-go/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

// OpenCensus Tag, Measure and View.
var (
	KeyMethod    = tag.MustNewKey("grpc_client_method")
	GFELatencyMs = stats.Int64("cloud.google.com/go/spanner/gfe_latency",
		"Latency between Google's network receives an RPC and reads back the first byte of the response", "ms")
	GFELatencyView = view.View{
		Name:        "cloud.google.com/go/spanner/gfe_latency",
		Measure:     GFELatencyMs,
		Description: "Latency between Google's network receives an RPC and reads back the first byte of the response",
		Aggregation: view.Distribution(0.0, 0.01, 0.05, 0.1, 0.3, 0.6, 0.8, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 8.0, 10.0, 13.0,
			16.0, 20.0, 25.0, 30.0, 40.0, 50.0, 65.0, 80.0, 100.0, 130.0, 160.0, 200.0, 250.0,
			300.0, 400.0, 500.0, 650.0, 800.0, 1000.0, 2000.0, 5000.0, 10000.0, 20000.0, 50000.0,
			100000.0),
		TagKeys: []tag.Key{KeyMethod}}
)

func queryWithGFELatency(w io.Writer, db string) error {
	projectID, _, _, err := parseDatabaseName(db)
	if err != nil {
		return err
	}

	ctx := context.Background()
	client, err := spanner.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	// Register OpenCensus views.
	err = view.Register(&GFELatencyView)
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

	// Create a session.
	req := &sppb.CreateSessionRequest{Database: db}
	session, err := client.CreateSession(ctx, req)
	if err != nil {
		return err
	}

	// Execute a SQL query and retrieve the GFE server-timing header in gRPC metadata.
	req2 := &sppb.ExecuteSqlRequest{
		Session: session.Name,
		Sql:     `SELECT SingerId, AlbumId, AlbumTitle FROM Albums`,
	}
	var md metadata.MD
	resultSet, err := client.ExecuteSql(ctx, req2, gax.WithGRPCOptions(grpc.Header(&md)))
	if err != nil {
		return err
	}
	for _, row := range resultSet.GetRows() {
		for _, value := range row.GetValues() {
			fmt.Fprintf(w, "%s ", value.GetStringValue())
		}
		fmt.Fprintf(w, "\n")
	}

	// The format is: "server-timing: gfet4t7; dur=[GFE latency in ms]"
	srvTiming := md.Get("server-timing")[0]
	gfeLtcy, err := strconv.Atoi(strings.TrimPrefix(srvTiming, "gfet4t7; dur="))
	if err != nil {
		return err
	}
	// Record GFE t4t7 latency with OpenCensus.
	ctx, err = tag.New(ctx, tag.Insert(KeyMethod, "ExecuteSql"))
	if err != nil {
		return err
	}
	stats.Record(ctx, GFELatencyMs.M(int64(gfeLtcy)))

	return nil
}

// [END spanner_opencensus_capture_gfe_metric]
