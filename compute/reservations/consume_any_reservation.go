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

// [START compute_consume_any_matching_reservation]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/protobuf/proto"
)

// consumeAnyReservation creates instance, consuming any available reservation
func consumeAnyReservation(w io.Writer, projectID, zone, instanceName, sourceTemplate string) error {
	// projectID := "your_project_id"
	// zone := "us-west3-a"
	// instanceName := "your_instance_name"
	// sourceTemplate: existing template path. Following formats are allowed:
	//  	- projects/{project_id}/global/instanceTemplates/{template_name}
	//  	- projects/{project_id}/regions/{region}/instanceTemplates/{template_name}
	//  	- https://www.googleapis.com/compute/v1/projects/{project_id}/global/instanceTemplates/instanceTemplate
	//  	- https://www.googleapis.com/compute/v1/projects/{project_id}/regions/{region}/instanceTemplates/instanceTemplate

	ctx := context.Background()

	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewInstancesRESTClient: %w", err)
	}
	defer instancesClient.Close()

	req := &computepb.InsertInstanceRequest{
		Project:                projectID,
		Zone:                   zone,
		SourceInstanceTemplate: proto.String(sourceTemplate),
		InstanceResource: &computepb.Instance{
			Name: proto.String(instanceName),
			// specifies that any matching reservation should be consumed
			ReservationAffinity: &computepb.ReservationAffinity{
				ConsumeReservationType: proto.String("ANY_RESERVATION"),
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

// [END compute_consume_any_matching_reservation]
