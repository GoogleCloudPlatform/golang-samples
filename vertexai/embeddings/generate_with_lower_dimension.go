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

// Package embeddings shows examples of working with multimodal embeddings in Vertex AI
package embeddings

// [START generativeaionvertexai_embeddings_specify_lower_dimension]
import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	aiplatform "cloud.google.com/go/aiplatform/apiv1beta1"
	aiplatformpb "cloud.google.com/go/aiplatform/apiv1beta1/aiplatformpb"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

// generateWithLowerDimension shows how to generate lower-dimensional embeddings for text and image inputs.
func generateWithLowerDimension(w io.Writer, project, location string) error {
	// location = "us-central1"
	ctx := context.Background()
	apiEndpoint := fmt.Sprintf("%s-aiplatform.googleapis.com:443", location)
	client, err := aiplatform.NewPredictionClient(ctx, option.WithEndpoint(apiEndpoint))
	if err != nil {
		return fmt.Errorf("failed to construct API client: %w", err)
	}
	defer client.Close()

	model := "multimodalembedding@001"
	endpoint := fmt.Sprintf("projects/%s/locations/%s/publishers/google/models/%s", project, location, model)

	// This is the input to the model's prediction call. For schema, see:
	// https://cloud.google.com/vertex-ai/generative-ai/docs/model-reference/multimodal-embeddings-api#request_body
	instance, err := structpb.NewValue(map[string]any{
		"image": map[string]any{
			// Image input can be provided either as a Google Cloud Storage URI or as
			// base64-encoded bytes using the "bytesBase64Encoded" field
			"gcsUri": "gs://cloud-samples-data/vertex-ai/llm/prompts/landmark1.png",
		},
		"text": "Colosseum",
	})
	if err != nil {
		return fmt.Errorf("failed to construct request payload: %w", err)
	}

	// TODO(developer): Try different dimenions: 128, 256, 512, 1408
	outputDimensionality := 128
	params, err := structpb.NewValue(map[string]any{
		"dimension": outputDimensionality,
	})
	if err != nil {
		return fmt.Errorf("failed to construct request params: %w", err)
	}

	req := &aiplatformpb.PredictRequest{
		Endpoint: endpoint,
		// The model supports only 1 instance per request
		Instances:  []*structpb.Value{instance},
		Parameters: params,
	}

	resp, err := client.Predict(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to generate embeddings: %w", err)
	}

	instanceEmbeddingsJson, err := protojson.Marshal(resp.GetPredictions()[0])
	if err != nil {
		return fmt.Errorf("failed to convert protobuf value to JSON: %w", err)
	}
	// For response schema, see:
	// https://cloud.google.com/vertex-ai/generative-ai/docs/model-reference/multimodal-embeddings-api#response-body
	var instanceEmbeddings struct {
		ImageEmbeddings []float32 `json:"imageEmbedding"`
		TextEmbeddings  []float32 `json:"textEmbedding"`
	}
	if err := json.Unmarshal(instanceEmbeddingsJson, &instanceEmbeddings); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	imageEmbedding := instanceEmbeddings.ImageEmbeddings
	textEmbedding := instanceEmbeddings.TextEmbeddings

	fmt.Fprintf(w, "Text embedding (length=%d): %v\n", len(textEmbedding), textEmbedding)
	fmt.Fprintf(w, "Image embedding (length=%d): %v\n", len(imageEmbedding), imageEmbedding)
	// Example response:
	// Text Embedding (length=128): [0.27469793 -0.14625867 0.022280363 ... ]
	// Image Embedding (length=128): [0.06225733 -0.040650766 0.02604402 ... ]

	return nil
}

// [END generativeaionvertexai_embeddings_specify_lower_dimension]
