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

// [START compute_snapshot_delete]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
)

// deleteSnapshot deletes a snapshot of a disk.
func deleteSnapshot(w io.Writer, projectID, snapshotName string) error {
	// projectID := "your_project_id"
	// snapshotName := "your_snapshot_name"

	ctx := context.Background()
	snapshotsClient, err := compute.NewSnapshotsRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewSnapshotsRESTClient: %w", err)
	}
	defer snapshotsClient.Close()

	req := &computepb.DeleteSnapshotRequest{
		Project:  projectID,
		Snapshot: snapshotName,
	}

	op, err := snapshotsClient.Delete(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to delete snapshot: %w", err)
	}

	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %w", err)
	}

	fmt.Fprintf(w, "Snapshot deleted\n")

	return nil
}

// [END compute_snapshot_delete]
