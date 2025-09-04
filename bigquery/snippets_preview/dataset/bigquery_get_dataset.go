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

// [START bigquery_get_dataset_preview]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/bigquery/v2/apiv2/bigquerypb"
	"cloud.google.com/go/bigquery/v2/apiv2_client"
	"google.golang.org/protobuf/encoding/protojson"
)

// getDataset demonstrates fetching dataset information.
func getDataset(client *apiv2_client.Client, w io.Writer, projectID, datasetID string) error {
	// client can be instantiated per-RPC service, or use cloud.google.com/go/bigquery/v2/apiv2_client to create
	// an aggregate client.
	//
	// projectID := "my-project-id"
	// datasetID := "mydataset"
	ctx := context.Background()

	req := &bigquerypb.GetDatasetRequest{
		ProjectId: projectID,
		DatasetId: datasetID,
		// Dataset supports fetching a subset of the dataset information depending
		// on whether you're interested in security information, basic metadata, or
		// both.  For the example, we'll request all the information.
		DatasetView: bigquerypb.GetDatasetRequest_FULL,
	}

	resp, err := client.GetDataset(ctx, req)
	if err != nil {
		return fmt.Errorf("GetDataset: %w", err)
	}

	// Print some of the information about the dataset to the provided writer.
	fmt.Fprintf(w, "Dataset %q has description %q\n",
		resp.GetDatasetReference().GetDatasetId(),
		resp.GetDescription())
	// Alternately, use the protojson package to print a more complete representation
	// of the dataset using a basic JSON mapping:
	fmt.Fprintf(w, "Dataset JSON representation:\n%s\n", protojson.Format(resp))
	return nil
}

// [END bigquery_get_dataset_preview]
