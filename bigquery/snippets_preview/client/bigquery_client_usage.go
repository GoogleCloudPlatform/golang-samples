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

package dataset

// [START bigquery_client_usage]
import (
	"context"
	"fmt"
	"log"

	bigquery "cloud.google.com/go/bigquery/apiv2"
	"cloud.google.com/go/bigquery/apiv2/bigquerypb"

	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// basicClientUsage demonstrates the basic steps of creating BigQuery clients, using them to invoke
// BQ RPCs, and shutting down.
func basicClientUsage(useGRPC bool) error {
	ctx := context.Background()

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

	// Each RPC service in BigQuery has a dedicated client, so this example will create several of the RPC clients
	// using the shared set of options.  We'll also

	var datasetClient *bigquery.DatasetClient
	var tableClient *bigquery.TableClient
	var err error

	// Create the clients.  Generally, you should favor reuse of clients rather than recreating them frequently
	// as they require a bit of setup. It is also important to properly close a client when you're done
	// using it.  This example uses Go's defer functionality to close the clients when the current
	// function exits, since the usage here is more trivial.
	if useGRPC {
		// Create clients that use gRPC as the underlying transport protocol.
		datasetClient, err = bigquery.NewDatasetClient(ctx, opts...)
		if err != nil {
			return fmt.Errorf("NewDatasetClient: %w", err)
		}
		defer datasetClient.Close()
		tableClient, err = bigquery.NewTableClient(ctx, opts...)
		if err != nil {
			return fmt.Errorf("NewDatasetClient: %w", err)
		}
		defer tableClient.Close()
	} else {
		// Fallback to using HTTP REST as the underlying transport protocol.
		datasetClient, err = bigquery.NewDatasetRESTClient(ctx, opts...)
		if err != nil {
			return fmt.Errorf("NewDatasetClient: %w", err)
		}
		defer datasetClient.Close()

		tableClient, err = bigquery.NewTableRESTClient(ctx, opts...)
		if err != nil {
			return fmt.Errorf("NewDatasetClient: %w", err)
		}
		defer tableClient.Close()
	}
	// With the instantiated clients, we can now call the various RPCs defined within the
	// BigQuery service.  We'll use our clients to get information about some of the resources
	// within BigQuery's public datasets.  More info about public datasets can be found at
	// https://cloud.google.com/bigquery/public-data

	req := &bigquerypb.ListDatasetsRequest{
		// This project contains many of the public datasets.
		ProjectId: "bigquery-public-data",
	}

	// haveListedTables keeps track of whether we've listed tables at least once in this example.
	haveListedTables := false

	// Create a dataset iterator, used too process the list of datasets present in the designated project.
	it := datasetClient.ListDatasets(ctx, req)
	for {
		dataset, err := it.Next()
		if err == iterator.Done {
			// We're reached the end of the iteration, break the loop.
			break
		}
		if err != nil {
			return fmt.Errorf("dataset iterator errored: %w", err)
		}
		// Log basic information about the encountered dataset.
		log.Printf("ListDatasets: dataset %q in location %q\n",
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
			tblIt := tableClient.ListTables(ctx, tblReq)
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
				// Log basic information about the encountered table.
				log.Printf("ListTables: table %q in dataset %q\n",
					table.GetTableReference().GetDatasetId(),
					table.GetTableReference().GetTableId())
			}
		}
	}
	return nil

}

// [END bigquery_client_usage]
