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

// [START compute_snapshot_schedule_remove]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
)

// removeSnapshotSchedule removes a snapshot schedule from a disk.
func removeSnapshotSchedule(w io.Writer, projectID, scheduleName, diskName, region string) error {
	// projectID := "your_project_id"
	// snapshotName := "your_snapshot_name"
	// diskName := "your_disk_name"
	// region := "europe-central2"

	ctx := context.Background()

	disksClient, err := compute.NewRegionDisksRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewDisksRESTClient: %w", err)
	}
	defer disksClient.Close()

	req := &computepb.RemoveResourcePoliciesRegionDiskRequest{
		Project: projectID,
		Region:  region,
		Disk:    diskName,
		RegionDisksRemoveResourcePoliciesRequestResource: &computepb.RegionDisksRemoveResourcePoliciesRequest{
			ResourcePolicies: []string{
				fmt.Sprintf("projects/%s/regions/%s/resourcePolicies/%s", projectID, region, scheduleName),
			},
		},
	}
	op, err := disksClient.RemoveResourcePolicies(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to remove schedule: %w", err)
	}

	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %w", err)
	}

	fmt.Fprint(w, "Snapshot schedule removed\n")

	return nil
}

// [END compute_snapshot_schedule_remove]
