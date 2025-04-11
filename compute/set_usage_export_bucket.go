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

// [START compute_usage_report_set]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/protobuf/proto"
)

// setUsageExportBucket sets the Compute Engine usage export bucket for the Cloud project.
// This sample presents how to interpret the default value for the report name prefix parameter.
func setUsageExportBucket(w io.Writer, projectID, bucketName, reportNamePrefix string) error {
	// projectID := "your_project_id"
	// bucketName := "your_bucket_name"
	// reportNamePrefix := ""
	ctx := context.Background()
	projectsClient, err := compute.NewProjectsRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewProjectsRESTClient: %w", err)
	}
	defer projectsClient.Close()

	// Updating the setting with an empty UsageExportLocationResource value disables the usage report generation.
	req := &computepb.SetUsageExportBucketProjectRequest{
		Project: projectID,
		UsageExportLocationResource: &computepb.UsageExportLocation{
			BucketName:       proto.String(bucketName),
			ReportNamePrefix: proto.String(reportNamePrefix),
		},
	}

	if reportNamePrefix == "" {
		// Sending an empty value for reportNamePrefix results in the next usage report being generated
		// with the default prefix value "usage_gce". (see: https://cloud.google.com/compute/docs/reference/rest/v1/projects/get)
		fmt.Fprintf(w, "Setting reportNamePrefix to empty value causes the report to have the default prefix value `usage_gce`.\n")
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

// [END compute_usage_report_set]
