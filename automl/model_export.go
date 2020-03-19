// Copyright 2019 Google LLC
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

// Package automl contains samples for Google Cloud AutoML API v1.
package automl

// [START automl_export_model]
import (
	"context"
	"fmt"
	"io"

	automl "cloud.google.com/go/automl/apiv1"
	automlpb "google.golang.org/genproto/googleapis/cloud/automl/v1"
)

// exportModel exports a model.
func exportModel(w io.Writer, projectID string, location string, modelID string, outputURI string) error {
	// projectID := "my-project-id"
	// location := "us-central1"
	// modelID := "TRL123456789..."
	// outputURI := "gs://BUCKET_ID/path_to_export/"

	ctx := context.Background()
	client, err := automl.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewClient: %v", err)
	}
	defer client.Close()

	req := &automlpb.ExportModelRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/models/%s", projectID, location, modelID),
		OutputConfig: &automlpb.ModelExportOutputConfig{
			Destination: &automlpb.ModelExportOutputConfig_GcsDestination{
				GcsDestination: &automlpb.GcsDestination{
					OutputUriPrefix: outputURI,
				},
			},
		},
	}

	op, err := client.ExportModel(ctx, req)
	if err != nil {
		return fmt.Errorf("ExportModel: %v", err)
	}
	fmt.Fprintf(w, "Processing operation name: %q\n", op.Name())

	if err := op.Wait(ctx); err != nil {
		return fmt.Errorf("Wait: %v", err)
	}

	fmt.Fprintf(w, "Model exported.\n")

	return nil
}

// [END automl_export_model]
