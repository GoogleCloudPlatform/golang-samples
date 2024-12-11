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

// [START compute_snapshot_schedule_edit]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/protobuf/proto"
)

// editSnapshotSchedule edits a snapshot schedule.
func editSnapshotSchedule(w io.Writer, projectID, scheduleName, region string) error {
	// projectID := "your_project_id"
	// snapshotName := "your_snapshot_name"
	// region := "eupore-central2"

	ctx := context.Background()

	snapshotsClient, err := compute.NewResourcePoliciesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewResourcePoliciesRESTClient: %w", err)
	}
	defer snapshotsClient.Close()

	req := &computepb.PatchResourcePolicyRequest{
		Project:        projectID,
		Region:         region,
		ResourcePolicy: scheduleName,
		ResourcePolicyResource: &computepb.ResourcePolicy{
			Name:        proto.String(scheduleName),
			Description: proto.String("MY HOURLY SNAPSHOT SCHEDULE"),
			SnapshotSchedulePolicy: &computepb.ResourcePolicySnapshotSchedulePolicy{
				Schedule: &computepb.ResourcePolicySnapshotSchedulePolicySchedule{
					HourlySchedule: &computepb.ResourcePolicyHourlyCycle{
						HoursInCycle: proto.Int32(12),
						StartTime:    proto.String("22:00"),
					},
				},
			},
		},
	}
	op, err := snapshotsClient.Patch(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to create snapshot schedule: %w", err)
	}

	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %w", err)
	}

	fmt.Fprint(w, "Snapshot schedule changed\n")

	return nil
}

// [END compute_snapshot_schedule_edit]
