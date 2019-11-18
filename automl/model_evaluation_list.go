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

// Package automl contains samples for Google Cloud AutoML API v1beta1.
package automl

// [START automl_list_model_evaluations]
import (
	"context"
	"fmt"
	"io"

	automl "cloud.google.com/go/automl/apiv1beta1"
	"google.golang.org/api/iterator"
	automlpb "google.golang.org/genproto/googleapis/cloud/automl/v1beta1"
)

// listModelEvaluation lists existing model evaluations.
func listModelEvaluations(w io.Writer, projectID string, location string, modelID string) error {
	// projectID := "my-project-id"
	// location := "us-central1"
	// modelID := "TRL123456789..."

	ctx := context.Background()
	client, err := automl.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewClient: %v", err)
	}
	defer client.Close()

	req := &automlpb.ListModelEvaluationsRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s/models/%s", projectID, location, modelID),
	}

	it := client.ListModelEvaluations(ctx, req)

	// Iterate over all results
	for {
		evaluation, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("ListModelEvaluations.Next: %v", err)
		}

		fmt.Fprintf(w, "Model evaluation name: %v\n", evaluation.GetName())
		fmt.Fprintf(w, "Model annotation spec id: %v\n", evaluation.GetAnnotationSpecId())
		fmt.Fprintf(w, "Create Time:\n")
		fmt.Fprintf(w, "\tseconds: %v\n", evaluation.GetCreateTime().GetSeconds())
		fmt.Fprintf(w, "\tnanos: %v\n", evaluation.GetCreateTime().GetNanos())
		fmt.Fprintf(w, "Evaluation example count: %v\n", evaluation.GetEvaluatedExampleCount())
		fmt.Fprintf(w, "Model evaluation metrics: %v\n", evaluation.GetTranslationEvaluationMetrics())
	}

	return nil
}

// [END automl_list_model_evaluations]
