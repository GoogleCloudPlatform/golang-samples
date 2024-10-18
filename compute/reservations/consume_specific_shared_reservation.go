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

// [START compute_consume_specific_shared_reservation]
import (
	"context"
	"fmt"
	"io"

	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/protobuf/proto"
)

// consumeSpecificSharedReservation consumes specific shared reservation in particular zone
func consumeSpecificSharedReservation(w io.Writer, client InstanceClientInterface, projectID, baseProjectId, zone, instanceName, reservationName string) error {
	// client, err := compute.NewInstancesRESTClient(ctx)
	// projectID := "your_project_id". Project where reservation is created.
	// baseProjectId := "shared_project_id". Project where instance will be consumed and created.
	// zone := "us-west3-a"
	// reservationName := "your_reservation_name"
	// instanceName := "your_instance_name"

	ctx := context.Background()
	machineType := fmt.Sprintf("zones/%s/machineTypes/%s", zone, "n2-standard-32")
	sourceImage := "projects/debian-cloud/global/images/family/debian-12"
	sharedReservation := fmt.Sprintf("projects/%s/reservations/%s", baseProjectId, reservationName)

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
				Values:                 []string{sharedReservation},
			},
		},
	}

	op, err := client.Insert(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to create instance: %w", err)
	}

	if op != nil {
		if err = op.Wait(ctx); err != nil {
			return fmt.Errorf("unable to wait for the operation: %w", err)
		}
	}
	fmt.Fprintf(w, "Instance created from shared reservation\n")

	return nil
}

// [END compute_consume_specific_shared_reservation]
