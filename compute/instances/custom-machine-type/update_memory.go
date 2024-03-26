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

package snippets

// [START compute_custom_machine_type_update_memory]
import (
	"context"
	"fmt"
	"io"
	"strings"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/protobuf/proto"
)

// modifyInstanceWithExtendedMemory sends an instance creation request
// to the Compute Engine API and waits for it to complete.
func modifyInstanceWithExtendedMemory(
	w io.Writer,
	projectID, zone, instanceName string,
	newMemory int,
) error {
	// projectID := "your_project_id"
	// zone := "europe-central2-b"
	// instanceName := "your_instance_name"
	// newMemory := 256 // the amount of memory for the VM instance, in megabytes.

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

	if !(strings.Contains(instance.GetMachineType(), "machineTypes/n1-") ||
		strings.Contains(instance.GetMachineType(), "machineTypes/n2-") ||
		strings.Contains(instance.GetMachineType(), "machineTypes/n2d-")) {
		return fmt.Errorf("extra memory is available only for N1, N2 and N2D CPUs")
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

	// Modify the machine definition, remember that extended memory
	// is available only for N1, N2 and N2D CPUs
	machineType := instance.GetMachineType()
	start := machineType[:strings.LastIndex(machineType, "-")]

	updateReq := &computepb.SetMachineTypeInstanceRequest{
		Project:  projectID,
		Zone:     zone,
		Instance: instanceName,
		InstancesSetMachineTypeRequestResource: &computepb.InstancesSetMachineTypeRequest{
			MachineType: proto.String(fmt.Sprintf("%s-%v-ext", start, newMemory)),
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

// [END compute_custom_machine_type_update_memory]
