// Copyright 2022 Google LLC
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

// [START bigqueryanalyticshub_quickstart]

// The bigquery_analyticshub_quickstart application demonstrates usage of the
// BigQuery analyticshub API.
package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"

	dataexchange "cloud.google.com/go/bigquery/dataexchange/apiv1beta1"
	dataexchangepb "google.golang.org/genproto/googleapis/cloud/bigquery/dataexchange/v1beta1"

	"google.golang.org/api/iterator"
)

func main() {

	// Define two command line flags for controlling the behavior of this quickstart.
	var (
		projectID = flag.String("project_id", "", "Cloud Project ID, used for session creation.")
		location  = flag.String("location", "US", "BigQuery location used for interactions.")
	)
	flag.Parse()
	if *projectID == "" {
		log.Fatal("empty --project_id specified, please provide a valid project ID")
	}

	ctx := context.Background()
	dataExchClient, err := dataexchange.NewAnalyticsHubClient(ctx)
	if err != nil {
		log.Fatalf("NewClient: %v", err)
	}
	defer dataExchClient.Close()

	s, err := reportDataExchanges(ctx, dataExchClient, *projectID, *location)
	if err != nil {
		log.Fatalf("listDataExchanges: %v", err)
	}
	fmt.Println(s)

}

func reportDataExchanges(ctx context.Context, client *dataexchange.AnalyticsHubClient, projectID, location string) (string, error) {
	var buf bytes.Buffer

	req := &dataexchangepb.ListDataExchangesRequest{
		Parent: fmt.Sprintf("projects/%s/location/%s", projectID, location),
	}

	it := client.ListDataExchanges(ctx, req)
	buf.WriteString("Data Exchanges:\n")
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return "", fmt.Errorf("error iterating results: %v", err)
		}
		buf.WriteString(fmt.Sprintf("Exchange %s has description: %s", resp.GetName(), resp.GetDescription()))
	}
	return buf.String(), nil
}
