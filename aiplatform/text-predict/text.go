// Copyright 2023 Google LLC
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

package snippets

// [START aiplatform_text_predictions]
// [START generativeaionvertexai_text_predictions]

import (
	"context"
	"fmt"
	"io"

	aiplatform "cloud.google.com/go/aiplatform/apiv1beta1"
	"cloud.google.com/go/aiplatform/apiv1beta1/aiplatformpb"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/types/known/structpb"
)

// textPredict generates text with certain prompt and configurations.
func textPredict(w io.Writer, projectID, location, model string) error {
	ctx := context.Background()

	prompt := "Hello, say something nice."
	publisher := "google"
	parameters := map[string]interface{}{
		"temperature":     0.8,
		"maxOutputTokens": 256,
		"topP":            0.4,
		"topK":            40,
	}

	apiEndpoint := fmt.Sprintf("%s-aiplatform.googleapis.com:443", location)

	client, err := aiplatform.NewPredictionClient(ctx, option.WithEndpoint(apiEndpoint))
	if err != nil {
		fmt.Fprintf(w, "unable to create prediction client: %v", err)
		return err
	}
	defer client.Close()

	// PredictRequest requires an endpoint, instances, and parameters
	// Endpoint
	base := fmt.Sprintf("projects/%s/locations/%s/publishers/%s/models", projectID, location, publisher)
	url := fmt.Sprintf("%s/%s", base, model)

	// Instances: the prompt to use with the text model
	promptValue, err := structpb.NewValue(map[string]interface{}{
		"prompt": prompt,
	})
	if err != nil {
		fmt.Fprintf(w, "unable to convert prompt to Value: %v", err)
		return err
	}

	// Parameters: the model configuration parameters
	parametersValue, err := structpb.NewValue(parameters)
	if err != nil {
		fmt.Fprintf(w, "unable to convert parameters to Value: %v", err)
		return err
	}

	// PredictRequest: create the model prediction request
	req := &aiplatformpb.PredictRequest{
		Endpoint:   url,
		Instances:  []*structpb.Value{promptValue},
		Parameters: parametersValue,
	}

	// PredictResponse: receive the response from the model
	resp, err := client.Predict(ctx, req)
	if err != nil {
		fmt.Fprintf(w, "error in prediction: %v", err)
		return err
	}

	fmt.Fprintf(w, "text-prediction response: %v", resp.Predictions[0])
	return nil
}

// [END aiplatform_text_predictions]
// [END generativeaionvertexai_text_predictions]
