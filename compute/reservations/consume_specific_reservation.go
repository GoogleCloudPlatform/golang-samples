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

// [START compute_consume_single_project_reservation]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/protobuf/proto"
)

// consumeSpecificReservation creates instance, consuming specific reservation
// Note: respective reservation should have SpecificReservationRequired: true
func consumeSpecificReservation(w io.Writer, projectID, zone, instanceName, reservationName string) error {
	// projectID := "your_project_id"
	// zone := "us-west3-a"
	// reservationName := "your_reservation_name"
	// instanceName := "your_instance_name"

	ctx := context.Background()
	machineType := fmt.Sprintf("zones/%s/machineTypes/%s", zone, "n2-standard-32")
	sourceImage := "projects/debian-cloud/global/images/family/debian-12"

	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewInstancesRESTClient: %w", err)
	}
	defer instancesClient.Close()

	req := &computepb.InsertInstanceRequest{
		Project: projectID,
		Zone:    zone,
		InstanceResource: &computepb.Instance{
			Disks: []*computepb.AttachedDisk{
				{
					InitializeParams: &computepb.AttachedDiskInitializeParams{
						DiskSizeGb:  proto.Int64(10),
						SourceImage: proto.String(sourceImage),
					},
					AutoDelete: proto.Bool(true),
					Boot:       proto.Bool(true),
					Type:       proto.String(computepb.AttachedDisk_PERSISTENT.String()),
				},
			},
			MachineType:    proto.String(machineType),
			MinCpuPlatform: proto.String("Intel Cascade Lake"),
			Name:           proto.String(instanceName),
			NetworkInterfaces: []*computepb.NetworkInterface{
				{
					Name: proto.String("global/networks/default"),
				},
			},
			// specifies particular reservation, which should be consumed
			ReservationAffinity: &computepb.ReservationAffinity{
				ConsumeReservationType: proto.String("SPECIFIC_RESERVATION"),
				Key:                    proto.String("compute.googleapis.com/reservation-name"),
				Values:                 []string{reservationName},
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
	fmt.Fprintf(w, "Instance created from reservation\n")

	return nil
}

// [END compute_consume_single_project_reservation]
