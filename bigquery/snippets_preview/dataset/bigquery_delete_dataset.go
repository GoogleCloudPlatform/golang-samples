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

// [START bigquery_delete_dataset_preview]
import (
	"context"
	"fmt"

	bigquery "cloud.google.com/go/bigquery/apiv2"
	"cloud.google.com/go/bigquery/apiv2/bigquerypb"
	"github.com/googleapis/gax-go/v2/apierror"

	"google.golang.org/grpc/codes"
)

// deleteDataset demonstrates deleting a dataset from BigQuery.
func deleteDataset(projectID, datasetID string) error {
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

	req := &bigquerypb.DeleteDatasetRequest{
		ProjectId: projectID,
		DatasetId: datasetID,
		// Deletion will fail if the dataset is not empty and DeleteContents is false.
		DeleteContents: true,
	}

	// Deleting a dataset doesn't return information, but it may produce an error.
	err = dsClient.DeleteDataset(ctx, req)
	if err != nil {
		if apierr, ok := apierror.FromError(err); ok {
			if status := apierr.GRPCStatus(); status.Code() == codes.NotFound {
				// The error indicates the dataset isn't present.  Possibly another process removed
				// the dataset, or perhaps there was a partial failure and this was handled via automatic retry.
				// In any case, treat this as a success.
				return nil
			}
		}
		return fmt.Errorf("PatchDataset: %w", err)
	}
	return nil
}

// [END bigquery_delete_dataset_preview]
