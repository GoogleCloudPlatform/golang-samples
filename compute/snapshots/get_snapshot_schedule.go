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

// [START compute_snapshot_schedule_get]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
)

// getSnapshotSchedule gets a snapshot schedule.
func getSnapshotSchedule(w io.Writer, projectID, scheduleName, region string) error {
	// projectID := "your_project_id"
	// snapshotName := "your_snapshot_name"
	// region := "eupore-central2"

	ctx := context.Background()

	snapshotsClient, err := compute.NewResourcePoliciesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewResourcePoliciesRESTClient: %w", err)
	}
	defer snapshotsClient.Close()

	req := &computepb.GetResourcePolicyRequest{
		Project:        projectID,
		Region:         region,
		ResourcePolicy: scheduleName,
	}
	schedule, err := snapshotsClient.Get(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to get snapshot schedule: %w", err)
	}

	fmt.Fprintf(w, "Found snapshot schedule: %s\n", schedule.GetName())

	return nil
}

// [END compute_snapshot_schedule_get]
