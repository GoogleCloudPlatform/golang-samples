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

// [START compute_ip_address_get_static_address]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
)

// getExternalAddress retrieves the external IP address of the given address.
func getRegionalExternal(w io.Writer, projectID, region, addressName string) (*computepb.Address, error) {
	// projectID := "your_project_id"
	// region := "europe-west3"
	// addressName := "your_address_name"

	ctx := context.Background()
	addressesClient, err := compute.NewAddressesRESTClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewAddressesRESTClient: %w", err)
	}
	defer addressesClient.Close()

	req := &computepb.GetAddressRequest{
		Project: projectID,
		Region:  region,
		Address: addressName,
	}

	address, err := addressesClient.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("unable to get address: %w", err)
	}

	fmt.Fprintf(w, "Regional address %s has external IP address: %s\n", addressName, address.GetAddress())

	return address, nil
}

func getGlobalExternal(w io.Writer, projectID, addressName string) (*computepb.Address, error) {

	ctx := context.Background()
	globalAddressesClient, err := compute.NewGlobalAddressesRESTClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewGlobalAddressesRESTClient: %w", err)
	}
	defer globalAddressesClient.Close()

	req := &computepb.GetGlobalAddressRequest{
		Project: projectID,
		Address: addressName,
	}

	address, err := globalAddressesClient.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("unable to get address: %w", err)
	}

	fmt.Fprintf(w, "Global address %s has external IP address: %s\n", addressName, address.GetAddress())

	return address, nil
}

// [END compute_ip_address_get_static_address]
