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

package model

// [START bigquery_get_model_preview]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/bigquery/v2/apiv2/bigquerypb"
	"cloud.google.com/go/bigquery/v2/apiv2_client"
	"google.golang.org/protobuf/encoding/protojson"
)

// getModel demonstrates fetching ML model information.
func getModel(client *apiv2_client.Client, w io.Writer, projectID, datasetID, modelID string) error {
	// client can be instantiated per-RPC service, or use cloud.google.com/go/bigquery/v2/apiv2_client to create
	// an aggregate client.
	//
	// projectID := "my-project-id"
	// datasetID := "mydataset"
	// modelID := "mymodel"
	ctx := context.Background()

	req := &bigquerypb.GetModelRequest{
		ProjectId: projectID,
		DatasetId: datasetID,
		ModelId:   modelID,
	}

	resp, err := client.GetModel(ctx, req)
	if err != nil {
		return fmt.Errorf("GetModel: %w", err)
	}

	// Print some of the information about the model to the provided writer.
	fmt.Fprintf(w, "Model %q has description %q\n",
		resp.GetModelReference().GetModelId(),
		resp.GetDescription())
	// Alternately, use the protojson package to print a more complete representation
	// of the model using a basic JSON mapping:
	fmt.Fprintf(w, "Model JSON representation:\n%s\n", protojson.Format(resp))
	return nil
}

// [END bigquery_get_model_preview]
