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

// [START compute_ip_address_promote_ephemeral]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/protobuf/proto"
)

// promoteEphemeralAddress promotes an ephemeral IP address to a reserved static external IP address.
func promoteEphemeralAddress(w io.Writer, projectID, region, ephemeralIP, addressName string) error {
	// projectID := "your_project_id"
	// region := "europe-central2"
	// ephemeral_ip := "214.123.100.121"
	ctx := context.Background()

	client, err := compute.NewAddressesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewAddressesRESTClient: %w", err)
	}
	defer client.Close()

	addressResource := &computepb.Address{
		Name:        proto.String(addressName),
		AddressType: proto.String("EXTERNAL"),
		Address:     proto.String(ephemeralIP),
	}

	req := &computepb.InsertAddressRequest{
		Project:         projectID,
		Region:          region,
		AddressResource: addressResource,
	}

	op, err := client.Insert(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to insert address promoted: %v", err)
	}

	// Wait for the operation to complete
	err = op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("failed to complete promotion operation: %v", err)
	}

	fmt.Fprintf(w, "Ephemeral IP %s address promoted successfully", ephemeralIP)
	return nil
}

// [END compute_ip_address_promote_ephemeral]
