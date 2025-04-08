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

// [START automl_language_entity_extraction_get_dataset]
// [START automl_language_sentiment_analysis_get_dataset]
// [START automl_language_text_classification_get_dataset]
// [START automl_translate_get_dataset]
// [START automl_vision_classification_get_dataset]
// [START automl_vision_object_detection_get_dataset]
import (
	"context"
	"fmt"
	"io"

	automl "cloud.google.com/go/automl/apiv1"
	"cloud.google.com/go/automl/apiv1/automlpb"
)

// getDataset gets a dataset.
func getDataset(w io.Writer, projectID string, location string, datasetID string) error {
	// projectID := "my-project-id"
	// location := "us-central1"
	// datasetID := "TRL123456789..."

	ctx := context.Background()
	client, err := automl.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewClient: %w", err)
	}
	defer client.Close()

	req := &automlpb.GetDatasetRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/datasets/%s", projectID, location, datasetID),
	}

	dataset, err := client.GetDataset(ctx, req)
	if err != nil {
		return fmt.Errorf("DeleteDataset: %w", err)
	}

	fmt.Fprintf(w, "Dataset name: %v\n", dataset.GetName())
	fmt.Fprintf(w, "Dataset display name: %v\n", dataset.GetDisplayName())
	fmt.Fprintf(w, "Dataset create time:\n")
	fmt.Fprintf(w, "\tseconds: %v\n", dataset.GetCreateTime().GetSeconds())
	fmt.Fprintf(w, "\tnanos: %v\n", dataset.GetCreateTime().GetNanos())

	// [END automl_language_sentiment_analysis_get_dataset]
	// [END automl_language_text_classification_get_dataset]
	// [END automl_translate_get_dataset]
	// [END automl_vision_classification_get_dataset]
	// [END automl_vision_object_detection_get_dataset]
	// Language entity extraction
	if metadata := dataset.GetTextExtractionDatasetMetadata(); metadata != nil {
		fmt.Fprintf(w, "Text extraction dataset metadata: %v\n", metadata)
	}
	// [END automl_language_entity_extraction_get_dataset]

	// [START automl_language_sentiment_analysis_get_dataset]
	// Language sentiment analysis
	if metadata := dataset.GetTextSentimentDatasetMetadata(); metadata != nil {
		fmt.Fprintf(w, "Text sentiment dataset metadata: %v\n", metadata)
	}
	// [END automl_language_sentiment_analysis_get_dataset]

	// [START automl_language_text_classification_get_dataset]
	// Language text classification
	if metadata := dataset.GetTextClassificationDatasetMetadata(); metadata != nil {
		fmt.Fprintf(w, "Text classification dataset metadata: %v\n", metadata)
	}
	// [END automl_language_text_classification_get_dataset]

	// [START automl_translate_get_dataset]
	// Translate
	if metadata := dataset.GetTranslationDatasetMetadata(); metadata != nil {
		fmt.Fprintf(w, "Translation dataset metadata:\n")
		fmt.Fprintf(w, "\tsource_language_code: %v\n", metadata.GetSourceLanguageCode())
		fmt.Fprintf(w, "\ttarget_language_code: %v\n", metadata.GetTargetLanguageCode())
	}
	// [END automl_translate_get_dataset]

	// [START automl_vision_classification_get_dataset]
	// Vision classification
	if metadata := dataset.GetImageClassificationDatasetMetadata(); metadata != nil {
		fmt.Fprintf(w, "Image classification dataset metadata: %v\n", metadata)
	}
	// [END automl_vision_classification_get_dataset]

	// [START automl_vision_object_detection_get_dataset]
	// Vision object detection
	if metadata := dataset.GetImageObjectDetectionDatasetMetadata(); metadata != nil {
		fmt.Fprintf(w, "Image object detection dataset metadata: %v\n", metadata)
	}
	// [START automl_language_entity_extraction_get_dataset]
	// [START automl_language_sentiment_analysis_get_dataset]
	// [START automl_language_text_classification_get_dataset]
	// [START automl_translate_get_dataset]
	// [START automl_vision_classification_get_dataset]

	return nil
}

// [END automl_language_entity_extraction_get_dataset]
// [END automl_language_sentiment_analysis_get_dataset]
// [END automl_language_text_classification_get_dataset]
// [END automl_translate_get_dataset]
// [END automl_vision_classification_get_dataset]
// [END automl_vision_object_detection_get_dataset]
