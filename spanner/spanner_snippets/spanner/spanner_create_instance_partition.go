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

// [START spanner_create_instance_partition]
import (
	"context"
	"fmt"
	"io"

	instance "cloud.google.com/go/spanner/admin/instance/apiv1"
	"cloud.google.com/go/spanner/admin/instance/apiv1/instancepb"
)

// Example of creating an instance partition with Go.
// projectID is the ID of the project that the new instance partition will be in.
// instanceID is the ID of the instance that the new instance partition will be in.
// instancePartitionID is the ID of the new instance partition to be created.
func createInstancePartition(w io.Writer, projectID, instanceID, instancePartitionID string) error {
	// projectID := "my-project-id"
	// instanceID := "my-instance"
	// instancePartitionID := "my-instance-partition"
	ctx := context.Background()
	instanceAdmin, err := instance.NewInstanceAdminClient(ctx)
	if err != nil {
		return err
	}
	defer instanceAdmin.Close()

	op, err := instanceAdmin.CreateInstancePartition(ctx, &instancepb.CreateInstancePartitionRequest{
		Parent:              fmt.Sprintf("projects/%s/instances/%s", projectID, instanceID),
		InstancePartitionId: instancePartitionID,
		InstancePartition: &instancepb.InstancePartition{
			Config:          fmt.Sprintf("projects/%s/instanceConfigs/%s", projectID, "nam3"),
			DisplayName:     "my-instance-partition",
			ComputeCapacity: &instancepb.InstancePartition_NodeCount{NodeCount: 1},
		},
	})
	if err != nil {
		return fmt.Errorf("could not create instance partition %s: %w", fmt.Sprintf("projects/%s/instances/%s/instancePartitions/%s", projectID, instanceID, instancePartitionID), err)
	}
	// Wait for the instance partition creation to finish.
	i, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("waiting for instance partition creation to finish failed: %w", err)
	}
	// The instance partition may not be ready to serve yet.
	if i.State != instancepb.InstancePartition_READY {
		fmt.Fprintf(w, "instance partition state is not READY yet. Got state %v\n", i.State)
	}
	fmt.Fprintf(w, "Created instance partition [%s]\n", instancePartitionID)
	return nil
}

// [END spanner_create_instance_partition]
