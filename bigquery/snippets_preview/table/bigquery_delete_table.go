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

// [START bigquery_delete_table_preview]
import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery/v2/apiv2/bigquerypb"
	"cloud.google.com/go/bigquery/v2/apiv2_client"
	"github.com/googleapis/gax-go/v2/apierror"

	"google.golang.org/grpc/codes"
)

// deleteTable demonstrates deleting a table from BigQuery.
func deleteTable(client *apiv2_client.Client, projectID, datasetID, tableID string) error {
	// client can be instantiated per-RPC service, or use cloud.google.com/go/bigquery/v2/apiv2_client to create
	// an aggregate client.
	//
	// projectID := "my-project-id"
	// datasetID := "mydataset"
	// tableID := "mytable"
	ctx := context.Background()

	req := &bigquerypb.DeleteTableRequest{
		ProjectId: projectID,
		DatasetId: datasetID,
		TableId:   tableID,
	}

	// Deleting a table doesn't return information, but it may produce an error.
	if err := client.DeleteTable(ctx, req); err != nil {
		if apierr, ok := apierror.FromError(err); ok {
			if status := apierr.GRPCStatus(); status.Code() == codes.NotFound {
				// The error indicates the table isn't present.  Possibly another process removed
				// the table, or perhaps there was a partial failure and this was handled via automatic retry.
				// In any case, treat this as a success.
				return nil
			}
		}
		return fmt.Errorf("DeleteTable: %w", err)
	}
	return nil
}

// [END bigquery_delete_table_preview]
