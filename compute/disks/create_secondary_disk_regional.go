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

// [START compute_disk_create_secondary_regional]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/protobuf/proto"
)

// createRegionalSecondaryDisk creates a new secondary disk in a project in given region.
// Note: secondary disk should be located in a different region, but within the same continent.
// More details: https://cloud.google.com/compute/docs/disks/async-pd/about#supported_region_pairs
func createRegionalSecondaryDisk(
	w io.Writer,
	projectID, region, diskName, primaryDiskName, primaryRegion string,
	replicaZones []string,
	diskSizeGb int64,
) error {
	// projectID := "your_project_id"
	// region := "europe-west1"
	// diskName := "your_disk_name"
	// primaryDiskName := "your_disk_name2"
	// primaryDiskRegion := "europe-west4"
	// replicaZones := []string{"europe-west1-a", "europe-west1-b"}
	// diskSizeGb := 200

	ctx := context.Background()
	disksClient, err := compute.NewRegionDisksRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewRegionDisksRESTClient: %w", err)
	}
	defer disksClient.Close()

	primaryFullDiskName := fmt.Sprintf("projects/%s/regions/%s/disks/%s", projectID, primaryRegion, primaryDiskName)

	// Exactly two replica zones must be specified
	replicaZoneURLs := []string{
		fmt.Sprintf("projects/%s/zones/%s", projectID, replicaZones[0]),
		fmt.Sprintf("projects/%s/zones/%s", projectID, replicaZones[1]),
	}

	req := &computepb.InsertRegionDiskRequest{
		Project: projectID,
		Region:  region,
		DiskResource: &computepb.Disk{
			Name:   proto.String(diskName),
			Region: proto.String(region),
			// The size must be at least 200 GB
			SizeGb: proto.Int64(diskSizeGb),
			AsyncPrimaryDisk: &computepb.DiskAsyncReplication{
				Disk: proto.String(primaryFullDiskName),
			},
			ReplicaZones: replicaZoneURLs,
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

// [END compute_disk_create_secondary_regional]
