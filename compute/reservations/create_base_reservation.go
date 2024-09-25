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

// [START compute_reservation_create]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/protobuf/proto"
)

// Creates the reservation with accelerated image
func createBaseReservation(w io.Writer, projectID, zone, reservationName string) error {
	// projectID := "your_project_id"
	// zone := "us-west3-a"
	// reservationName := "your_reservation_name"

	ctx := context.Background()
	reservationsClient, err := compute.NewReservationsRESTClient(ctx)
	if err != nil {
		return err
	}
	defer reservationsClient.Close()

	// Creating reservation based on direct properties
	req := &computepb.InsertReservationRequest{
		Project: projectID,
		ReservationResource: &computepb.Reservation{
			Name: proto.String(reservationName),
			Zone: proto.String(zone),
			SpecificReservation: &computepb.AllocationSpecificSKUReservation{
				Count: proto.Int64(2),
				// Properties, which allows customising instances
				InstanceProperties: &computepb.AllocationSpecificSKUAllocationReservedInstanceProperties{
					// Attaching GPUs to the reserved VMs
					// Read more: https://cloud.google.com/compute/docs/gpus#n1-gpus
					GuestAccelerators: []*computepb.AcceleratorConfig{
						{
							AcceleratorCount: proto.Int32(1),
							AcceleratorType:  proto.String("nvidia-tesla-t4"),
						},
					},
					// Including local SSD disks
					LocalSsds: []*computepb.AllocationSpecificSKUAllocationAllocatedInstancePropertiesReservedDisk{
						{
							DiskSizeGb: proto.Int64(375),
							Interface:  proto.String("NVME"),
						},
					},
					MachineType: proto.String("n1-standard-2"),
					// Specifying minimum CPU platform
					// Read more: https://cloud.google.com/compute/docs/instances/specify-min-cpu-platform
					MinCpuPlatform: proto.String("Intel Skylake"),
				},
			},
		},
		Zone: zone,
	}

	op, err := reservationsClient.Insert(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to create reservation: %w", err)
	}

	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %w", err)
	}

	fmt.Fprintf(w, "Reservation created\n")

	return nil
}

// [END compute_reservation_create]
