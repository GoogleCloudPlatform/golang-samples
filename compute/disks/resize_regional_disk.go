// Copyright 2023 Google LLC
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

// [START compute_regional_disk_resize]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
)

// resizeRegionalDisk changes the size of a regional disk.
// After you resize the disk, you must also resize the file system
// so that the operating system can access the additional space.
func resizeRegionalDisk(w io.Writer, projectID, region, diskName string, newSizeGb int64) error {
	// projectID := "your_project_id"
	// region := "us-west3"
	// diskName := "your_disk_name"
	// newSizeGb := 20

	ctx := context.Background()
	disksClient, err := compute.NewRegionDisksRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewRegionDisksRESTClient: %w", err)
	}
	defer disksClient.Close()

	req := &computepb.ResizeRegionDiskRequest{
		Disk:    diskName,
		Project: projectID,
		Region:  region,
		RegionDisksResizeRequestResource: &computepb.RegionDisksResizeRequest{
			SizeGb: &newSizeGb,
		},
	}

	op, err := disksClient.Resize(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to resize disk: %w", err)
	}

	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %w", err)
	}

	fmt.Fprintf(w, "Disk resized\n")

	return nil
}

// [END compute_regional_disk_resize]
