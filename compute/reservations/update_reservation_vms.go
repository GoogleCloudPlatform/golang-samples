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

// [START compute_reservation_vms_update]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/protobuf/proto"
)

// Update reservation's number of allocated resources
func updateReservationVMS(w io.Writer, projectID, zone, reservationName string, numberOfVMs int64) error {
	// projectID := "your_project_id"
	// zone := "us-west3-a"
	// reservationName := "your_reservation_name"
	// numberOfVMs := 5

	ctx := context.Background()
	reservationsClient, err := compute.NewReservationsRESTClient(ctx)
	if err != nil {
		return err
	}
	defer reservationsClient.Close()

	req := &computepb.ResizeReservationRequest{
		Project:     projectID,
		Reservation: reservationName,
		ReservationsResizeRequestResource: &computepb.ReservationsResizeRequest{
			SpecificSkuCount: proto.Int64(numberOfVMs),
		},
		Zone: zone,
	}

	op, err := reservationsClient.Resize(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to create reservation: %w", err)
	}

	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %w", err)
	}

	return nil
}

// [END compute_reservation_vms_update]
