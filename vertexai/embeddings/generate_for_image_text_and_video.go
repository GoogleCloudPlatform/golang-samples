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

// [START generativeaionvertexai_multimodal_embedding_image_video_text]
import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	aiplatform "cloud.google.com/go/aiplatform/apiv1beta1"
	aiplatformpb "cloud.google.com/go/aiplatform/apiv1beta1/aiplatformpb"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

// generateForImageTextAndVideo shows how use the multimodal model to generate embeddings for
// text and image data.
func generateForImageTextAndVideo(w io.Writer, project, location string) error {
	// location = "us-central1"

	// The default context timeout may be not enough to process a video input
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	apiEndpoint := fmt.Sprintf("%s-aiplatform.googleapis.com:443", location)
	client, err := aiplatform.NewPredictionClient(ctx, option.WithEndpoint(apiEndpoint))
	if err != nil {
		return fmt.Errorf("failed to construct API client: %w", err)
	}
	defer client.Close()

	model := "multimodalembedding@001"
	endpoint := fmt.Sprintf("projects/%s/locations/%s/publishers/google/models/%s", project, location, model)

	// This is the input to the model's prediction call. The schema of any single instance
	// will be specified by the endpoint's deployed model:
	// https://storage.googleapis.com/google-cloud-aiplatform/schema/predict/instance/vision_embedding_model_1.0.0.yaml
	instance, err := structpb.NewValue(map[string]any{
		"text": "Domestic cats in natural conditions",
		"image": map[string]any{
			// Image and video inputs can be provided either as a Google Cloud Storage URI or as
			// base64-encoded bytes using the "bytesBase64Encoded" field
			"gcsUri": "gs://cloud-samples-data/generative-ai/image/320px-Felis_catus-cat_on_snow.jpg",
		},
		"video": map[string]any{
			"gcsUri": "gs://cloud-samples-data/video/cat.mp4",
		},
	})
	if err != nil {
		return fmt.Errorf("failed to construct request payload: %w", err)
	}

	req := &aiplatformpb.PredictRequest{
		Endpoint: endpoint,
		// The model supports only 1 instance per request
		Instances: []*structpb.Value{instance},
	}

	resp, err := client.Predict(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to generate embeddings: %w", err)
	}

	instanceEmbeddingsJson, err := protojson.Marshal(resp.GetPredictions()[0])
	if err != nil {
		return fmt.Errorf("failed to convert protobuf value to JSON: %w", err)
	}
	// Check the response schema of the model:
	// https://cloud.google.com/vertex-ai/generative-ai/docs/model-reference/multimodal-embeddings-api#response-body
	var instanceEmbeddings struct {
		ImageEmbeddings []float32 `json:"imageEmbedding"`
		TextEmbeddings  []float32 `json:"textEmbedding"`
		VideoEmbeddings []struct {
			Embedding      []float32 `json:"embedding"`
			StartOffsetSec float64   `json:"startOffsetSec"`
			EndOffsetSec   float64   `json:"endOffsetSec"`
		} `json:"videoEmbeddings"`
	}
	if err := json.Unmarshal(instanceEmbeddingsJson, &instanceEmbeddings); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	imageEmbedding := instanceEmbeddings.ImageEmbeddings
	textEmbedding := instanceEmbeddings.TextEmbeddings
	// Get the embedding for our single video segment (`.videoEmbeddings` object has one entry per
	// each processed segment)
	videoEmbedding := instanceEmbeddings.VideoEmbeddings[0].Embedding

	fmt.Fprintf(w, "Image embedding (length=%d): %v\n", len(imageEmbedding), imageEmbedding)
	fmt.Fprintf(w, "Text embedding (length=%d): %v\n", len(textEmbedding), textEmbedding)
	fmt.Fprintf(w, "Video embedding (length=%d): %v\n", len(videoEmbedding), videoEmbedding)
	// Example response:
	// Image embedding (length=1408): [-0.01558477 0.0258355 0.016342038 ... ]
	// Text embedding (length=1408): [-0.005894961 0.008349559 0.015355394 ... ]
	// Video embedding (length=1408): [-0.018867437 0.013997682 0.0012682161 ... ]

	return nil
}

// [END generativeaionvertexai_multimodal_embedding_image_video_text]
