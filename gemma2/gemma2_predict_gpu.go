// Copyright 2024 Google LLC
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

package snippets

// [START generativeaionvertexai_gemma2_predict_gpu]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/aiplatform/apiv1/aiplatformpb"
	"github.com/googleapis/gax-go/v2"

	"google.golang.org/protobuf/types/known/structpb"
)

// predictGPU demonstrates how to run interference on a Gemma2 model deployed to a Vertex AI endpoint with GPU accelerators.
func predictGPU(w io.Writer, client ClientInterface, projectID, location, endpointID string) error {
	ctx := context.Background()

	// Note: client can be initialised in the following way:
	// apiEndpoint := fmt.Sprintf("%s-aiplatform.googleapis.com:443", location)
	// client, err := aiplatform.NewPredictionClient(ctx, option.WithEndpoint(apiEndpoint))
	// if err != nil {
	// 	return fmt.Errorf("unable to create prediction client: %v", err)
	// }
	// defer client.Close()

	gemma2Endpoint := fmt.Sprintf("projects/%s/locations/%s/endpoints/%s", projectID, location, endpointID)
	prompt := "Why is the sky blue?"
	parameters := map[string]interface{}{
		"temperature":     0.9,
		"maxOutputTokens": 1024,
		"topP":            1.0,
		"topK":            1,
	}

	// Encapsulate the prompt in a correct format for TPUs.
	// Pay attention that prompt should be set in "inputs" field.
	// Example format: [{'inputs': 'Why is the sky blue?', 'parameters': {'temperature': 0.9}}]
	promptValue, err := structpb.NewValue(map[string]interface{}{
		"inputs":     prompt,
		"parameters": parameters,
	})
	if err != nil {
		fmt.Fprintf(w, "unable to convert prompt to Value: %v", err)
		return err
	}

	req := &aiplatformpb.PredictRequest{
		Endpoint:  gemma2Endpoint,
		Instances: []*structpb.Value{promptValue},
	}

	resp, err := client.Predict(ctx, req)
	if err != nil {
		return err
	}

	prediction := resp.GetPredictions()
	value := prediction[0].GetStringValue()
	fmt.Fprintf(w, "%v", value)

	return nil
}

// [END generativeaionvertexai_gemma2_predict_gpu]

type ClientInterface interface {
	Close() error
	Predict(ctx context.Context, req *aiplatformpb.PredictRequest, opts ...gax.CallOption) (*aiplatformpb.PredictResponse, error)
}
