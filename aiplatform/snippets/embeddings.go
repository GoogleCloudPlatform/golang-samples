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

// [START aiplatform_text_embeddings]
// [START generativeaionvertexai_text_embeddings]
import (
	"context"
	"fmt"
	"io"

	aiplatform "cloud.google.com/go/aiplatform/apiv1"
	"cloud.google.com/go/aiplatform/apiv1/aiplatformpb"

	"google.golang.org/api/option"
	"google.golang.org/protobuf/types/known/structpb"
)

// embedTexts shows how embeddings are set for gemini-embedding-001 model
func embedTexts(w io.Writer, project, location string) error {
	// location := "us-central1"
	ctx := context.Background()

	apiEndpoint := fmt.Sprintf("%s-aiplatform.googleapis.com:443", location)
	dimensionality := 3072
	model := "gemini-embedding-001"
	texts := []string{"banana muffins? ", "banana bread? banana muffins?"}

	client, err := aiplatform.NewPredictionClient(ctx, option.WithEndpoint(apiEndpoint))
	if err != nil {
		return err
	}
	defer client.Close()

	endpoint := fmt.Sprintf("projects/%s/locations/%s/publishers/google/models/%s", project, location, model)
	allEmbeddings := make([][]float32, 0, len(texts))
	// gemini-embedding-001 takes 1 input at a time
	for i, text := range texts {
		instances := make([]*structpb.Value, 1)
		instances[0] = structpb.NewStructValue(&structpb.Struct{
			Fields: map[string]*structpb.Value{
				"content":   structpb.NewStringValue(text),
				"task_type": structpb.NewStringValue("QUESTION_ANSWERING"),
			},
		})

		params := structpb.NewStructValue(&structpb.Struct{
			Fields: map[string]*structpb.Value{
				"outputDimensionality": structpb.NewNumberValue(float64(dimensionality)),
			},
		})

		req := &aiplatformpb.PredictRequest{
			Endpoint:   endpoint,
			Instances:  instances,
			Parameters: params,
		}
		resp, err := client.Predict(ctx, req)
		if err != nil {
			return err
		}

		// Process the prediction for the single text
		// The response will contain one prediction because we sent one instance.
		if len(resp.Predictions) == 0 {
			return fmt.Errorf("no predictions returned for text \"%s\"", text)
		}

		prediction := resp.Predictions[0]
		embeddingValues := prediction.GetStructValue().Fields["embeddings"].GetStructValue().Fields["values"].GetListValue().Values

		currentEmbedding := make([]float32, len(embeddingValues))
		for j, value := range embeddingValues {
			currentEmbedding[j] = float32(value.GetNumberValue())
		}
		allEmbeddings = append(allEmbeddings, currentEmbedding)
	}

	if len(allEmbeddings) > 0 {
		fmt.Fprintf(w, "Dimensionality: %d. Embeddings length: %d", len(allEmbeddings[0]), len(allEmbeddings))
	} else {
		fmt.Fprintln(w, "No texts were processed.")
	}
	return nil
}

// [END aiplatform_text_embeddings]
// [END generativeaionvertexai_text_embeddings]
