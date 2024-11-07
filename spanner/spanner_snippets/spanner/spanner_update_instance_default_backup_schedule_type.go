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

package spanner

// [START spanner_update_instance_default_backup_schedule_type]
import (
	"context"
	"fmt"
	"io"

	instance "cloud.google.com/go/spanner/admin/instance/apiv1"
	"cloud.google.com/go/spanner/admin/instance/apiv1/instancepb"
	"google.golang.org/genproto/protobuf/field_mask"
)

// This function updates instance default backup schedule type to AUTOMATIC.  This means a default
// backup schedule will be created automatically on creation of a database within the instance.
func updateInstanceDefaultBackupScheduleType(w io.Writer, projectID, instanceID string) error {
	// projectID := "my-project-id"
	// instanceID := "my-instance"
	ctx := context.Background()
	instanceAdmin, err := instance.NewInstanceAdminClient(ctx)
	if err != nil {
		return err
	}
	defer instanceAdmin.Close()

	// Updates the default backup schedule type field of an instance.  The field mask is required to
	// indicate which field is being updated.
	req := &instancepb.UpdateInstanceRequest{
		Instance: &instancepb.Instance{
			Name: fmt.Sprintf("projects/%s/instances/%s", projectID, instanceID),
			// Controls the default backup behavior for new databases within the instance.
			DefaultBackupScheduleType: instancepb.Instance_AUTOMATIC,
		},
		FieldMask: &field_mask.FieldMask{
			Paths: []string{"default_backup_schedule_type"},
		},
	}
	op, err := instanceAdmin.UpdateInstance(ctx, req)
	if err != nil {
		return fmt.Errorf("could not update instance %s: %w", fmt.Sprintf("projects/%s/instances/%s", projectID, instanceID), err)
	}
	// Wait for the instance update to finish.
	_, err = op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("waiting for instance update to finish failed: %w", err)
	}

	fmt.Fprintf(w, "Updated instance [%s]\n", instanceID)
	return nil
}

// [END spanner_update_instance_default_backup_schedule_type]
