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

package table

// [START bigquery_update_table_preview]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/bigquery/v2/apiv2/bigquerypb"
	"cloud.google.com/go/bigquery/v2/apiv2_client"
	"github.com/googleapis/gax-go/v2/apierror"
	"github.com/googleapis/gax-go/v2/callctx"

	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// updateTable demonstrates making partial updates to an existing table's metadata.
func updateTable(client *apiv2_client.Client, w io.Writer, projectID, datasetID, tableID string) error {
	// client can be instantiated per-RPC service, or use cloud.google.com/go/bigquery/v2/apiv2_client to create
	// an aggregate client.
	//
	// projectID := "my-project-id"
	// datasetID := "mydataset"
	ctx := context.Background()

	// Fetch the existing table metadata prior to making any modifications.
	// This allows us to use optimistic concurrency controls to avoid overwriting
	// other changes.
	meta, err := client.GetTable(ctx, &bigquerypb.GetTableRequest{
		ProjectId: projectID,
		DatasetId: datasetID,
		TableId:   tableID,
	})
	if err != nil {
		return fmt.Errorf("GetTable: %w", err)
	}

	// Construct an update request, populating many of the available configurations.
	req := &bigquerypb.UpdateOrPatchTableRequest{
		ProjectId: projectID,
		DatasetId: datasetID,
		TableId:   tableID,
		Table: &bigquerypb.Table{
			Description: &wrapperspb.StringValue{
				Value: "An updated description of the dataset.",
			},
		},
	}
	// Now, use the ETag from the original metadata to guard against conflicting writes.
	// The callctx package let's us inject headers in a transport agnostic fashion (gRPC or HTTP).
	patchCtx := callctx.SetHeaders(ctx, "if-match", meta.GetEtag())
	resp, err := client.PatchTable(patchCtx, req)
	if err != nil {
		if apierr, ok := apierror.FromError(err); ok {
			if status := apierr.GRPCStatus(); status.Code() == codes.FailedPrecondition {
				// The error was due to precondition failing (the If-Match constraint).
				// For this example we're not doing anything overly stateful with the dataset
				// so we simply return a more readable outer error.
				return fmt.Errorf("table etag changed between Get and Patch: %w", err)
			}
		}
		return fmt.Errorf("PatchTable: %w", err)
	}
	// Print the values we expected to be modified.
	fmt.Fprintf(w, "Description: %s\n", resp.GetDescription())
	return nil
}

// [END bigquery_update_table_preview]
