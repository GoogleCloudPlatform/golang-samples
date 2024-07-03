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

// [START compute_custom_machine_type_create_without_helper]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/protobuf/proto"
)

// createInstanceWithCustomMachineTypeWithoutHelper create new VM instances
// without using a CustomMachineType struct.
func createInstanceWithCustomMachineTypeWithoutHelper(
	w io.Writer,
	projectID, zone, instanceName, cpuSeries string,
	coreCount, memory int,
) error {
	// projectID := "your_project_id"
	// zone := "europe-central2-b"
	// instanceName := "your_instance_name"
	// cpuSeries := "N1"
	// coreCount := 2 // number of CPU cores you want to use.
	// memory := 256 // the amount of memory for the VM instance, in megabytes.

	// The coreCount and memory values are not validated anywhere and can be rejected by the API.

	ctx := context.Background()
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewInstancesRESTClient: %w", err)
	}
	defer instancesClient.Close()

	mt := fmt.Sprintf("zones/%s/machineTypes/%v-%v-%v", zone, cpuSeries, coreCount, memory)
	inst := &computepb.Instance{
		Name: proto.String(instanceName),
		Disks: []*computepb.AttachedDisk{
			{
				InitializeParams: &computepb.AttachedDiskInitializeParams{
					DiskSizeGb: proto.Int64(10),
					SourceImage: proto.String(
						"projects/debian-cloud/global/images/family/debian-12",
					),
				},
				AutoDelete: proto.Bool(true),
				Boot:       proto.Bool(true),
			},
		},
		MachineType: proto.String(mt),
		NetworkInterfaces: []*computepb.NetworkInterface{
			{
				Name: proto.String("global/networks/default"),
			},
		},
	}

	req := &computepb.InsertInstanceRequest{
		Project:          projectID,
		Zone:             zone,
		InstanceResource: inst,
	}

	op, err := instancesClient.Insert(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to create instance: %w", err)
	}

	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %w", err)
	}

	fmt.Fprintf(w, "Instance created\n")

	return nil
}

// [END compute_custom_machine_type_create_without_helper]
