// Copyright 2023 Google LLC
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

// [START compute_instance_machine_type_update]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
)

// changeMachineType changes the machine type of the instance.
// If the instance is running, it will shut it down
// so that the machine type could be changed.
func changeMachineType(
	w io.Writer,
	projectID, zone, instanceName, newInstanceType string,
) error {
	// projectID := "your_project_id"
	// zone := "europe-central2-b"
	// instanceName := "your_instance_name"
	// newInstanceType = "e2-standard-2"

	ctx := context.Background()
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewInstancesRESTClient: %w", err)
	}
	defer instancesClient.Close()

	reqInstance := &computepb.GetInstanceRequest{
		Project:  projectID,
		Zone:     zone,
		Instance: instanceName,
	}

	instance, err := instancesClient.Get(ctx, reqInstance)
	if err != nil {
		return fmt.Errorf("unable to get instance: %w", err)
	}

	containsString := func(s []string, str string) bool {
		for _, v := range s {
			if v == str {
				return true
			}
		}

		return false
	}

	// Make sure that the machine is turned off
	if !containsString([]string{"TERMINATED", "STOPPED"}, instance.GetStatus()) {
		reqStop := &computepb.StopInstanceRequest{
			Project:  projectID,
			Zone:     zone,
			Instance: instanceName,
		}

		op, err := instancesClient.Stop(ctx, reqStop)
		if err != nil {
			return fmt.Errorf("unable to stop instance: %w", err)
		}

		if err = op.Wait(ctx); err != nil {
			return fmt.Errorf("unable to wait for the operation: %w", err)
		}
	}

	machineTypeUrl := fmt.Sprintf("zones/%s/machineTypes/%s", zone, newInstanceType)

	// Modify the machine definition with the new instance type
	updateReq := &computepb.SetMachineTypeInstanceRequest{
		Project:  projectID,
		Zone:     zone,
		Instance: instanceName,
		InstancesSetMachineTypeRequestResource: &computepb.InstancesSetMachineTypeRequest{
			MachineType: &machineTypeUrl,
		},
	}
	op, err := instancesClient.SetMachineType(ctx, updateReq)
	if err != nil {
		return fmt.Errorf("unable to update instance: %w", err)
	}

	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %w", err)
	}

	fmt.Fprintf(w, "Instance updated\n")

	return nil

}

// [END compute_instance_machine_type_update]
