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
	"fmt"
	"io"

	aiplatform "cloud.google.com/go/aiplatform/apiv1beta1"
	aiplatformpb "cloud.google.com/go/aiplatform/apiv1beta1/aiplatformpb"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/types/known/structpb"
)

// generateWithLowerDimension shows how to generate lower-dimensional embeddings for
// multimodal (text and image) inputs
func generateWithLowerDimension(w io.Writer, project, location string) ([][]float32, error) {
	// location = "us-central1"
	ctx := context.Background()
	apiEndpoint := fmt.Sprintf("%s-aiplatform.googleapis.com:443", location)
	client, err := aiplatform.NewPredictionClient(ctx, option.WithEndpoint(apiEndpoint))
	if err != nil {
		return nil, err
	}
	defer client.Close()

	model := "multimodalembedding@001"
	endpoint := fmt.Sprintf("projects/%s/locations/%s/publishers/google/models/%s", project, location, model)

	// This is the input to the model's prediction call. The schema of any single instance
	// will be specified by the endpoint's deployed model:
	// https://storage.googleapis.com/google-cloud-aiplatform/schema/predict/instance/vision_embedding_model_1.0.0.yaml
	instances := []*structpb.Value{
		structpb.NewStructValue(&structpb.Struct{
			Fields: map[string]*structpb.Value{
				"text": structpb.NewStringValue("Colosseum"),
				// Image input can be provided either as a Google Cloud Storage URI or as base64-encoded
				// bytes using the "bytesBase64Encoded" field
				"image": structpb.NewStructValue(&structpb.Struct{
					Fields: map[string]*structpb.Value{
						"gcsUri": structpb.NewStringValue("gs://cloud-samples-data/vertex-ai/llm/prompts/landmark1.png"),
					},
				}),
			},
		}),
	}

	// TODO(developer): Try different dimenions: 128, 256, 512, 1408
	outputDimensionality := 128
	req := &aiplatformpb.PredictRequest{
		Endpoint:  endpoint,
		Instances: instances,
		Parameters: structpb.NewStructValue(&structpb.Struct{
			Fields: map[string]*structpb.Value{
				"dimension": structpb.NewNumberValue(float64(outputDimensionality)),
			},
		}),
	}

	resp, err := client.Predict(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to generate embeddings: %v", err)
	}
	// The list of response predictions contains one prediction per input instance.
	// In this case, we sent only one input instance, so we access its prediction at
	// index 0.
	instanceEmbeddings := resp.GetPredictions()[0].GetStructValue().GetFields()
	textEmbedding := make([]float32, outputDimensionality)
	for i, v := range instanceEmbeddings["textEmbedding"].GetListValue().GetValues() {
		textEmbedding[i] = float32(v.GetNumberValue())
	}

	imageEmbedding := make([]float32, outputDimensionality)
	for i, v := range instanceEmbeddings["imageEmbedding"].GetListValue().GetValues() {
		imageEmbedding[i] = float32(v.GetNumberValue())
	}

	fmt.Fprintf(w, "Text embedding (length=%d): %v\n", len(textEmbedding), textEmbedding)
	fmt.Fprintf(w, "Image embedding (length=%d): %v\n", len(imageEmbedding), imageEmbedding)
	// Example response:
	// Text Embedding (length=128): [0.27469793 -0.14625867 0.022280363 ... ]
	// Image Embedding (length=128): [0.06225733 -0.040650766 0.02604402 ... ]

	return [][]float32{textEmbedding, imageEmbedding}, nil
}

// [END generativeaionvertexai_embeddings_specify_lower_dimension]
