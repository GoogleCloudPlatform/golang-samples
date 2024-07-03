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
	"google.golang.org/protobuf/proto"
)

// [START compute_ip_address_assign_static_external_new_vm]

// assignStaticExternalToNewVM creates a new VM instance and assigns a static external IP address to it.
// NOTE: ip address is expected to exist and be located in the same region as new VM
func assignStaticExternalToNewVM(w io.Writer, projectID, zone, instanceName, ipAddress string) error {
	// projectID := "your_project_id"
	// zone := "europe-central2-b"
	// instanceName := "your_instance_name"
	// ipAddress := 301.222.11.123

	ctx := context.Background()
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewInstancesRESTClient: %w", err)
	}
	defer instancesClient.Close()

	imagesClient, err := compute.NewImagesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewImagesRESTClient: %w", err)
	}
	defer imagesClient.Close()

	// List of public operating system (OS) images: https://cloud.google.com/compute/docs/images/os-details.
	newestDebianReq := &computepb.GetFromFamilyImageRequest{
		Project: "debian-cloud",
		Family:  "debian-12",
	}
	newestDebian, err := imagesClient.GetFromFamily(ctx, newestDebianReq)
	if err != nil {
		return fmt.Errorf("unable to get image from family: %w", err)
	}

	req := &computepb.InsertInstanceRequest{
		Project: projectID,
		Zone:    zone,
		InstanceResource: &computepb.Instance{
			Name: proto.String(instanceName),
			Disks: []*computepb.AttachedDisk{
				{
					InitializeParams: &computepb.AttachedDiskInitializeParams{
						DiskSizeGb:  proto.Int64(10),
						SourceImage: newestDebian.SelfLink,
						DiskType:    proto.String(fmt.Sprintf("zones/%s/diskTypes/pd-standard", zone)),
					},
					AutoDelete: proto.Bool(true),
					Boot:       proto.Bool(true),
					Type:       proto.String(computepb.AttachedDisk_PERSISTENT.String()),
				},
			},
			MachineType: proto.String(fmt.Sprintf("zones/%s/machineTypes/n1-standard-1", zone)),
			NetworkInterfaces: []*computepb.NetworkInterface{
				{
					AccessConfigs: []*computepb.AccessConfig{
						{
							Type:        proto.String(computepb.AccessConfig_ONE_TO_ONE_NAT.String()),
							Name:        proto.String("External NAT"),
							NetworkTier: proto.String(computepb.AccessConfig_PREMIUM.String()),
							NatIP:       proto.String(ipAddress),
						},
					},
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

	fmt.Fprintf(w, "Static address %s assigned to new VM", ipAddress)

	return nil
}

// [END compute_ip_address_assign_static_external_new_vm]
