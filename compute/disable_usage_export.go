// Copyright 2021 Google LLC
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

// [START compute_usage_report_disable]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
)

// disableUsageExport disables Compute Engine usage export bucket for the Cloud Project.
func disableUsageExport(w io.Writer, projectID string) error {
	// projectID := "your_project_id"
	ctx := context.Background()
	projectsClient, err := compute.NewProjectsRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewProjectsRESTClient: %w", err)
	}
	defer projectsClient.Close()

	// Providing an empty UsageExportLocationResource value disables the usage report generation.
	req := &computepb.SetUsageExportBucketProjectRequest{
		Project: projectID,
	}

	op, err := projectsClient.SetUsageExportBucket(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to set usage export bucket %w", err)
	}

	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %w", err)
	}

	fmt.Fprintf(w, "Usage export bucket has been set\n")

	return nil
}

// [END compute_usage_report_disable]
