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

// getImageFromFamily retrieves the newest image that is part of a given family in a project.
func getImageFromFamily(project, family string) (*computepb.Image, error) {
	ctx := context.Background()
	client, err := compute.NewImagesRESTClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewImagesRESTClient: %w", err)
	}
	defer client.Close()

	req := &computepb.GetFromFamilyImageRequest{
		Project: project,
		Family:  family,
	}

	resp, err := client.GetFromFamily(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("GetFromFamily: %w", err)
	}
	return resp, nil
}

// diskFromImage creates an AttachedDisk object to be used in VM instance creation using an image as the source.
func diskFromImage(diskType string, diskSizeGb int64, boot bool, sourceImage string, autoDelete bool) *computepb.AttachedDisk {
	return &computepb.AttachedDisk{
		AutoDelete: proto.Bool(autoDelete),
		Boot:       proto.Bool(boot),
		InitializeParams: &computepb.AttachedDiskInitializeParams{
			DiskSizeGb:  proto.Int64(diskSizeGb),
			DiskType:    proto.String(diskType),
			SourceImage: proto.String(sourceImage),
		},
		Type: proto.String(computepb.AttachedDisk_PERSISTENT.String()),
	}
}

// createInstance sends an instance creation request to the Compute Engine API and waits for it to complete.
func createInstance(ctx context.Context, projectID, zone, instanceName string, disks []*computepb.AttachedDisk, machineType, networkLink string, spot bool) (*computepb.Instance, error) {
	client, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewInstancesRESTClient: %w", err)
	}
	defer client.Close()

	networkInterface := &computepb.NetworkInterface{
		Name: proto.String(networkLink),
	}

	instance := &computepb.Instance{
		Name:              proto.String(instanceName),
		Disks:             disks,
		MachineType:       proto.String(fmt.Sprintf("zones/%s/machineTypes/%s", zone, machineType)),
		NetworkInterfaces: []*computepb.NetworkInterface{networkInterface},
		Scheduling: &computepb.Scheduling{
			ProvisioningModel: proto.String(computepb.Scheduling_SPOT.String()),
		},
	}

	req := &computepb.InsertInstanceRequest{
		Project:          projectID,
		Zone:             zone,
		InstanceResource: instance,
	}

	op, err := client.Insert(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("Insert: %w", err)
	}

	opClient, err := compute.NewZoneOperationsRESTClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewZoneOperationsRESTClient: %w", err)
	}
	defer opClient.Close()

	if err = op.Wait(ctx); err != nil {
		return nil, fmt.Errorf("unable to wait for the operation: %w", err)
	}

	return client.Get(ctx, &computepb.GetInstanceRequest{
		Project:  projectID,
		Zone:     zone,
		Instance: instanceName,
	})
}

// createSpotInstance creates a new Spot VM instance with Debian 10 operating system.
func createSpotInstance(w io.Writer, projectID, zone, instanceName string) error {
	// projectID := "your_project_id"
	// zone := "europe-central2-b"
	// instanceName := "your_instance_name"

	ctx := context.Background()
	image, err := getImageFromFamily("debian-cloud", "debian-11")
	if err != nil {
		return fmt.Errorf("getImageFromFamily: %w", err)
	}

	diskType := fmt.Sprintf("zones/%s/diskTypes/pd-standard", zone)
	disks := []*computepb.AttachedDisk{
		diskFromImage(diskType, 10, true, image.GetSelfLink(), true),
	}

	instance, err := createInstance(ctx, projectID, zone, instanceName, disks, "n1-standard-1", "global/networks/default", true)
	if err != nil {
		return fmt.Errorf("createInstance: %w", err)
	}

	fmt.Fprintf(w, "Instance created: %v\n", instance)
	return nil
}

// [END compute_spot_create]
