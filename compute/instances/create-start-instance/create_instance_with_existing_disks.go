// Copyright 2021 Google LLC
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

// [START compute_instances_create_with_existing_disks]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/protobuf/proto"
)

// createWithExistingDisks create a new VM instance using selected disks.
// The first disk in diskNames will be used as boot disk.
func createWithExistingDisks(
	w io.Writer,
	projectID, zone, instanceName string,
	diskNames []string,
) error {
	// projectID := "your_project_id"
	// zone := "europe-central2-b"
	// instanceName := "your_instance_name"
	// diskNames := []string{"boot_disk", "disk1", "disk2"}

	ctx := context.Background()
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewInstancesRESTClient: %w", err)
	}
	defer instancesClient.Close()

	disksClient, err := compute.NewDisksRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewDisksRESTClient: %w", err)
	}
	defer disksClient.Close()

	disks := [](*computepb.Disk){}

	for _, diskName := range diskNames {
		reqDisk := &computepb.GetDiskRequest{
			Project: projectID,
			Zone:    zone,
			Disk:    diskName,
		}

		disk, err := disksClient.Get(ctx, reqDisk)
		if err != nil {
			return fmt.Errorf("unable to get disk: %w", err)
		}

		disks = append(disks, disk)
	}

	attachedDisks := [](*computepb.AttachedDisk){}

	for _, disk := range disks {
		attachedDisk := &computepb.AttachedDisk{
			Source: proto.String(disk.GetSelfLink()),
		}
		attachedDisks = append(attachedDisks, attachedDisk)
	}

	attachedDisks[0].Boot = proto.Bool(true)

	instanceResource := &computepb.Instance{
		Name:        proto.String(instanceName),
		Disks:       attachedDisks,
		MachineType: proto.String(fmt.Sprintf("zones/%s/machineTypes/n1-standard-1", zone)),
		NetworkInterfaces: []*computepb.NetworkInterface{
			{
				Name: proto.String("global/networks/default"),
			},
		},
	}

	req := &computepb.InsertInstanceRequest{
		Project:          projectID,
		Zone:             zone,
		InstanceResource: instanceResource,
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

// [END compute_instances_create_with_existing_disks]
