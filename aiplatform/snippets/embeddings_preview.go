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

package snippets

// [START generativeaionvertexai_sdk_embedding]
import (
	"context"
	"fmt"
	"regexp"

	aiplatform "cloud.google.com/go/aiplatform/apiv1"
	"cloud.google.com/go/aiplatform/apiv1/aiplatformpb"

	"google.golang.org/api/option"
	"google.golang.org/protobuf/types/known/structpb"
)

// Embeds code query with a pre-trained, foundational model by specifying the task type as 'CODE_RETRIEVAL_QUERY'. e.g. 'Retrieve a function that adds two numbers'.
// Embeds code block with a pre-trained, foundational model by specifying the task type as 'RETRIEVAL_DOCUMENT'. e.g. 'texts := []string{"def func(a, b): return a + b", "def func(a, b): return a - b", "def func(a, b): return (a ** 2 + b ** 2) ** 0.5"}'.
func embedTextsPreview(
	apiEndpoint, project, model string, texts []string,
	task string, dimensionality *int) ([][]float32, error) {
	ctx := context.Background()

	client, err := aiplatform.NewPredictionClient(ctx, option.WithEndpoint(apiEndpoint))
	if err != nil {
		return nil, err
	}
	defer client.Close()

	match := regexp.MustCompile(`^(\w+-\w+)`).FindStringSubmatch(apiEndpoint)
	location := "us-central1"
	if match != nil {
		location = match[1]
	}
	endpoint := fmt.Sprintf("projects/%s/locations/%s/publishers/google/models/%s", project, location, model)
	instances := make([]*structpb.Value, len(texts))
	for i, text := range texts {
		instances[i] = structpb.NewStructValue(&structpb.Struct{
			Fields: map[string]*structpb.Value{
				"content":   structpb.NewStringValue(text),
				"task_type": structpb.NewStringValue(task),
			},
		})
	}
	outputDimensionality := structpb.NewNullValue()
	if dimensionality != nil {
		outputDimensionality = structpb.NewNumberValue(float64(*dimensionality))
	}
	params := structpb.NewStructValue(&structpb.Struct{
		Fields: map[string]*structpb.Value{"outputDimensionality": outputDimensionality},
	})

	req := &aiplatformpb.PredictRequest{
		Endpoint:   endpoint,
		Instances:  instances,
		Parameters: params,
	}
	resp, err := client.Predict(ctx, req)
	if err != nil {
		return nil, err
	}
	embeddings := make([][]float32, len(resp.Predictions))
	for i, prediction := range resp.Predictions {
		values := prediction.GetStructValue().Fields["embeddings"].GetStructValue().Fields["values"].GetListValue().Values
		embeddings[i] = make([]float32, len(values))
		for j, value := range values {
			embeddings[i][j] = float32(value.GetNumberValue())
		}
	}
	return embeddings, nil
}

// [END generativeaionvertexai_sdk_embedding]
