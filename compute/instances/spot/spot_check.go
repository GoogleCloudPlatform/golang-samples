//  Copyright 2024 Google LLC
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package snippets

// [START compute_spot_check]

import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
)

// isSpotVM checks if a given instance is a Spot VM or not.
func isSpotVM(w io.Writer, projectID, zone, instanceName string) (bool, error) {
	// projectID := "your_project_id"
	// zone := "europe-central2-b"
	// instanceName := "your_instance_name"
	ctx := context.Background()
	client, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return false, fmt.Errorf("NewInstancesRESTClient: %w", err)
	}
	defer client.Close()

	req := &computepb.GetInstanceRequest{
		Project:  projectID,
		Zone:     zone,
		Instance: instanceName,
	}

	instance, err := client.Get(ctx, req)
	if err != nil {
		return false, fmt.Errorf("GetInstance: %w", err)
	}

	isSpot := instance.GetScheduling().GetProvisioningModel() == computepb.Scheduling_SPOT.String()

	var isSpotMessage string
	if !isSpot {
		isSpotMessage = " not"
	}
	fmt.Fprintf(w, "Instance %s is%s spot\n", instanceName, isSpotMessage)

	return instance.GetScheduling().GetProvisioningModel() == computepb.Scheduling_SPOT.String(), nil
}

// [END compute_spot_check]
