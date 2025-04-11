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

// [START compute_hyperdisk_pool_create]

// createHyperdiskStoragePool creates a new Hyperdisk storage pool in the specified project and zone.
func createHyperdiskStoragePool(w io.Writer, projectId, zone, storagePoolName, storagePoolType string) error {
	// projectID := "your_project_id"
	// zone := "europe-west4-b"
	// storagePoolName := "your_storage_pool_name"
	// storagePoolType := "projects/**your_project_id**/zones/europe-west4-b/diskTypes/hyperdisk-balanced"

	ctx := context.Background()
	client, err := compute.NewStoragePoolsRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewStoragePoolsRESTClient: %v", err)
	}
	defer client.Close()

	// Create the storage pool resource
	resource := &computepb.StoragePool{
		Name:                        proto.String(storagePoolName),
		Zone:                        proto.String(zone),
		StoragePoolType:             proto.String(storagePoolType),
		CapacityProvisioningType:    proto.String("advanced"),
		PerformanceProvisioningType: proto.String("advanced"),
		PoolProvisionedCapacityGb:   proto.Int64(10240),
		PoolProvisionedIops:         proto.Int64(10000),
		PoolProvisionedThroughput:   proto.Int64(1024),
	}

	// Create the insert storage pool request
	req := &computepb.InsertStoragePoolRequest{
		Project:             projectId,
		Zone:                zone,
		StoragePoolResource: resource,
	}

	// Send the insert storage pool request
	op, err := client.Insert(ctx, req)
	if err != nil {
		return fmt.Errorf("Insert storage pool request failed: %v", err)
	}

	// Wait for the insert storage pool operation to complete
	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %w", err)
	}

	// Retrieve and return the created storage pool
	storagePool, err := client.Get(ctx, &computepb.GetStoragePoolRequest{
		Project:     projectId,
		Zone:        zone,
		StoragePool: storagePoolName,
	})
	if err != nil {
		return fmt.Errorf("Get storage pool request failed: %v", err)
	}

	fmt.Fprintf(w, "Hyperdisk Storage Pool created: %v\n", storagePool.GetName())
	return nil
}

// [END compute_hyperdisk_pool_create]
