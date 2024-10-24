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

// [START compute_reservation_get]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
)

// Get certain reservation for given project and zone
func getReservation(w io.Writer, projectID, zone, reservationName string) (*computepb.Reservation, error) {
	// projectID := "your_project_id"
	// zone := "us-west3-a"
	// reservationName := "your_reservation_name"

	ctx := context.Background()
	reservationsClient, err := compute.NewReservationsRESTClient(ctx)
	if err != nil {
		return nil, err
	}
	defer reservationsClient.Close()

	req := &computepb.GetReservationRequest{
		Project:     projectID,
		Reservation: reservationName,
		Zone:        zone,
	}

	reservation, err := reservationsClient.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("unable to delete reservation: %w", err)
	}

	fmt.Fprintf(w, "Reservation: %s\n", reservation.GetName())

	return reservation, nil
}

// [END compute_reservation_get]
