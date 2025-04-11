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

// [START compute_ip_address_assign_static_existing_vm]
import (
	"context"
	"fmt"
	"io"

	"google.golang.org/protobuf/proto"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
)

// assignStaticAddressToExistingVM assigns a static external IP address to an existing VM instance.
// Note: VM and assigned IP must be in the same region.
func assignStaticAddressToExistingVM(w io.Writer, projectID, zone, instanceName, IPAddress, networkInterfaceName string) error {
	// projectID := "your_project_id"
	// zone := "europe-central2-b"
	// instanceName := "your_instance_name"
	// IPAddress := "34.111.222.333"
	// networkInterfaceName := "nic0"

	ctx := context.Background()
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewInstancesRESTClient: %w", err)
	}
	defer instancesClient.Close()

	reqGet := &computepb.GetInstanceRequest{
		Project:  projectID,
		Zone:     zone,
		Instance: instanceName,
	}

	instance, err := instancesClient.Get(ctx, reqGet)
	if err != nil {
		return fmt.Errorf("could not get instance: %w", err)
	}

	var networkInterface *computepb.NetworkInterface
	for _, ni := range instance.NetworkInterfaces {
		if *ni.Name == networkInterfaceName {
			networkInterface = ni
			break
		}
	}

	if networkInterface == nil {
		return fmt.Errorf("No network interface named '%s' found on instance %s", networkInterfaceName, instanceName)
	}

	var accessConfig *computepb.AccessConfig
	for _, ac := range networkInterface.AccessConfigs {
		if *ac.Type == computepb.AccessConfig_ONE_TO_ONE_NAT.String() {
			accessConfig = ac
			break
		}
	}

	if accessConfig != nil {
		// network interface is immutable - deletion stage is required in case of any assigned ip (static or ephemeral).
		reqDelete := &computepb.DeleteAccessConfigInstanceRequest{
			Project:          projectID,
			Zone:             zone,
			Instance:         instanceName,
			AccessConfig:     *accessConfig.Name,
			NetworkInterface: networkInterfaceName,
		}

		opDelete, err := instancesClient.DeleteAccessConfig(ctx, reqDelete)
		if err != nil {
			return fmt.Errorf("unable to delete access config: %w", err)
		}

		if err = opDelete.Wait(ctx); err != nil {
			return fmt.Errorf("unable to wait for the operation: %w", err)
		}
	}

	reqAdd := &computepb.AddAccessConfigInstanceRequest{
		Project:  projectID,
		Zone:     zone,
		Instance: instanceName,
		AccessConfigResource: &computepb.AccessConfig{
			NatIP: &IPAddress,
			Type:  proto.String(computepb.AccessConfig_ONE_TO_ONE_NAT.String()),
		},
		NetworkInterface: networkInterfaceName,
	}

	opAdd, err := instancesClient.AddAccessConfig(ctx, reqAdd)
	if err != nil {
		return fmt.Errorf("unable to add access config: %w", err)
	}

	if err = opAdd.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %w", err)
	}

	fmt.Fprintf(w, "Static address %s assigned to the instance %s\n", IPAddress, instanceName)

	return nil
}

// [END compute_ip_address_assign_static_existing_vm]
