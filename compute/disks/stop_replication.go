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

// [START compute_disk_stop_replication]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
)

// stopReplication stops primary disk replication in a project for a given zone.
func stopReplication(
	w io.Writer,
	projectID, primaryDiskName, primaryZone string,
) error {
	// projectID := "your_project_id"
	// primaryDiskName := "your_disk_name2"
	// primaryZone := "europe-west2-b"

	ctx := context.Background()
	disksClient, err := compute.NewDisksRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewDisksRESTClient: %w", err)
	}
	defer disksClient.Close()

	req := &computepb.StopAsyncReplicationDiskRequest{
		Project: projectID,
		Zone:    primaryZone,
		Disk:    primaryDiskName,
	}

	op, err := disksClient.StopAsyncReplication(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to create disk: %w", err)
	}

	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %w", err)
	}

	fmt.Fprintf(w, "Replication stopped\n")

	return nil
}

// [END compute_disk_stop_replication]
