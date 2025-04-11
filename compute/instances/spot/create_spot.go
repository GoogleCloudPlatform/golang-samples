//  Copyright 2024 Google LLC
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package snippets

// [START compute_spot_create]

import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/protobuf/proto"
)

// createSpotInstance creates a new Spot VM instance with Debian 10 operating system.
func createSpotInstance(w io.Writer, projectID, zone, instanceName string) error {
	// projectID := "your_project_id"
	// zone := "europe-central2-b"
	// instanceName := "your_instance_name"

	ctx := context.Background()
	imagesClient, err := compute.NewImagesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewImagesRESTClient: %w", err)
	}
	defer imagesClient.Close()

	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewInstancesRESTClient: %w", err)
	}
	defer instancesClient.Close()

	req := &computepb.GetFromFamilyImageRequest{
		Project: "debian-cloud",
		Family:  "debian-11",
	}

	image, err := imagesClient.GetFromFamily(ctx, req)
	if err != nil {
		return fmt.Errorf("getImageFromFamily: %w", err)
	}

	diskType := fmt.Sprintf("zones/%s/diskTypes/pd-standard", zone)
	disks := []*computepb.AttachedDisk{
		{
			AutoDelete: proto.Bool(true),
			Boot:       proto.Bool(true),
			InitializeParams: &computepb.AttachedDiskInitializeParams{
				DiskSizeGb:  proto.Int64(10),
				DiskType:    proto.String(diskType),
				SourceImage: proto.String(image.GetSelfLink()),
			},
			Type: proto.String(computepb.AttachedDisk_PERSISTENT.String()),
		},
	}

	req2 := &computepb.InsertInstanceRequest{
		Project: projectID,
		Zone:    zone,
		InstanceResource: &computepb.Instance{
			Name:        proto.String(instanceName),
			Disks:       disks,
			MachineType: proto.String(fmt.Sprintf("zones/%s/machineTypes/%s", zone, "n1-standard-1")),
			NetworkInterfaces: []*computepb.NetworkInterface{
				{
					Name: proto.String("global/networks/default"),
				},
			},
			Scheduling: &computepb.Scheduling{
				ProvisioningModel: proto.String(computepb.Scheduling_SPOT.String()),
			},
		},
	}
	op, err := instancesClient.Insert(ctx, req2)
	if err != nil {
		return fmt.Errorf("insert: %w", err)
	}

	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %w", err)
	}

	instance, err := instancesClient.Get(ctx, &computepb.GetInstanceRequest{
		Project:  projectID,
		Zone:     zone,
		Instance: instanceName,
	})

	if err != nil {
		return fmt.Errorf("createInstance: %w", err)
	}

	fmt.Fprintf(w, "Instance created: %v\n", instance)
	return nil
}

// [END compute_spot_create]
