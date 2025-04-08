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

// [START compute_usage_report_get]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/protobuf/proto"
)

// getUsageExportBucket retrieves the Compute Engine usage export bucket for the Cloud project. Replaces the empty value returned by the API with the default value used to generate report file names.
func getUsageExportBucket(w io.Writer, projectID string) error {
	// projectID := "your_project_id"
	ctx := context.Background()
	projectsClient, err := compute.NewProjectsRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewProjectsRESTClient: %w", err)
	}
	defer projectsClient.Close()

	// Updating the setting with an empty UsageExportLocationResource value disables the usage report generation.
	req := &computepb.GetProjectRequest{
		Project: projectID,
	}

	project, err := projectsClient.Get(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to get project: %w", err)
	}

	usageExportLocation := project.GetUsageExportLocation()

	if usageExportLocation == nil || usageExportLocation.GetBucketName() == "" {
		// The usage reports are disabled.
		return nil
	}

	if usageExportLocation.ReportNamePrefix != nil {
		// Although the server explicitly sent the empty string value,
		// the next usage report generated with these settings still has the default prefix value `usage_gce`.
		// (see https://cloud.google.com/compute/docs/reference/rest/v1/projects/get)
		fmt.Fprintf(w, "Report name prefix not set, replacing with default value of `usage_gce`.\n")
		usageExportLocation.ReportNamePrefix = proto.String("usage_gce")
	}

	fmt.Fprintf(w, "Returned ReportNamePrefix: %s\n", usageExportLocation.GetReportNamePrefix())

	return nil
}

// [END compute_usage_report_get]
