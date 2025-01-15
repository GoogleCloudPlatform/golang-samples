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

// [START compute_regional_disk_create_from_snapshot]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/protobuf/proto"
)

// createRegionalDiskFromSnapshot creates a new regional disk with the contents of
// an already existitng zonal disk snapshot.
func createRegionalDiskFromSnapshot(
	w io.Writer,
	projectID, region string, replicaZones []string,
	diskName, diskType, snapshotLink string,
	diskSizeGb int64,
) error {
	// projectID := "your_project_id"
	// region := "us-west3" // should match diskType below
	// replicaZones := []string{"us-west3-a", "us-west3-b"}
	// diskName := "your_disk_name"
	// diskType := "regions/us-west3/diskTypes/pd-ssd"
	// snapshotLink := "projects/your_project_id/global/snapshots/snapshot_name"
	// diskSizeGb := 120

	// Exactly two replica zones must be specified
	replicaZoneURLs := []string{
		fmt.Sprintf("projects/%s/zones/%s", projectID, replicaZones[0]),
		fmt.Sprintf("projects/%s/zones/%s", projectID, replicaZones[1]),
	}

	ctx := context.Background()
	disksClient, err := compute.NewRegionDisksRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewRegionDisksRESTClient: %w", err)
	}
	defer disksClient.Close()

	req := &computepb.InsertRegionDiskRequest{
		Project: projectID,
		Region:  region,
		DiskResource: &computepb.Disk{
			Name:           proto.String(diskName),
			Region:         proto.String(region),
			Type:           proto.String(diskType),
			SourceSnapshot: proto.String(snapshotLink),
			SizeGb:         proto.Int64(diskSizeGb),
			ReplicaZones:   replicaZoneURLs,
		},
	}

	op, err := disksClient.Insert(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to create disk: %w", err)
	}

	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %w", err)
	}

	fmt.Fprintf(w, "Disk created\n")

	return nil
}

// [END compute_regional_disk_create_from_snapshot]
