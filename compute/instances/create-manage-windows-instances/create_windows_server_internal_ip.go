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

// [START compute_create_windows_instance_internal_ip]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/protobuf/proto"
)

// createWndowsServerInstanceInternalIP creates a new Windows Server instance
// that has only an internal IP address.
func createWndowsServerInstanceInternalIP(
	w io.Writer,
	projectID, zone, instanceName, machineType, sourceImageFamily, networkLink, subnetworkLink string,
) error {
	// projectID := "your_project_id"
	// zone := "europe-central2-b"
	// instanceName := "your_instance_name"
	// machineType := "n1-standard-1"
	// sourceImageFamily := "windows-2022"
	// networkLink := "global/networks/default"
	// subnetworkLink := "regions/europe-central2/subnetworks/default"

	ctx := context.Background()
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewInstancesRESTClient: %w", err)
	}
	defer instancesClient.Close()

	disk := &computepb.AttachedDisk{
		// Describe the size and source image of the boot disk to attach to the instance.
		InitializeParams: &computepb.AttachedDiskInitializeParams{
			DiskSizeGb: proto.Int64(64),
			SourceImage: proto.String(
				fmt.Sprintf(
					"projects/windows-cloud/global/images/family/%s",
					sourceImageFamily,
				),
			),
		},
		AutoDelete: proto.Bool(true),
		Boot:       proto.Bool(true),
	}

	network := &computepb.NetworkInterface{
		// You must verify or configure routes and firewall rules in your VPC network
		// to allow access to kms.windows.googlecloud.com.
		// More information about access to kms.windows.googlecloud.com:
		// https://cloud.google.com/compute/docs/instances/windows/creating-managing-windows-instances#kms-server

		// Additionally, you must enable Private Google Access for subnets in your VPC network
		// that contain Windows instances with only internal IP addresses.
		// More information about Private Google Access:
		// https://cloud.google.com/vpc/docs/configure-private-google-access#enabling
		Name:       proto.String(networkLink),
		Subnetwork: proto.String(subnetworkLink),
	}

	inst := &computepb.Instance{
		Name: proto.String(instanceName),
		Disks: []*computepb.AttachedDisk{
			disk,
		},
		MachineType: proto.String(fmt.Sprintf("zones/%s/machineTypes/%s", zone, machineType)),
		NetworkInterfaces: []*computepb.NetworkInterface{
			network,
		},
		// If you chose an image that supports Shielded VM,
		// you can optionally change the instance's Shielded VM settings.
		// ShieldedInstanceConfig: &computepb.ShieldedInstanceConfig{
		// 	EnableSecureBoot: proto.Bool(true),
		// 	EnableVtpm: proto.Bool(true),
		// 	EnableIntegrityMonitoring: proto.Bool(true),
		// },
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

// [END compute_create_windows_instance_internal_ip]
