// Copyright 2020 Google LLC
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

// [START spanner_create_instance_without_default_backup_schedule_type]
import (
	"context"
	"fmt"
	"io"

	instance "cloud.google.com/go/spanner/admin/instance/apiv1"
	"cloud.google.com/go/spanner/admin/instance/apiv1/instancepb"
)

func createInstanceWithoutDefaultBackupScheduleType(w io.Writer, projectID, instanceID string) error {
	// projectID := "my-project-id"
	// instanceID := "my-instance"
	ctx := context.Background()
	instanceAdmin, err := instance.NewInstanceAdminClient(ctx)
	if err != nil {
		return err
	}
	defer instanceAdmin.Close()

	op, err := instanceAdmin.CreateInstance(ctx, &instancepb.CreateInstanceRequest{
		Parent:     fmt.Sprintf("projects/%s", projectID),
		InstanceId: instanceID,
		Instance: &instancepb.Instance{
			Config:                    fmt.Sprintf("projects/%s/instanceConfigs/%s", projectID, "regional-us-central1"),
			DisplayName:               instanceID,
			NodeCount:                 1,
			Labels:                    map[string]string{"cloud_spanner_samples": "true"},
			DefaultBackupScheduleType: instancepb.Instance_NONE,
		},
	})
	if err != nil {
		return fmt.Errorf("could not create instance %s: %w", fmt.Sprintf("projects/%s/instances/%s", projectID, instanceID), err)
	}
	// Wait for the instance creation to finish.
	i, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("waiting for instance creation to finish failed: %w", err)
	}
	// The instance may not be ready to serve yet.
	if i.State != instancepb.Instance_READY {
		fmt.Fprintf(w, "instance state is not READY yet. Got state %v\n", i.State)
	}
	fmt.Fprintf(w, "Created instance [%s]\n", instanceID)
	return nil
}

// [END spanner_create_instance_without_default_backup_schedule_type]
