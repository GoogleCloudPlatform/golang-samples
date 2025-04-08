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

package snippets

// [START compute_ip_address_list_static_external]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	"google.golang.org/api/iterator"

	"cloud.google.com/go/compute/apiv1/computepb"
)

// listRegionalExternal retrieves list external IP addresses in Google Cloud Platform region.
func listRegionalExternal(w io.Writer, projectID, region string) ([]*computepb.Address, error) {
	// projectID := "your_project_id"
	// region := "europe-west3"

	ctx := context.Background()
	// Create the service client.
	addressesClient, err := compute.NewAddressesRESTClient(ctx)
	if err != nil {
		return nil, err
	}
	defer addressesClient.Close()

	// Build the request.
	req := &computepb.ListAddressesRequest{
		Project: projectID,
		Region:  region,
	}

	// List the addresses.
	it := addressesClient.List(ctx, req)

	// Iterate over the results.
	var addresses []*computepb.Address
	for {
		address, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, address)
	}

	// Print the addresses.
	fmt.Fprint(w, "Fetched addresses: \n")
	for _, address := range addresses {
		fmt.Fprintf(w, "%s\n", *address.Name)
	}

	return addresses, nil
}

// listGlobalExternal retrieves list external global IP addresses in Google Cloud Platform.
func listGlobalExternal(w io.Writer, projectID string) ([]*computepb.Address, error) {
	// projectID := "your_project_id"

	ctx := context.Background()
	// Create the service client.
	addressesClient, err := compute.NewGlobalAddressesRESTClient(ctx)
	if err != nil {
		return nil, err
	}
	defer addressesClient.Close()

	// Build the request.
	req := &computepb.ListGlobalAddressesRequest{
		Project: projectID,
	}

	// List the addresses.
	it := addressesClient.List(ctx, req)

	// Iterate over the results.
	var addresses []*computepb.Address
	for {
		address, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, address)
	}

	// Print the addresses.
	fmt.Fprint(w, "Fetched addresses: \n")
	for _, address := range addresses {
		fmt.Fprintf(w, "%s\n", *address.Name)
	}

	return addresses, nil
}

// [END compute_ip_address_list_static_external]
