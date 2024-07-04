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

import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/protobuf/proto"
)

// [START compute_hyperdisk_create]
// createHyperdisk creates a new Hyperdisk in the specified project and zone.
func createHyperdisk(w io.Writer, projectId, zone, diskName string) error {
	//   projectID := "your_project_id"
	//   zone := "europe-central2-b"
	//   diskName := "your_disk_name"

	ctx := context.Background()
	client, err := compute.NewDisksRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewDisksRESTClient: %v", err)
	}
	defer client.Close()

	// use format "zones/{zone}/diskTypes/(hyperdisk-balanced|hyperdisk-throughput)".
	diskType := fmt.Sprintf("zones/%s/diskTypes/hyperdisk-balanced", zone)

	// Create the disk
	disk := &computepb.Disk{
		Name:   proto.String(diskName),
		Type:   proto.String(diskType),
		SizeGb: proto.Int64(10),
		Zone:   proto.String(zone),
	}

	req := &computepb.InsertDiskRequest{
		Project:      projectId,
		Zone:         zone,
		DiskResource: disk,
	}

	op, err := client.Insert(ctx, req)
	if err != nil {
		return fmt.Errorf("Insert disk request failed: %v", err)
	}

	// Wait for the insert disk operation to complete
	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %w", err)
	}

	fmt.Fprintf(w, "Hyperdisk created: %v\n", diskName)
	return nil
}

// [END compute_hyperdisk_create]
