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

import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
)

// unasignStaticAddresFromExistingVM removes a static external IP address from an existing VM instance in network interface.
func unassignStaticAddressFromExistingVM(w io.Writer, projectID, zone, instanceName, networkInterfaceName string) error {
	// projectID := "your_project_id"
	// zone := "europe-central2-b"
	// instanceName := "your_instance_name"
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

	fmt.Fprintf(w, "Static address %s unassigned from the instance %s\n", *accessConfig.NatIP, instanceName)

	return nil
}
