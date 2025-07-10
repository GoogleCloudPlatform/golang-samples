// Copyright 2025 Google LLC
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

package client

// [START bigquery_client_usage]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/bigquery/v2/apiv2/bigquerypb"
	"cloud.google.com/go/bigquery/v2/apiv2_client"

	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

const datasetListLimit = 10
const tableListLimit = 10

// basicClientUsage demonstrates the basic steps of creating BigQuery clients, using them to invoke
// BQ RPCs, and shutting down.
func basicClientUsage(ctx context.Context, w io.Writer, useGRPC bool) error {

	// Client instantiation accepts options that affect the behavior of the
	// generated clients.  See https://pkg.go.dev/google.golang.org/api/option#ClientOption for
	// the set of defined options.
	//
	// While options exist for explicitly declaring how the client authenticates, in this
	// case we'll rely on Application Default Credentials (ADC), which resolves credentials implicitly
	// by searching a well defined set of credential sources.  For more information on ADC,
	// see https://cloud.google.com/docs/authentication/application-default-credentials
	//
	// For this example, we'll define a more benign option (WithUserAgent) that can be used to
	// help identify what particular piece of software generated the request by attaching the
	// additional information to outgoing requests.  Users of GCP features such as audit logging
	// will be able to see this additional information in audit log entries about specific API requests.
	opts := []option.ClientOption{
		option.WithUserAgent("my-golang-sample"),
	}

	// Each RPC service in BigQuery has a dedicated client, but for this example we'll
	// use our aggregated client which exposes access to all the individual services.
	var bqClient *apiv2_client.Client
	var err error

	// Create the client.  Generally, you should favor reuse of clients rather than recreating them frequently
	// as they require a bit of setup. It is also important to properly close a client when you're done
	// using it.  This example uses Go's defer functionality to close the clients when the current
	// function exits, since the usage here is more trivial.
	if useGRPC {
		// Create clients that use gRPC as the underlying transport protocol.
		bqClient, err = apiv2_client.NewClient(ctx, opts...)
		if err != nil {
			return fmt.Errorf("NewClient: %w", err)
		}
		defer bqClient.Close()
	} else {
		// Fallback to using HTTP REST as the underlying transport protocol.
		bqClient, err = apiv2_client.NewRESTClient(ctx, opts...)
		if err != nil {
			return fmt.Errorf("NewRESTClient: %w", err)
		}
		defer bqClient.Close()
	}
	// With the instantiated client, we can now call the various RPCs defined within the
	// BigQuery service.  We'll use our clients to get information about some of the resources
	// within BigQuery's public datasets.  More info about public datasets can be found at
	// https://cloud.google.com/bigquery/public-data

	req := &bigquerypb.ListDatasetsRequest{
		// This project contains many of the public datasets hosted within BigQuery.
		ProjectId: "bigquery-public-data",
	}

	// haveListedTables keeps track of whether we've listed tables at least once in this example.
	haveListedTables := false
	// dsCount and tblCount keep track of how many of each resource this sample has listed.
	var dsCount, tblCount int

	// Create a dataset iterator, used to process the list of datasets present in the designated project.
	it := bqClient.ListDatasets(ctx, req)
	for {
		dataset, err := it.Next()
		if err == iterator.Done {
			// We're reached the end of the iteration, break the loop.
			break
		}
		if err != nil {
			return fmt.Errorf("dataset iterator errored: %w", err)
		}
		dsCount = dsCount + 1
		if dsCount > datasetListLimit {
			break
		}
		// Log basic information about the encountered dataset.
		fmt.Fprintf(w, "ListDatasets: dataset %q in location %q\n",
			dataset.GetDatasetReference().GetDatasetId(),
			dataset.GetLocation())

		// We only demonstrate listing tables once in this example, but one could recurse through
		// the API to enumerate all tables in all encountered datasets.  However, this kind of traversal
		// may be better done through other mechanisms like INFORMATION_SCHEMA queries depending on the
		// specific needs.
		if !haveListedTables {
			// Now, we'll use table-related functionality, and list the tables
			// within the first dataset encountered.
			tblReq := &bigquerypb.ListTablesRequest{
				ProjectId: dataset.GetDatasetReference().GetProjectId(),
				DatasetId: dataset.DatasetReference.GetDatasetId(),
			}
			tblIt := bqClient.ListTables(ctx, tblReq)
			haveListedTables = true
			for {
				table, err := tblIt.Next()
				if err == iterator.Done {
					// We're reached the end of the iteration, break the loop.
					break
				}
				if err != nil {
					return fmt.Errorf("table iterator errored: %w", err)
				}
				tblCount = tblCount + 1
				if tblCount > tableListLimit {
					break
				}
				// Log basic information about the encountered table.
				fmt.Fprintf(w, "ListTables: table %q in dataset %q\n",
					table.GetTableReference().GetTableId(),
					table.GetTableReference().GetDatasetId(),
				)
			}
		}
	}
	return nil
}

// [END bigquery_client_usage]
