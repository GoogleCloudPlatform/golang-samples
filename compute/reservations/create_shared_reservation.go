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

// [START compute_reservation_create_shared]
import (
	"context"
	"fmt"
	"io"

	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/protobuf/proto"
)

// Creates shared reservation from given template in particular zone
func createSharedReservation(w io.Writer, client ClientInterface, projectID, baseProjectId, zone, reservationName, sourceTemplate string) error {
	// client, err := compute.NewReservationsRESTClient(ctx)
	// projectID := "your_project_id". Destination of sharing.
	// baseProjectId := "your_project_id2". Project where the reservation will be created.
	// zone := "us-west3-a"
	// reservationName := "your_reservation_name"
	// sourceTemplate: existing template path. Following formats are allowed:
	//  	- projects/{project_id}/global/instanceTemplates/{template_name}
	//  	- projects/{project_id}/regions/{region}/instanceTemplates/{template_name}
	//  	- https://www.googleapis.com/compute/v1/projects/{project_id}/global/instanceTemplates/instanceTemplate
	//  	- https://www.googleapis.com/compute/v1/projects/{project_id}/regions/{region}/instanceTemplates/instanceTemplate

	ctx := context.Background()

	shareSettings := map[string]*computepb.ShareSettingsProjectConfig{
		projectID: {ProjectId: proto.String(projectID)},
	}

	req := &computepb.InsertReservationRequest{
		Project: baseProjectId,
		ReservationResource: &computepb.Reservation{
			Name: proto.String(reservationName),
			Zone: proto.String(zone),
			SpecificReservation: &computepb.AllocationSpecificSKUReservation{
				Count:                  proto.Int64(2),
				SourceInstanceTemplate: proto.String(sourceTemplate),
			},
			ShareSettings: &computepb.ShareSettings{
				ProjectMap: shareSettings,
				ShareType:  proto.String("SPECIFIC_PROJECTS"),
			},
		},
		Zone: zone,
	}

	op, err := client.Insert(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to create reservation: %w", err)
	}

	if op != nil {
		if err = op.Wait(ctx); err != nil {
			return fmt.Errorf("unable to wait for the operation: %w", err)
		}
	}

	fmt.Fprintf(w, "Reservation created\n")

	return nil
}

// [END compute_reservation_create_shared]
