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

// generateForVideo shows how to use the multimodal model to generate embeddings for
// video input
func generateForVideo(w io.Writer, project, location string) ([]float32, error) {
	// location = "us-central1"

	// The default context timeout may be not enough to process a video input
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	apiEndpoint := fmt.Sprintf("%s-aiplatform.googleapis.com:443", location)
	client, err := aiplatform.NewPredictionClient(ctx, option.WithEndpoint(apiEndpoint))
	if err != nil {
		return nil, fmt.Errorf("failed to construct API client: %v", err)
	}
	defer client.Close()

	model := "multimodalembedding@001"
	endpoint := fmt.Sprintf("projects/%s/locations/%s/publishers/google/models/%s", project, location, model)

	// This is the input to the model's prediction call. The schema of any single instance
	// may be specified by the endpoint's deployed model, e.g.:
	// https://storage.googleapis.com/google-cloud-aiplatform/schema/predict/instance/vision_embedding_model_1.0.0.yaml
	instances, err := structpb.NewValue(map[string]any{
		"video": map[string]any{
			"gcsUri": "gs://cloud-samples-data/vertex-ai-vision/highway_vehicles.mp4",
			"videoSegmentConfig": map[string]any{
				"startOffsetSec": 1,
				"endOffsetSec":   5,
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to construct request payload: %v", err)
	}

	req := &aiplatformpb.PredictRequest{
		Endpoint:  endpoint,
		Instances: []*structpb.Value{instances},  // The model supports only 1 instance per request
	}
	resp, err := client.Predict(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to generate embeddings: %v", err)
	}

	instanceEmbeddingsJson, err := protojson.Marshal(resp.GetPredictions()[0])
	if err != nil {
		return nil, fmt.Errorf("failed to convert protobuf value to JSON: %v", err)
	}
	// Check the response schema of the model:
	// https://storage.googleapis.com/google-cloud-aiplatform/schema/predict/prediction/vision_embedding_model_1.0.0.yaml
	var instanceEmbeddings struct {
		VideoEmbeddings []struct {
			Embedding      []float32 `json:"embedding"`
			StartOffsetSec float64   `json:"startOffsetSec"`
			EndOffsetSec   float64   `json:"endOffsetSec"`
		} `json:"videoEmbeddings"`
	}
	if err := json.Unmarshal(instanceEmbeddingsJson, &instanceEmbeddings); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %v", err)
	}
	// Get the embedding for our single video interval (videoEmbeddings has one entry per interval)
	videoEmbedding := instanceEmbeddings.VideoEmbeddings[0]

	fmt.Fprintf(w, "Video embedding (seconds: %.f-%.f; length=%d): %v\n",
		videoEmbedding.StartOffsetSec,
		videoEmbedding.EndOffsetSec,
		len(videoEmbedding.Embedding),
		videoEmbedding.Embedding,
	)
	// Example response:
	// Video embedding (seconds: 1-5; length=1408): [-0.016427778 0.032878537 -0.030755188 ... ]

	return videoEmbedding.Embedding, nil
}

// [END generativeaionvertexai_multimodal_embedding_video]
