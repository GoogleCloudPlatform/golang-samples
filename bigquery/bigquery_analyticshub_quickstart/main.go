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

// [START analyticshub_quickstart]

// The analyticshub quickstart application demonstrates usage of the
// Analytics hub API by creating an example data exchange and listing.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	analyticshub "cloud.google.com/go/bigquery/analyticshub/apiv1"
	"cloud.google.com/go/bigquery/analyticshub/apiv1/analyticshubpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {

	// Define the command line flags for controlling the behavior of this quickstart.
	var (
		projectID            = flag.String("project_id", "", "Cloud Project ID, used for session creation.")
		location             = flag.String("location", "US", "BigQuery location used for interactions.")
		exchangeID           = flag.String("exchange_id", "ExampleDataExchange", "identifier of the example data exchange")
		listingID            = flag.String("listing_id", "ExampleDataExchange", "identifier of the example data exchange")
		exampleDatasetSource = flag.String("dataset_source", "", "dataset source in the form projects/myproject/datasets/mydataset")
		delete               = flag.Bool("delete_exchange", true, "delete exchange at the end of quickstart")
	)
	flag.Parse()
	// Perform simple validation of the specified flags.
	if *projectID == "" {
		log.Fatal("empty --project_id specified, please provide a valid project ID")
	}
	if *exampleDatasetSource == "" {
		log.Fatalf("empty --dataset_source specified, please provide in the form \"projects/myproject/datasets/mydataset\"")
	}

	// Instantiate the client.
	ctx := context.Background()
	ahubClient, err := analyticshub.NewClient(ctx)
	if err != nil {
		log.Fatalf("NewClient: %v", err)
	}
	defer ahubClient.Close()

	// Then, create the data exchange (or return information about one already bearing the example name), and
	// print information about it.
	exchange, err := createOrGetDataExchange(ctx, ahubClient, *projectID, *location, *exchangeID)
	if err != nil {
		log.Fatalf("failed to get information about the exchange: %v", err)
	}
	fmt.Printf("\nData Exchange Information\n")
	fmt.Printf("Exchange Name: %s\n", exchange.GetName())
	if desc := exchange.GetDescription(); desc != "" {
		fmt.Printf("Exchange Description: %s", desc)
	}

	// Finally, create a listing within the data exchange and print information about it.
	listing, err := createListing(ctx, ahubClient, *projectID, *location, *exchangeID, *listingID, *exampleDatasetSource)
	if err != nil {
		log.Fatalf("failed to create the listing within the exchange: %v", err)
	}
	fmt.Printf("\n\nListing Information\n")
	fmt.Printf("Listing Name: %s\n", listing.GetName())
	if desc := listing.GetDescription(); desc != "" {
		fmt.Printf("Listing Description: %s\n", desc)
	}
	fmt.Printf("Listing State: %s\n", listing.GetState().String())
	if source := listing.GetSource(); source != nil {
		if dsSource, ok := source.(*analyticshubpb.Listing_BigqueryDataset); ok && dsSource.BigqueryDataset != nil {
			if dataset := dsSource.BigqueryDataset.GetDataset(); dataset != "" {
				fmt.Printf("Source is a bigquery dataset: %s", dataset)
			}
		}
	}
	// Optionally, delete the data exchange at the end of the quickstart to clean up the resources used.
	if *delete {
		fmt.Printf("\n\n")
		if err := deleteDataExchange(ctx, ahubClient, *projectID, *location, *exchangeID); err != nil {
			log.Fatalf("failed to delete exchange: %v", err)
		}
		fmt.Printf("Exchange projects/%s/locations/%s/dataExchanges/%s was deleted.\n", *projectID, *location, *exchangeID)
	}
	fmt.Printf("\nQuickstart completed.\n")
}

// createOrGetDataExchange creates an example data exchange, or returns information about the exchange already bearing
// the example identifier.
func createOrGetDataExchange(ctx context.Context, client *analyticshub.Client, projectID, location, exchangeID string) (*analyticshubpb.DataExchange, error) {
	req := &analyticshubpb.CreateDataExchangeRequest{
		Parent:         fmt.Sprintf("projects/%s/locations/%s", projectID, location),
		DataExchangeId: exchangeID,
		DataExchange: &analyticshubpb.DataExchange{
			DisplayName:    "Example Data Exchange",
			Description:    "Exchange created as part of an API quickstart",
			PrimaryContact: "",
			Documentation:  "https://link.to.optional.documentation/",
		},
	}

	resp, err := client.CreateDataExchange(ctx, req)
	if err != nil {
		// We'll handle one specific error case specially, the case of the exchange already existing.  In this instance,
		// we'll issue a second request to fetch the exchange information for the already present exchange and return it.
		if code := status.Code(err); code == codes.AlreadyExists {
			getReq := &analyticshubpb.GetDataExchangeRequest{
				Name: fmt.Sprintf("projects/%s/locations/%s/dataExchanges/%s", projectID, location, exchangeID),
			}
			resp, err = client.GetDataExchange(ctx, getReq)
			if err != nil {
				return nil, fmt.Errorf("error getting dataExchange: %w", err)
			}
			return resp, nil
		}
		// For all other cases, return the error from creation request.
		return nil, err
	}
	return resp, nil
}

// createListing creates an example listing within the specified exchange using the provided source dataset.
func createListing(ctx context.Context, client *analyticshub.Client, projectID, location, exchangeID, listingID, sourceDataset string) (*analyticshubpb.Listing, error) {
	req := &analyticshubpb.CreateListingRequest{
		Parent:    fmt.Sprintf("projects/%s/locations/%s/dataExchanges/%s", projectID, location, exchangeID),
		ListingId: listingID,
		Listing: &analyticshubpb.Listing{
			DisplayName: "Example Exchange Listing",
			Description: "Example listing created as part of an API quickstart",
			Categories: []analyticshubpb.Listing_Category{
				analyticshubpb.Listing_CATEGORY_OTHERS,
			},
			Source: &analyticshubpb.Listing_BigqueryDataset{
				BigqueryDataset: &analyticshubpb.Listing_BigQueryDatasetSource{
					Dataset: sourceDataset,
				},
			},
		},
	}
	return client.CreateListing(ctx, req)
}

// deleteDataExchange deletes a data exchange.
func deleteDataExchange(ctx context.Context, client *analyticshub.Client, projectID, location, exchangeID string) error {
	req := &analyticshubpb.DeleteDataExchangeRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/dataExchanges/%s", projectID, location, exchangeID),
	}
	return client.DeleteDataExchange(ctx, req)
}

// [END analyticshub_quickstart]
