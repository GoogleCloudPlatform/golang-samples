// Copyright 2025 Google LLC
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

package dataset

// [START bigquery_update_dataset_preview]
import (
	"context"
	"fmt"
	"io"

	bigquery "cloud.google.com/go/bigquery/apiv2"
	"cloud.google.com/go/bigquery/apiv2/bigquerypb"
	"github.com/googleapis/gax-go/v2/apierror"
	"github.com/googleapis/gax-go/v2/callctx"

	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// updateDataset demonstrates making partial updates to an existing dataset.
func updateDataset(w io.Writer, projectID, datasetID string) error {
	// projectID := "my-project-id"
	// datasetID := "mydataset"
	ctx := context.Background()

	// Construct a gRPC-based client.
	// To construct a REST-based client, use NewDatasetRESTClient instead.
	dsClient, err := bigquery.NewDatasetClient(ctx)
	if err != nil {
		return fmt.Errorf("bigquery.NewDatasetClient: %w", err)
	}
	defer dsClient.Close()

	// Fetch the existing dataset prior to making any modifications.
	// This allows us to use optimistic concurrency controls to avoid overwriting
	// other changes.
	meta, err := dsClient.GetDataset(ctx, &bigquerypb.GetDatasetRequest{
		ProjectId: projectID,
		DatasetId: datasetID,
		// Only fetch dataset metadata, not the permisions.
		DatasetView: bigquerypb.GetDatasetRequest_METADATA,
	})
	if err != nil {
		return fmt.Errorf("GetDataset: %w", err)
	}

	// Construct an update request, populating many of the available configurations.
	req := &bigquerypb.UpdateOrPatchDatasetRequest{
		ProjectId: projectID,
		DatasetId: datasetID,
		Dataset: &bigquerypb.Dataset{
			Description: &wrapperspb.StringValue{
				Value: "An updated description of the dataset.",
			},
			// Changing DefaultTableExpirationMs does not affect existing tables in the
			// dataset, but newly created tables will inherit this as defaults.
			DefaultTableExpirationMs: &wrapperspb.Int64Value{
				Value: 90 * 86400 * 1000, // 90 days in milliseconds.
			},
		},
	}
	// Now, use the ETag from the original metadata to guard against conflicting writes.
	// The callctx package let's us inject headers in a transport agnostic fashion (gRPC or HTTP).
	patchCtx := callctx.SetHeaders(ctx, "If-Match", meta.GetEtag())
	resp, err := dsClient.PatchDataset(patchCtx, req)
	if err != nil {
		if apierr, ok := apierror.FromError(err); ok {
			if status := apierr.GRPCStatus(); status.Code() == codes.FailedPrecondition {
				// The error was due to precondition failing (the If-Match constraint).
				// For this example we're not doing anything overly stateful with the dataset
				// so we simply return a more readable outer error.
				return fmt.Errorf("dataset etag changed between Get and Patch: %w", err)
			}
		}
		return fmt.Errorf("PatchDataset: %w", err)
	}
	// Print the values we expected to be modified.
	fmt.Fprintf(w, "Description: %s\n", resp.GetDescription())
	if expiration := resp.GetDefaultTableExpirationMs(); expiration != nil {
		fmt.Fprintf(w, "DefaultTableExpirationMs: %d", expiration.Value)
	}
	return nil
}

// [END bigquery_update_dataset_preview]
