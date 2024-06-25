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

// [START compute_ip_address_release_static_external]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
)

// releaseRegionalStaticExternal releases a static regional external IP address.
func releaseRegionalStaticExternal(w io.Writer, projectID, region, addressName string) error {
	// projectID := "your_project_id"
	// region := "us-central1"
	// addressName := "your_address_name"

	ctx := context.Background()
	addressesClient, err := compute.NewAddressesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewAddressesRESTClient: %w", err)
	}
	defer addressesClient.Close()

	req := &computepb.DeleteAddressRequest{
		Project: projectID,
		Region:  region,
		Address: addressName,
	}

	op, err := addressesClient.Delete(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to release static external IP address: %w", err)
	}

	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %w", err)
	}

	fmt.Fprintf(w, "Static external IP address released\n")

	return nil
}

// releaseGlobalStaticExternal releases a static global external IP address.
func releaseGlobalStaticExternal(w io.Writer, projectID, addressName string) error {
	// projectID := "your_project_id"
	// addressName := "your_address_name"

	ctx := context.Background()
	addressesClient, err := compute.NewGlobalAddressesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewGlobalAddressesRESTClient: %w", err)
	}
	defer addressesClient.Close()

	req := &computepb.DeleteGlobalAddressRequest{
		Project: projectID,
		Address: addressName,
	}

	op, err := addressesClient.Delete(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to release static external IP address: %w", err)
	}

	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %w", err)
	}

	fmt.Fprintf(w, "Static external IP address released\n")

	return nil
}

// [END compute_ip_address_release_static_external]
