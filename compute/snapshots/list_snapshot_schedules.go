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

// [START compute_snapshot_schedule_list]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/api/iterator"
	"google.golang.org/protobuf/proto"
)

// listSnapshotSchedule retrieves a list of snapshot schedules.
func listSnapshotSchedule(w io.Writer, projectID, region, filter string) error {
	// projectID := "your_project_id"
	// snapshotName := "your_snapshot_name"
	// region := "eupore-central2"

	// Formatting for filters:
	// https://cloud.google.com/python/docs/reference/compute/latest/google.cloud.compute_v1.types.ListResourcePoliciesRequest

	ctx := context.Background()

	snapshotsClient, err := compute.NewResourcePoliciesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewResourcePoliciesRESTClient: %w", err)
	}
	defer snapshotsClient.Close()

	req := &computepb.ListResourcePoliciesRequest{
		Project: projectID,
		Region:  region,
		Filter:  proto.String(filter),
	}
	it := snapshotsClient.List(ctx, req)

	for {
		policy, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "- %s", policy.GetName())
	}
	return nil
}

// [END compute_snapshot_schedule_list]
