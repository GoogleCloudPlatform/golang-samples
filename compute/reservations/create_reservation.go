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

// [START compute_reservation_create_template]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/protobuf/proto"
)

// Creates the reservation from given template in particular zone
func createReservation(w io.Writer, projectID, zone, reservationName, sourceTemplate string) error {
	// projectID := "your_project_id"
	// zone := "us-west3-a"
	// reservationName := "your_reservation_name"
	// template: existing template path. Following formats are allowed:
	//  	- projects/{project_id}/global/instanceTemplates/{template_name}
	//  	- projects/{project_id}/regions/{region}/instanceTemplates/{template_name}
	//  	- https://www.googleapis.com/compute/v1/projects/{project_id}/global/instanceTemplates/instanceTemplate
	//  	- https://www.googleapis.com/compute/v1/projects/{project_id}/regions/{region}/instanceTemplates/instanceTemplate

	ctx := context.Background()
	reservationsClient, err := compute.NewReservationsRESTClient(ctx)
	if err != nil {
		return err
	}
	defer reservationsClient.Close()

	req := &computepb.InsertReservationRequest{
		Project: projectID,
		ReservationResource: &computepb.Reservation{
			Name: proto.String(reservationName),
			Zone: proto.String(zone),
			SpecificReservation: &computepb.AllocationSpecificSKUReservation{
				Count:                  proto.Int64(2),
				SourceInstanceTemplate: proto.String(sourceTemplate),
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

// [END compute_reservation_create_template]
