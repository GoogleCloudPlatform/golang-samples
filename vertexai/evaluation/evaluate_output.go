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

// [START generativeaionvertexai_evaluation_pointwise]
import (
	context_pkg "context"
	"fmt"
	"io"

	aiplatform "cloud.google.com/go/aiplatform/apiv1beta1"
	aiplatformpb "cloud.google.com/go/aiplatform/apiv1beta1/aiplatformpb"
	"google.golang.org/api/option"
)

// evaluateOutput evaluates an output of LLM using the groundedness metric
func evaluateOutput(w io.Writer, projectID, location string) error {
	ctx := context_pkg.Background()
	apiEndpoint := fmt.Sprintf("%s-aiplatform.googleapis.com:443", location)
	client, err := aiplatform.NewEvaluationClient(ctx, option.WithEndpoint(apiEndpoint))

	if err != nil {
		return fmt.Errorf("unable to create aiplatform client: %w", err)
	}
	defer client.Close()

	context := `
As part of a comprehensive initiative to tackle urban congestion and foster
sustainable urban living, a major city has revealed ambitious plans for an
extensive overhaul of its public transportation system. The project aims not
only to improve the efficiency and reliability of public transit but also to
reduce the city\'s carbon footprint and promote eco-friendly commuting options.
City officials anticipate that this strategic investment will enhance
accessibility for residents and visitors alike, ushering in a new era of
efficient, environmentally conscious urban transportation.
`
	modelResponse := `
The city is undertaking a major project to revamp its public transportation system.
This initiative is designed to improve efficiency, reduce carbon emissions, and promote
eco-friendly commuting. The city expects that this investment will enhance accessibility
and usher in a new era of sustainable urban transportation.
`
	req := aiplatformpb.EvaluateInstancesRequest{
		Location: fmt.Sprintf("projects/%s/locations/%s", projectID, location),
		MetricInputs: &aiplatformpb.EvaluateInstancesRequest_GroundednessInput{
			GroundednessInput: &aiplatformpb.GroundednessInput{
				MetricSpec: &aiplatformpb.GroundednessSpec{},
				Instance: &aiplatformpb.GroundednessInstance{
					Context:    &context,
					Prediction: &modelResponse,
				},
			},
		},
	}

	resp, err := client.EvaluateInstances(ctx, &req)
	if err != nil {
		return fmt.Errorf("evaluateInstances failed: %v", err)
	}

	results := resp.GetGroundednessResult()
	fmt.Fprintf(w, "score: %.2f\n", results.GetScore())
	fmt.Fprintf(w, "confidence: %.2f\n", results.GetConfidence())
	fmt.Fprintf(w, "explanation:\n%s\n", results.GetExplanation())
	// Example response:
	// score: 1.00
	// confidence: 1.00
	// explanation:
	// STEP 1: All aspects of the response are found in the context.
	// The response accurately summarizes the city's plan to overhaul its public transportation system, highlighting the goals of ...
	// STEP 2: According to the rubric, the response is scored 1 because all aspects of the response are attributable to the context.

	return nil
}

// [END generativeaionvertexai_evaluation_pointwise]
