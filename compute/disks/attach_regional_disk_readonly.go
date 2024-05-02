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

// [START compute_regional_disk_attach_read_only]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
)

// Attaches the provided regional disk in read-only mode to the given VM.
// Read-only disks can be attached to multiple VMs at once.
func attachRegionalDiskReadOnly(w io.Writer, projectID, zone, instanceName, diskUrl string) error {
	// projectID := "your_project_id"
	// zone := "us-west3-a" // refers to the instance, not the disk
	// instanceName := "your_instance_name"
	// diskUrl := "projects/your_project/regions/europe-west3/disks/your_disk"

	ctx := context.Background()
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return err
	}
	defer instancesClient.Close()

	mode := "READ_ONLY"

	req := &computepb.AttachDiskInstanceRequest{
		AttachedDiskResource: &computepb.AttachedDisk{
			Source: &diskUrl,
			Mode:   &mode,
		},
		Instance: instanceName,
		Project:  projectID,
		Zone:     zone,
	}

	op, err := instancesClient.AttachDisk(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to attach disk: %w", err)
	}

	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %w", err)
	}

	fmt.Fprintf(w, "Disk attached\n")

	return nil
}

// [END compute_regional_disk_attach_read_only]
