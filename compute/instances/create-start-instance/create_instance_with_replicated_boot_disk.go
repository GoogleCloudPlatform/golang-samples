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

// [START compute_instance_create_replicated_boot_disk]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/protobuf/proto"
)

// createReplicatedBootDisk creates a new VM with replicated boot disk
func createReplicatedBootDisk(
	w io.Writer,
	projectID, zone, snapshotLink, diskName, instanceName string,
	replicaZones []string,
) error {
	// projectID := "your_project_id"
	// zone := "europe-central2-b"
	// instanceName := "your_instance_name"
	// snapshotLink := "projects/project_name/global/snapshots/snapshot_name"

	ctx := context.Background()
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewInstancesRESTClient: %w", err)
	}
	defer instancesClient.Close()

	replicaZoneURLs := []string{
		fmt.Sprintf("projects/%s/zones/%s", projectID, replicaZones[0]),
		fmt.Sprintf("projects/%s/zones/%s", projectID, replicaZones[1]),
	}

	req := &computepb.InsertInstanceRequest{
		Project: projectID,
		Zone:    zone,
		InstanceResource: &computepb.Instance{
			Name: proto.String(instanceName),
			Disks: []*computepb.AttachedDisk{
				{
					InitializeParams: &computepb.AttachedDiskInitializeParams{
						DiskName:       proto.String(diskName),
						DiskSizeGb:     proto.Int64(200),
						SourceSnapshot: proto.String(snapshotLink),
						DiskType:       proto.String(fmt.Sprintf("zones/%s/diskTypes/pd-standard", zone)),
						ReplicaZones:   replicaZoneURLs,
					},
					Boot: proto.Bool(true),
					Type: proto.String(computepb.AttachedDisk_PERSISTENT.String()),
				},
			},
			MachineType: proto.String(fmt.Sprintf("zones/%s/machineTypes/n1-standard-1", zone)),
			NetworkInterfaces: []*computepb.NetworkInterface{
				{
					Name: proto.String("global/networks/default"),
				},
			},
		},
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

// [END compute_instance_create_replicated_boot_disk]
