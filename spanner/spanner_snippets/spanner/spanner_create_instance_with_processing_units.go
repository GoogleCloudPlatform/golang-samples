// Copyright 2021 Google LLC
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

package spanner

// [START spanner_create_instance_with_processing_units]
import (
	"context"
	"fmt"
	"io"

	instance "cloud.google.com/go/spanner/admin/instance/apiv1"
	instancepb "google.golang.org/genproto/googleapis/spanner/admin/instance/v1"
	"google.golang.org/genproto/protobuf/field_mask"
)

func createInstanceWithProcessingUnits(w io.Writer, projectID, instanceID string) error {
	// projectID := "my-project-id"
	// instanceID := "my-instance"
	ctx := context.Background()
	instanceAdmin, err := instance.NewInstanceAdminClient(ctx)
	if err != nil {
		return fmt.Errorf("could not create instance admin client for project %s: %v", projectID, err)
	}
	defer instanceAdmin.Close()

	instanceName := fmt.Sprintf("projects/%s/instances/%s", projectID, instanceID)
	fmt.Fprintf(w, "Creating instance %s.", instanceName)

	op, err := instanceAdmin.CreateInstance(ctx, &instancepb.CreateInstanceRequest{
		Parent:     fmt.Sprintf("projects/%s", projectID),
		InstanceId: instanceID,
		Instance: &instancepb.Instance{
			Config:          fmt.Sprintf("projects/%s/instanceConfigs/%s", projectID, "regional-us-central1"),
			DisplayName:     "This is a display name.",
			ProcessingUnits: 500,
			Labels:          map[string]string{"cloud_spanner_samples": "true"},
		},
	})
	if err != nil {
		return fmt.Errorf("could not create instance %s: %v", fmt.Sprintf("projects/%s/instances/%s", projectID, instanceID), err)
	}
	fmt.Fprintf(w, "Waiting for operation on %s to complete...", instanceID)
	// Wait for the instance creation to finish.
	i, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("waiting for instance creation to finish failed: %v", err)
	}
	// The instance may not be ready to serve yet.
	if i.State != instancepb.Instance_READY {
		fmt.Fprintf(w, "instance state is not READY yet. Got state %v\n", i.State)
	}
	fmt.Fprintf(w, "Created instance [%s].\n", instanceID)

	instance, err := instanceAdmin.GetInstance(ctx, &instancepb.GetInstanceRequest{
		Name:      instanceName,
		FieldMask: &field_mask.FieldMask{Paths: []string{"processing_units"}},
	})
	if err != nil {
		return fmt.Errorf("failed to get instance [%s]: %v", instanceName, err)
	}
	fmt.Fprintf(w, "Instance %s has %d processing units.", instanceID, instance.ProcessingUnits)
	return nil
}

// [END spanner_create_instance_with_processing_units]
