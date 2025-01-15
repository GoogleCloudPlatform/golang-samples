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

package snippets

// [START compute_snapshot_create]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/protobuf/proto"
)

// createSnapshot creates a snapshot of a disk.
func createSnapshot(
	w io.Writer,
	projectID, diskName, snapshotName, zone, region, location, diskProjectID string,
) error {
	// projectID := "your_project_id"
	// diskName := "your_disk_name"
	// snapshotName := "your_snapshot_name"
	// zone := "europe-central2-b"
	// region := "eupore-central2"
	// location = "eupore-central2"
	// diskProjectID = "YOUR_DISK_PROJECT_ID"

	ctx := context.Background()

	snapshotsClient, err := compute.NewSnapshotsRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewSnapshotsRESTClient: %w", err)
	}
	defer snapshotsClient.Close()

	if zone == "" && region == "" {
		return fmt.Errorf("you need to specify `zone` or `region` for this function to work")
	}

	if zone != "" && region != "" {
		return fmt.Errorf("you can't set both `zone` and `region` parameters")
	}

	if diskProjectID == "" {
		diskProjectID = projectID
	}

	disk := &computepb.Disk{}
	locations := []string{}
	if location != "" {
		locations = append(locations, location)
	}

	if zone != "" {
		disksClient, err := compute.NewDisksRESTClient(ctx)
		if err != nil {
			return fmt.Errorf("NewDisksRESTClient: %w", err)
		}
		defer disksClient.Close()

		getDiskReq := &computepb.GetDiskRequest{
			Project: projectID,
			Zone:    zone,
			Disk:    diskName,
		}

		disk, err = disksClient.Get(ctx, getDiskReq)
		if err != nil {
			return fmt.Errorf("unable to get disk: %w", err)
		}
	} else {
		regionDisksClient, err := compute.NewRegionDisksRESTClient(ctx)
		if err != nil {
			return fmt.Errorf("NewRegionDisksRESTClient: %w", err)
		}
		defer regionDisksClient.Close()

		getDiskReq := &computepb.GetRegionDiskRequest{
			Project: projectID,
			Region:  region,
			Disk:    diskName,
		}

		disk, err = regionDisksClient.Get(ctx, getDiskReq)
		if err != nil {
			return fmt.Errorf("unable to get disk: %w", err)
		}
	}

	req := &computepb.InsertSnapshotRequest{
		Project: projectID,
		SnapshotResource: &computepb.Snapshot{
			Name:             proto.String(snapshotName),
			SourceDisk:       proto.String(disk.GetSelfLink()),
			StorageLocations: locations,
		},
	}

	op, err := snapshotsClient.Insert(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to create snapshot: %w", err)
	}

	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %w", err)
	}

	fmt.Fprintf(w, "Snapshot created\n")

	return nil
}

// [END compute_snapshot_create]
