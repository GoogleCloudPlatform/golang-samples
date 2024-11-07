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

// [START generativeaionvertexai_evaluation_pairwise_summarization_quality]
import (
	context_pkg "context"
	"fmt"
	"io"

	aiplatform_v1beta "cloud.google.com/go/aiplatform/apiv1beta1"
	aiplatformpb_v1beta "cloud.google.com/go/aiplatform/apiv1beta1/aiplatformpb"
	"google.golang.org/api/option"
)

// pairwiseEvaluation lets the judge model to compare the responses of two models and pick the better one
func pairwiseEvaluation(w io.Writer, projectID, location string) error {
	ctx := context_pkg.Background()
	apiEndpoint := fmt.Sprintf("%s-aiplatform.googleapis.com:443", location)
	client, err := aiplatform_v1beta.NewEvaluationClient(ctx, option.WithEndpoint(apiEndpoint))

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
	instruction := "Summarize the text such that a five-year-old can understand."
	baselineResponse := `
The city wants to make it easier for people to get around without using cars.
They're going to make the buses and trains better and faster, so people will want to
use them more. This will help the air be cleaner and make the city a better place to live.
`
	candidateResponse := `
The city is making big changes to how people get around. They want to make the buses and
trains work better and be easier for everyone to use. This will also help the environment
by getting people to use less gas. The city thinks these changes will make it easier for
everyone to get where they need to go.
`

	req := aiplatformpb_v1beta.EvaluateInstancesRequest{
		Location: fmt.Sprintf("projects/%s/locations/%s", projectID, location),
		MetricInputs: &aiplatformpb_v1beta.EvaluateInstancesRequest_PairwiseSummarizationQualityInput{
			PairwiseSummarizationQualityInput: &aiplatformpb_v1beta.PairwiseSummarizationQualityInput{
				MetricSpec: &aiplatformpb_v1beta.PairwiseSummarizationQualitySpec{},
				Instance: &aiplatformpb_v1beta.PairwiseSummarizationQualityInstance{
					Context:            &context,
					Instruction:        &instruction,
					Prediction:         &candidateResponse,
					BaselinePrediction: &baselineResponse,
				},
			},
		},
	}

	resp, err := client.EvaluateInstances(ctx, &req)
	if err != nil {
		return fmt.Errorf("evaluateInstances failed: %v", err)
	}

	results := resp.GetPairwiseSummarizationQualityResult()
	fmt.Fprintf(w, "choice: %s\n", results.GetPairwiseChoice())
	fmt.Fprintf(w, "confidence: %.2f\n", results.GetConfidence())
	fmt.Fprintf(w, "explanation:\n%s\n", results.GetExplanation())
	// Example response:
	// choice: BASELINE
	// confidence: 0.50
	// explanation:
	// BASELINE response is easier to understand. For example, the phrase "..." is easier to understand than "...". Thus, BASELINE response is ...

	return nil
}

// [END generativeaionvertexai_evaluation_pairwise_summarization_quality]
