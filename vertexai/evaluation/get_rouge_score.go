// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// evaluation package shows examples of working with Gen AI Evaluation API
package evaluation

// [START generativeaionvertexai_evaluation_get_rouge_score]
import (
	"context"
	"fmt"
	"io"

	aiplatform "cloud.google.com/go/aiplatform/apiv1beta1"
	aiplatformpb "cloud.google.com/go/aiplatform/apiv1beta1/aiplatformpb"
	"google.golang.org/api/option"
)

// getROUGEScore evaluates a model response against a reference (ground truth) using the ROUGE metric
func getROUGEScore(w io.Writer, projectID, location string) error {
	// location = "us-central1"
	ctx := context.Background()
	apiEndpoint := fmt.Sprintf("%s-aiplatform.googleapis.com:443", location)
	client, err := aiplatform.NewEvaluationClient(ctx, option.WithEndpoint(apiEndpoint))

	if err != nil {
		return fmt.Errorf("unable to create aiplatform client: %w", err)
	}
	defer client.Close()

	modelResponse := `
The Great Barrier Reef, the world's largest coral reef system located in Australia,
is a vast and diverse ecosystem. However, it faces serious threats from climate change,
ocean acidification, and coral bleaching, endangering its rich marine life.
`
	reference := `
The Great Barrier Reef, the world's largest coral reef system, is
located off the coast of Queensland, Australia. It's a vast
ecosystem spanning over 2,300 kilometers with thousands of reefs
and islands. While it harbors an incredible diversity of marine
life, including endangered species, it faces serious threats from
climate change, ocean acidification, and coral bleaching.
`
	req := aiplatformpb.EvaluateInstancesRequest{
		Location: fmt.Sprintf("projects/%s/locations/%s", projectID, location),
		MetricInputs: &aiplatformpb.EvaluateInstancesRequest_RougeInput{
			RougeInput: &aiplatformpb.RougeInput{
				// Check the API reference for the list of supported ROUGE metric types:
				// https://cloud.google.com/vertex-ai/docs/reference/rpc/google.cloud.aiplatform.v1beta1#rougespec
				MetricSpec: &aiplatformpb.RougeSpec{
					RougeType: "rouge1",
				},
				Instances: []*aiplatformpb.RougeInstance{
					{
						Prediction: &modelResponse,
						Reference:  &reference,
					},
				},
			},
		},
	}

	resp, err := client.EvaluateInstances(ctx, &req)
	if err != nil {
		return fmt.Errorf("evaluateInstances failed: %v", err)
	}

	fmt.Fprintln(w, "evaluation results:")
	fmt.Fprintln(w, resp.GetRougeResults().GetRougeMetricValues())
	// Example response:
	// [score:0.6597938]

	return nil
}

// [END generativeaionvertexai_evaluation_get_rouge_score]
