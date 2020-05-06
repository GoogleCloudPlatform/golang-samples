// Copyright 2019 Google LLC
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

// Package automl contains samples for Google Cloud AutoML API v1.
package automl

// [START automl_language_entity_extraction_predict]
import (
	"context"
	"fmt"
	"io"

	automl "cloud.google.com/go/automl/apiv1"
	automlpb "google.golang.org/genproto/googleapis/cloud/automl/v1"
)

// languageEntityExtractionPredict does a prediction for text entity extraction.
func languageEntityExtractionPredict(w io.Writer, projectID string, location string, modelID string, content string) error {
	// projectID := "my-project-id"
	// location := "us-central1"
	// modelID := "TEN123456789..."
	// content := "text to extract entities"

	ctx := context.Background()
	client, err := automl.NewPredictionClient(ctx)
	if err != nil {
		return fmt.Errorf("NewPredictionClient: %v", err)
	}
	defer client.Close()

	req := &automlpb.PredictRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/models/%s", projectID, location, modelID),
		Payload: &automlpb.ExamplePayload{
			Payload: &automlpb.ExamplePayload_TextSnippet{
				TextSnippet: &automlpb.TextSnippet{
					Content:  content,
					MimeType: "text/plain", // Types: "text/plain", "text/html"
				},
			},
		},
	}

	resp, err := client.Predict(ctx, req)
	if err != nil {
		return fmt.Errorf("Predict: %v", err)
	}

	for _, payload := range resp.GetPayload() {
		fmt.Fprintf(w, "Text extract entity types: %v\n", payload.GetDisplayName())
		fmt.Fprintf(w, "Text score: %v\n", payload.GetTextExtraction().GetScore())
		textSegment := payload.GetTextExtraction().GetTextSegment()
		fmt.Fprintf(w, "Text extract entity content: %v\n", textSegment.GetContent())
		fmt.Fprintf(w, "Text start offset: %v\n", textSegment.GetStartOffset())
		fmt.Fprintf(w, "Text end offset: %v\n", textSegment.GetEndOffset())
	}

	return nil
}

// [END automl_language_entity_extraction_predict]
