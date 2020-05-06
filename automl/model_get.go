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

// [START automl_get_model]
import (
	"context"
	"fmt"
	"io"

	automl "cloud.google.com/go/automl/apiv1"
	automlpb "google.golang.org/genproto/googleapis/cloud/automl/v1"
)

// getModel gets a model.
func getModel(w io.Writer, projectID string, location string, modelID string) error {
	// projectID := "my-project-id"
	// location := "us-central1"
	// modelID := "TRL123456789..."

	ctx := context.Background()
	client, err := automl.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewClient: %v", err)
	}
	defer client.Close()

	req := &automlpb.GetModelRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/models/%s", projectID, location, modelID),
	}

	model, err := client.GetModel(ctx, req)
	if err != nil {
		return fmt.Errorf("GetModel: %v", err)
	}

	// Retrieve deployment state.
	deploymentState := "undeployed"
	if model.GetDeploymentState() == automlpb.Model_DEPLOYED {
		deploymentState = "deployed"
	}

	// Display the model information.
	fmt.Fprintf(w, "Model name: %v\n", model.GetName())
	fmt.Fprintf(w, "Model display name: %v\n", model.GetDisplayName())
	fmt.Fprintf(w, "Model create time:\n")
	fmt.Fprintf(w, "\tseconds: %v\n", model.GetCreateTime().GetSeconds())
	fmt.Fprintf(w, "\tnanos: %v\n", model.GetCreateTime().GetNanos())
	fmt.Fprintf(w, "Model deployment state: %v\n", deploymentState)

	return nil
}

// [END automl_get_model]
