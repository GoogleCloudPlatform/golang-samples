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

// [START compute_instances_create]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
	"google.golang.org/protobuf/proto"
)

// createInstance sends an instance creation request to the Compute Engine API and waits for it to complete.
func createInstance(w io.Writer, projectID, zone, instanceName, machineType, sourceImage, networkName string) error {
	// projectID := "your_project_id"
	// zone := "europe-central2-b"
	// instanceName := "your_instance_name"
	// machineType := "n1-standard-1"
	// sourceImage := "projects/debian-cloud/global/images/family/debian-10"
	// networkName := "global/networks/default"

	ctx := context.Background()
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewInstancesRESTClient: %v", err)
	}
	defer instancesClient.Close()

	req := &computepb.InsertInstanceRequest{
		Project: projectID,
		Zone:    zone,
		InstanceResource: &computepb.Instance{
			Name: proto.String(instanceName),
			Disks: []*computepb.AttachedDisk{
				{
					InitializeParams: &computepb.AttachedDiskInitializeParams{
						DiskSizeGb:  proto.Int64(10),
						SourceImage: proto.String(sourceImage),
					},
					AutoDelete: proto.Bool(true),
					Boot:       proto.Bool(true),
					Type:       computepb.AttachedDisk_PERSISTENT.Enum(),
				},
			},
			MachineType: proto.String(fmt.Sprintf("zones/%s/machineTypes/%s", zone, machineType)),
			NetworkInterfaces: []*computepb.NetworkInterface{
				{
					Name: proto.String(networkName),
				},
			},
		},
	}

	op, err := instancesClient.Insert(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to create instance: %v", err)
	}

	zoneOperationsClient, err := compute.NewZoneOperationsRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewZoneOperationsRESTClient: %v", err)
	}
	defer zoneOperationsClient.Close()

	for {
		waitReq := &computepb.WaitZoneOperationRequest{
			Operation: op.Proto().GetName(),
			Project:   projectID,
			Zone:      zone,
		}
		zoneOp, err := zoneOperationsClient.Wait(ctx, waitReq)
		if err != nil {
			return fmt.Errorf("unable to wait for the operation: %v", err)
		}

		if *zoneOp.Status.Enum() == computepb.Operation_DONE {
			fmt.Fprintf(w, "Instance created\n")
			break
		}
	}

	return nil
}

// [END compute_instances_create]
