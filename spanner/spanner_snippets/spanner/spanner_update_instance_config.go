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

// [START spanner_update_instance_config]

import (
	"context"
	"fmt"
	"io"
	"time"

	instance "cloud.google.com/go/spanner/admin/instance/apiv1"
	"cloud.google.com/go/spanner/admin/instance/apiv1/instancepb"
	"google.golang.org/genproto/protobuf/field_mask"
)

// updateInstanceConfig updates the custom spanner instance config
func updateInstanceConfig(w io.Writer, projectID, userConfigID string) error {
	// projectID := "my-project-id"
	// userConfigID := "custom-config", custom config names must start with the prefix “custom-”.

	// Add timeout to context.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	adminClient, err := instance.NewInstanceAdminClient(ctx)
	if err != nil {
		return err
	}
	defer adminClient.Close()
	config, err := adminClient.GetInstanceConfig(ctx, &instancepb.GetInstanceConfigRequest{
		Name: fmt.Sprintf("projects/%s/instanceConfigs/%s", projectID, userConfigID),
	})
	if err != nil {
		return fmt.Errorf("updateInstanceConfig.GetInstanceConfig: %w", err)
	}
	config.DisplayName = "updated custom instance config"
	config.Labels["updated"] = "true"
	op, err := adminClient.UpdateInstanceConfig(ctx, &instancepb.UpdateInstanceConfigRequest{
		InstanceConfig: config,
		UpdateMask: &field_mask.FieldMask{
			Paths: []string{"display_name", "labels"},
		},
		ValidateOnly: false,
	})
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Waiting for update operation on %s to complete...\n", userConfigID)
	// Wait for the instance configuration creation to finish.
	i, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("Waiting for instance config creation to finish failed: %w", err)
	}
	// The instance configuration may not be ready to serve yet.
	if i.State != instancepb.InstanceConfig_READY {
		fmt.Fprintf(w, "InstanceConfig state is not READY yet. Got state %v\n", i.State)
	}
	fmt.Fprintf(w, "Updated instance configuration [%s]\n", config.Name)
	return nil
}

// [END spanner_update_instance_config]
