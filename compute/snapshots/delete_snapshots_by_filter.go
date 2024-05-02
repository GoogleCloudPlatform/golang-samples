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

// [START compute_snapshot_delete_by_filter]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/api/iterator"
)

// Deletes ALL disk snapshots in the project that match the filter
func deleteByFilter(w io.Writer, projectID, filter string) error {
	// projectID := "your_project_id"
	// filter := ""
	// Learn more about filters:
	// https://cloud.google.com/python/docs/reference/compute/latest/google.cloud.compute_v1.types.ListSnapshotsRequest

	ctx := context.Background()
	snapshotsClient, err := compute.NewSnapshotsRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewSnapshotsRESTClient: %w", err)
	}
	defer snapshotsClient.Close()

	// List snapshots in the project with the filter applied
	req := &computepb.ListSnapshotsRequest{
		Project: projectID,
		Filter:  &filter,
	}
	it := snapshotsClient.List(ctx, req)

	// Iterate over the list of snapshots
	for {
		snapshot, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}

		// Delete each snapshot that matches the filter
		req := &computepb.DeleteSnapshotRequest{
			Project:  projectID,
			Snapshot: *snapshot.Name,
		}

		op, err := snapshotsClient.Delete(ctx, req)
		if err != nil {
			return fmt.Errorf("unable to delete snapshot: %w", err)
		}

		if err = op.Wait(ctx); err != nil {
			return fmt.Errorf("unable to wait for the operation: %w", err)
		}

		fmt.Fprintf(w, "Snapshot %s deleted\n", *snapshot.Name)
	}
	return nil
}

// [END compute_snapshot_delete_by_filter]
