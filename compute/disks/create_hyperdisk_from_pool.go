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
	"time"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/protobuf/proto"
)

// [START compute_hyperdisk_create_from_pool]

// createDiskInStoragePool creates a new Hyperdisk in the specified storage pool.
func createDiskInStoragePool(w io.Writer, projectId, zone, diskName, storagePoolName, diskType string, diskSizeGb, provisionedIops, provisionedThroughput int64) error {
	// Example usage:
	//   projectID := "your_project_id"
	//   zone := "europe-central2-b"
	//   diskName := "your_disk_name"
	//   storagePoolName := "https://www.googleapis.com/compute/v1/projects/your_project_id/zones/europe-central2-b/storagePools/your_storage_pool"
	//   diskType := "zones/europe-central2-b/diskTypes/hyperdisk-balanced"
	//   diskSizeGb := int64(10)
	//   provisionedIops := int64(3000)
	//   provisionedThroughput := int64(140)

	ctx := context.Background()
	client, err := compute.NewDisksRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewDisksRESTClient: %v", err)
	}
	defer client.Close()

	// Create the disk resource
	disk := &computepb.Disk{
		Name:                  proto.String(diskName),
		Type:                  proto.String(diskType),
		SizeGb:                proto.Int64(diskSizeGb),
		Zone:                  proto.String(zone),
		StoragePool:           proto.String(storagePoolName),
		ProvisionedIops:       proto.Int64(provisionedIops),
		ProvisionedThroughput: proto.Int64(provisionedThroughput),
	}

	// Create the insert disk request
	req := &computepb.InsertDiskRequest{
		Project:      projectId,
		Zone:         zone,
		DiskResource: disk,
	}

	// Send the insert disk request
	op, err := client.Insert(ctx, req)
	if err != nil {
		return fmt.Errorf("Insert disk request failed: %v", err)
	}

	// Wait for the insert disk operation to complete
	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %w", err)
	}

	// Wait for server update
	time.Sleep(10 * time.Second)

	// Retrieve and return the created disk
	createdDisk, err := client.Get(ctx, &computepb.GetDiskRequest{
		Project: projectId,
		Zone:    zone,
		Disk:    diskName,
	})
	if err != nil {
		return fmt.Errorf("Get disk request failed: %v", err)
	}

	fmt.Fprintf(w, "Disk created in storage pool: %v\n", createdDisk.GetName())
	return nil
}

// [END compute_hyperdisk_create_from_pool]
