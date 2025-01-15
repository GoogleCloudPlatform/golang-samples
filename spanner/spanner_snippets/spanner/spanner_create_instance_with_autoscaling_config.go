// Copyright 2023 Google LLC
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

// [START spanner_create_instance_with_autoscaling_config]
import (
	"context"
	"fmt"
	"io"

	instance "cloud.google.com/go/spanner/admin/instance/apiv1"
	"cloud.google.com/go/spanner/admin/instance/apiv1/instancepb"
	"google.golang.org/genproto/protobuf/field_mask"
)

// Example of creating an autoscaling instance with Go.
// projectID is the ID of the project that the new instance will be in.
// instanceID is the ID of the new instance to be created.
func createInstanceWithAutoscalingConfig(w io.Writer, projectID, instanceID string) error {
	// projectID := "my-project-id"
	// instanceID := "my-instance"
	ctx := context.Background()
	instanceAdmin, err := instance.NewInstanceAdminClient(ctx)
	if err != nil {
		return fmt.Errorf("could not create instance admin client for project %s: %w", projectID, err)
	}
	defer instanceAdmin.Close()

	instanceName := fmt.Sprintf("projects/%s/instances/%s", projectID, instanceID)
	fmt.Fprintf(w, "Creating instance %s.", instanceName)

	op, err := instanceAdmin.CreateInstance(ctx, &instancepb.CreateInstanceRequest{
		Parent:     fmt.Sprintf("projects/%s", projectID),
		InstanceId: instanceID,
		Instance: &instancepb.Instance{
			Config:      fmt.Sprintf("projects/%s/instanceConfigs/%s", projectID, "regional-us-central1"),
			DisplayName: "Create instance example",
			AutoscalingConfig: &instancepb.AutoscalingConfig{
				AutoscalingLimits: &instancepb.AutoscalingConfig_AutoscalingLimits{
					MinLimit: &instancepb.AutoscalingConfig_AutoscalingLimits_MinNodes{
						MinNodes: 1,
					},
					MaxLimit: &instancepb.AutoscalingConfig_AutoscalingLimits_MaxNodes{
						MaxNodes: 2,
					},
				},
				AutoscalingTargets: &instancepb.AutoscalingConfig_AutoscalingTargets{
					HighPriorityCpuUtilizationPercent: 65,
					StorageUtilizationPercent:         95,
				},
			},
			Labels:  map[string]string{"cloud_spanner_samples": "true"},
			Edition: instancepb.Instance_ENTERPRISE_PLUS,
		},
	})
	if err != nil {
		return fmt.Errorf("could not create instance %s: %w", instanceName, err)
	}
	fmt.Fprintf(w, "Waiting for operation on %s to complete...", instanceID)
	// Wait for the instance creation to finish.
	i, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("waiting for instance creation to finish failed: %w", err)
	}
	// The instance may not be ready to serve yet.
	if i.State != instancepb.Instance_READY {
		fmt.Fprintf(w, "instance state is not READY yet. Got state %v\n", i.State)
	}
	fmt.Fprintf(w, "Created instance [%s].\n", instanceID)

	instance, err := instanceAdmin.GetInstance(ctx, &instancepb.GetInstanceRequest{
		Name: instanceName,
		// Get the autoscaling_config field from the newly created instance.
		FieldMask: &field_mask.FieldMask{Paths: []string{"autoscaling_config"}},
	})
	if err != nil {
		return fmt.Errorf("failed to get instance [%s]: %w", instanceName, err)
	}
	fmt.Fprintf(w, "Instance %s has autoscaling_config: %s.", instanceID, instance.AutoscalingConfig)
	return nil
}

// [END spanner_create_instance_with_autoscaling_config]
