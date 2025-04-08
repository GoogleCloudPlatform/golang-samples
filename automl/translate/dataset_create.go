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

// [START automl_translate_create_dataset]
import (
	"context"
	"fmt"
	"io"

	automl "cloud.google.com/go/automl/apiv1"
	"cloud.google.com/go/automl/apiv1/automlpb"
)

// translateCreateDataset creates a dataset for translate.
func translateCreateDataset(w io.Writer, projectID string, location string, datasetName string, sourceLanguageCode string, targetLanguageCode string) error {
	// projectID := "my-project-id"
	// location := "us-central1"
	// datasetName := "dataset_display_name"

	// Supported languages:
	//   https://cloud.google.com/translate/automl/docs/languages
	// sourceLanguageCode := "en"
	// targetLanguageCode := "ja"

	ctx := context.Background()
	client, err := automl.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewClient: %w", err)
	}
	defer client.Close()

	req := &automlpb.CreateDatasetRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, location),
		Dataset: &automlpb.Dataset{
			DisplayName: datasetName,
			DatasetMetadata: &automlpb.Dataset_TranslationDatasetMetadata{
				TranslationDatasetMetadata: &automlpb.TranslationDatasetMetadata{
					SourceLanguageCode: sourceLanguageCode,
					TargetLanguageCode: targetLanguageCode,
				},
			},
		},
	}

	op, err := client.CreateDataset(ctx, req)
	if err != nil {
		return fmt.Errorf("CreateDataset: %w", err)
	}
	fmt.Fprintf(w, "Processing operation name: %q\n", op.Name())

	dataset, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("Wait: %w", err)
	}

	fmt.Fprintf(w, "Dataset name: %v\n", dataset.GetName())

	return nil
}

// [END automl_translate_create_dataset]
