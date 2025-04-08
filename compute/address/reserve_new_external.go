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

// [START compute_ip_address_reserve_new_external]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
)

// reserveNewRegionalExternal reserves a new regional external IP address in Google Cloud Platform.
func reserveNewRegionalExternal(w io.Writer, projectID, region, addressName string, isPremium bool) (*computepb.Address, error) {
	// projectID := "your_project_id"
	// region := "europe-central2"
	// addressName := "your_address_name"
	// isPremium := true
	ctx := context.Background()

	networkTier := computepb.AccessConfig_STANDARD.String()
	if isPremium {
		networkTier = computepb.AccessConfig_PREMIUM.String()
	}

	address := &computepb.Address{
		Name:        &addressName,
		NetworkTier: &networkTier,
	}

	client, err := compute.NewAddressesRESTClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewAddressesRESTClient: %w", err)
	}
	defer client.Close()

	req := &computepb.InsertAddressRequest{
		Project:         projectID,
		Region:          region,
		AddressResource: address,
	}

	op, err := client.Insert(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("unable to reserve regional address: %w", err)
	}

	err = op.Wait(ctx)
	if err != nil {
		return nil, fmt.Errorf("waiting for the regional address reservation operation to complete: %w", err)
	}

	addressResult, err := client.Get(ctx, &computepb.GetAddressRequest{
		Project: projectID,
		Region:  region,
		Address: addressName,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to get reserved regional address: %w", err)
	}

	fmt.Fprintf(w, "Regional address %v reserved: %v", addressName, addressResult.GetAddress())

	return addressResult, err
}

// reserveNewGlobalExternal reserves a new global external IP address in Google Cloud Platform.
func reserveNewGlobalExternal(w io.Writer, projectID, addressName string, isV6 bool) (*computepb.Address, error) {
	// projectID := "your_project_id"
	// addressName := "your_address_name"
	// isV6 := false
	ctx := context.Background()
	ipVersion := computepb.Address_IPV4.String()
	if isV6 {
		ipVersion = computepb.Address_IPV6.String()
	}

	address := &computepb.Address{
		Name:      &addressName,
		IpVersion: &ipVersion,
	}

	client, err := compute.NewGlobalAddressesRESTClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewGlobalAddressesRESTClient: %w", err)
	}
	defer client.Close()

	req := &computepb.InsertGlobalAddressRequest{
		Project:         projectID,
		AddressResource: address,
	}

	op, err := client.Insert(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("unable to reserve global address: %w", err)
	}

	err = op.Wait(ctx)
	if err != nil {
		return nil, fmt.Errorf("waiting for the global address reservation operation to complete: %w", err)
	}

	addressResult, err := client.Get(ctx, &computepb.GetGlobalAddressRequest{
		Project: projectID,
		Address: addressName,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to get reserved global address: %w", err)
	}

	fmt.Fprintf(w, "Global address %v reserved: %v", addressName, addressResult.GetAddress())

	return addressResult, nil
}

// [END compute_ip_address_reserve_new_external]
