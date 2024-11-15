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
	"fmt"
	"io"
	"time"

	aiplatform "cloud.google.com/go/aiplatform/apiv1beta1"
	aiplatformpb "cloud.google.com/go/aiplatform/apiv1beta1/aiplatformpb"
	"google.golang.org/api/option"
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
				// Video input can be provided either as a Google Cloud Storage URI or as base64-encoded
				// bytes using the "bytesBase64Encoded" field
				"video": structpb.NewStructValue(&structpb.Struct{
					Fields: map[string]*structpb.Value{
						"gcsUri": structpb.NewStringValue("gs://cloud-samples-data/vertex-ai-vision/highway_vehicles.mp4"),
						"videoSegmentConfig": structpb.NewStructValue(&structpb.Struct{
							Fields: map[string]*structpb.Value{
								"startOffsetSec": structpb.NewNumberValue(1),
								"endOffsetSec":   structpb.NewNumberValue(5),
							},
						}),
					},
				}),
			},
		}),
	}

	req := &aiplatformpb.PredictRequest{
		Endpoint:  endpoint,
		Instances: instances,
	}

	resp, err := client.Predict(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to generate embeddings: %v", err)
	}
	// The list of response predictions contains one prediction per input instance.
	// In this case, we sent only one input instance, so we access its prediction at
	// index 0. Check the response schema of the model for more details:
	// https://storage.googleapis.com/google-cloud-aiplatform/schema/predict/prediction/vision_embedding_model_1.0.0.yaml
	instanceEmbeddings := resp.GetPredictions()[0].GetStructValue().GetFields()
	videoEmbeddingsList := instanceEmbeddings["videoEmbeddings"].GetListValue().GetValues()
	// The list of "videoEmbeddings" contains one embedding per processed video interval.
	// In this case, our entire video input should be processed as one interval, so we
	// access it at index 0
	segmentData := videoEmbeddingsList[0].GetStructValue().GetFields()

	// By default, an embedding request returns a 1408-dimensional vector of float values
	videoEmbedding := make([]float32, 1408)
	for i, v := range segmentData["embedding"].GetListValue().GetValues() {
		videoEmbedding[i] = float32(v.GetNumberValue())
	}
	videoSegment := fmt.Sprintf(
		"seconds: %.f-%.f",
		segmentData["startOffsetSec"].GetNumberValue(),
		segmentData["endOffsetSec"].GetNumberValue(),
	)

	fmt.Fprintf(w, "Video embedding (%s; length=%d): %v\n", videoSegment, len(videoEmbedding), videoEmbedding)
	// Example response:
	// Video embedding (seconds: 1-5; length=1408): [-0.016427778 0.032878537 -0.030755188 ... ]

	return videoEmbedding, nil
}

// [END generativeaionvertexai_multimodal_embedding_video]
