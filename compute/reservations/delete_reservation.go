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

// [START compute_reservation_delete]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
)

// Deletes the reservation for given project and zone
func deleteReservation(w io.Writer, projectID, zone, reservationName string) error {
	// projectID := "your_project_id"
	// zone := "us-west3-a"
	// reservationName := "your_reservation_name"

	ctx := context.Background()
	reservationsClient, err := compute.NewReservationsRESTClient(ctx)
	if err != nil {
		return err
	}
	defer reservationsClient.Close()

	req := &computepb.DeleteReservationRequest{
		Project:     projectID,
		Reservation: reservationName,
		Zone:        zone,
	}

	op, err := reservationsClient.Delete(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to delete reservation: %w", err)
	}

	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %w", err)
	}

	fmt.Fprintf(w, "Reservation deleted\n")

	return nil
}

// [END compute_reservation_delete]
