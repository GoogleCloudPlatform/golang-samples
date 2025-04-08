// Copyright 2022 Google LLC
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

// [START spanner_create_instance_config]

import (
	"context"
	"fmt"
	"io"
	"time"

	instance "cloud.google.com/go/spanner/admin/instance/apiv1"
	"cloud.google.com/go/spanner/admin/instance/apiv1/instancepb"
)

// createInstanceConfig creates a custom spanner instance config
func createInstanceConfig(w io.Writer, projectID, userConfigID, baseConfigID string) error {
	// projectID := "my-project-id"
	// userConfigID := "custom-config", custom config names must start with the prefix “custom-”.
	// baseConfigID := "my-base-config"

	// Add timeout to context.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	adminClient, err := instance.NewInstanceAdminClient(ctx)
	if err != nil {
		return err
	}
	defer adminClient.Close()
	baseConfig, err := adminClient.GetInstanceConfig(ctx, &instancepb.GetInstanceConfigRequest{
		Name: fmt.Sprintf("projects/%s/instanceConfigs/%s", projectID, baseConfigID),
	})
	if err != nil {
		return fmt.Errorf("createInstanceConfig.GetInstanceConfig: %w", err)
	}
	if baseConfig.OptionalReplicas == nil || len(baseConfig.OptionalReplicas) == 0 {
		return fmt.Errorf("CreateInstanceConfig expects base config with at least from the list of optional replicas")
	}
	op, err := adminClient.CreateInstanceConfig(ctx, &instancepb.CreateInstanceConfigRequest{
		Parent: fmt.Sprintf("projects/%s", projectID),
		// Custom config names must start with the prefix “custom-”.
		InstanceConfigId: userConfigID,
		InstanceConfig: &instancepb.InstanceConfig{
			Name:        fmt.Sprintf("projects/%s/instanceConfigs/%s", projectID, userConfigID),
			DisplayName: "custom-golang-samples",
			ConfigType:  instancepb.InstanceConfig_USER_MANAGED,
			// The replicas for the custom instance configuration must include all the replicas of the base
			// configuration, in addition to at least one from the list of optional replicas of the base
			// configuration.
			Replicas:   append(baseConfig.Replicas, baseConfig.OptionalReplicas...),
			BaseConfig: baseConfig.Name,
			Labels:     map[string]string{"go_cloud_spanner_samples": "true"},
		},
	})
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Waiting for create operation on projects/%s/instanceConfigs/%s to complete...\n", projectID, userConfigID)
	// Wait for the instance configuration creation to finish.
	i, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("Waiting for instance config creation to finish failed: %w", err)
	}
	// The instance configuration may not be ready to serve yet.
	if i.State != instancepb.InstanceConfig_READY {
		fmt.Fprintf(w, "InstanceConfig state is not READY yet. Got state %v\n", i.State)
	}
	fmt.Fprintf(w, "Created instance configuration [%s]\n", userConfigID)
	return nil
}

// [END spanner_create_instance_config]
