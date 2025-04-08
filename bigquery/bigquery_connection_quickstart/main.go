// Copyright 2020 Google LLC
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

// [START bigqueryconnection_quickstart]

// The bigquery_connection_quickstart application demonstrates basic usage of the
// BigQuery connection API.
package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	connection "cloud.google.com/go/bigquery/connection/apiv1"
	"cloud.google.com/go/bigquery/connection/apiv1/connectionpb"
	"google.golang.org/api/iterator"
)

func main() {

	// Define two command line flags for controlling the behavior of this quickstart.
	projectID := flag.String("project_id", "", "Cloud Project ID, used for session creation.")
	location := flag.String("location", "US", "BigQuery location used for interactions.")

	// Parse flags and do some minimal validation.
	flag.Parse()
	if *projectID == "" {
		log.Fatal("empty --project_id specified, please provide a valid project ID")
	}
	if *location == "" {
		log.Fatal("empty --location specified, please provide a valid location")
	}

	ctx := context.Background()
	connClient, err := connection.NewClient(ctx)
	if err != nil {
		log.Fatalf("NewClient: %v", err)
	}
	defer connClient.Close()

	s, err := reportConnections(ctx, connClient, *projectID, *location)
	if err != nil {
		log.Fatalf("printCapacityCommitments: %v", err)
	}
	fmt.Println(s)
}

// reportConnections gathers basic information about existing connections in a given project and location.
func reportConnections(ctx context.Context, client *connection.Client, projectID, location string) (string, error) {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Current connections defined in project %s in location %s:\n", projectID, location)

	req := &connectionpb.ListConnectionsRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, location),
	}
	totalConnections := 0
	it := client.ListConnections(ctx, req)
	for {
		conn, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return "", err
		}
		fmt.Fprintf(&buf, "\tConnection %s was created %s\n", conn.GetName(), unixMillisToTime(conn.GetCreationTime()).Format(time.RFC822Z))
		totalConnections++
	}
	fmt.Fprintf(&buf, "\n%d connections processed.\n", totalConnections)
	return buf.String(), nil
}

// unixMillisToTime converts epoch-millisecond representations used by the API into a time.Time representation.
func unixMillisToTime(m int64) time.Time {
	if m == 0 {
		return time.Time{}
	}
	return time.Unix(0, m*1e6)
}

// [END bigqueryconnection_quickstart]
