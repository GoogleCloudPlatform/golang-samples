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

// [START generativeaionvertexai_multimodal_embedding_video]
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

// generateForVideo shows how to use the multimodal model to generate embeddings for video input.
func generateForVideo(w io.Writer, project, location string) error {
	// location = "us-central1"

	// The default context timeout may be not enough to process a video input.
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

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
	instances, err := structpb.NewValue(map[string]any{
		"video": map[string]any{
			// Video input can be provided either as a Google Cloud Storage URI or as base64-encoded
			// bytes using the "bytesBase64Encoded" field.
			"gcsUri": "gs://cloud-samples-data/vertex-ai-vision/highway_vehicles.mp4",
			"videoSegmentConfig": map[string]any{
				"startOffsetSec": 1,
				"endOffsetSec":   5,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to construct request payload: %w", err)
	}

	req := &aiplatformpb.PredictRequest{
		Endpoint: endpoint,
		// The model supports only 1 instance per request.
		Instances: []*structpb.Value{instances},
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
		VideoEmbeddings []struct {
			Embedding      []float32 `json:"embedding"`
			StartOffsetSec float64   `json:"startOffsetSec"`
			EndOffsetSec   float64   `json:"endOffsetSec"`
		} `json:"videoEmbeddings"`
	}
	if err := json.Unmarshal(instanceEmbeddingsJson, &instanceEmbeddings); err != nil {
		return fmt.Errorf("failed to unmarshal json: %w", err)
	}
	// Get the embedding for our single video segment (`.videoEmbeddings` object has one entry per
	// each processed segment).
	videoEmbedding := instanceEmbeddings.VideoEmbeddings[0]

	fmt.Fprintf(w, "Video embedding (seconds: %.f-%.f; length=%d): %v\n",
		videoEmbedding.StartOffsetSec,
		videoEmbedding.EndOffsetSec,
		len(videoEmbedding.Embedding),
		videoEmbedding.Embedding,
	)
	// Example response:
	// Video embedding (seconds: 1-5; length=1408): [-0.016427778 0.032878537 -0.030755188 ... ]

	return nil
}

// [END generativeaionvertexai_multimodal_embedding_video]
