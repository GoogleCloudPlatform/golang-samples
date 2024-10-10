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

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/aiplatform/apiv1/aiplatformpb"
	gax "github.com/googleapis/gax-go/v2"
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	GPUEndpointID     = "123456789"
	GPUEndpointRegion = "us-east1"
	TPUEndpointID     = "987654321"
	TPUEndpointRegion = "us-west1"
)

type PredictionsClient struct{}

func (client PredictionsClient) Close() error {
	return nil
}

func (client PredictionsClient) Predict(ctx context.Context, req *aiplatformpb.PredictRequest, opts ...gax.CallOption) (*aiplatformpb.PredictResponse, error) {
	mockedResponse := `
	The sky appears blue due to a phenomenon called **Rayleigh scattering**.

	**Here's how it works:**

	* **Sunlight is white:**  Sunlight actually contains all the colors of the rainbow.

	* **Scattering:** When sunlight enters the Earth's atmosphere, it collides with tiny gas molecules (mostly nitrogen and oxygen). These collisions cause the light to scatter in different directions.

	* **Blue light scatters most:**  Blue light has a shorter wavelength
	`
	response := &aiplatformpb.PredictResponse{
		Predictions: []*structpb.Value{structpb.NewStringValue(mockedResponse)},
	}

	instance := req.Instances[0].GetStructValue()
	if ok := strings.Contains(req.Endpoint, fmt.Sprintf("locations/%s/endpoints/%s", GPUEndpointRegion, GPUEndpointID)); ok {
		if err := instance.Fields["inputs"].GetStringValue(); err == "" {
			return nil, fmt.Errorf("invalid request")
		}
	} else if ok := strings.Contains(req.Endpoint, fmt.Sprintf("locations/%s/endpoints/%s", TPUEndpointRegion, TPUEndpointID)); ok {
		if err := instance.Fields["prompt"].GetStringValue(); err == "" {
			return nil, fmt.Errorf("invalid request")
		}
	}

	params := req.Instances[0].GetStructValue().Fields["parameters"].GetStructValue()

	if err := params.Fields["temperature"].GetNumberValue(); err == 0 {
		return nil, fmt.Errorf("invalid request")
	}
	if err := params.Fields["maxOutputTokens"].GetNumberValue(); err == 0 {
		return nil, fmt.Errorf("invalid request")
	}
	if err := params.Fields["topP"].GetNumberValue(); err == 0 {
		return nil, fmt.Errorf("invalid request")
	}
	if err := params.Fields["topK"].GetNumberValue(); err == 0 {
		return nil, fmt.Errorf("invalid request")
	}
	return response, nil
}
