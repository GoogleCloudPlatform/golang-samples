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

// Package automl contains samples for Google Cloud AutoML API v1beta1.
package automl

// [START automl_list_datasets]
import (
	"context"
	"fmt"
	"io"

	automl "cloud.google.com/go/automl/apiv1beta1"
	"google.golang.org/api/iterator"
	automlpb "google.golang.org/genproto/googleapis/cloud/automl/v1beta1"
)

// listDatasets lists existing datasets.
func listDatasets(w io.Writer, projectID string, location string) error {
	// projectID := "my-project-id"
	// location := "us-central1"

	ctx := context.Background()
	client, err := automl.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewClient: %v", err)
	}
	defer client.Close()

	req := &automlpb.ListDatasetsRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, location),
	}

	it := client.ListDatasets(ctx, req)

	// Iterate over all results
	for {
		dataset, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("ListGlossaries.Next: %v", err)
		}

		fmt.Fprintf(w, "Dataset name: %v\n", dataset.GetName())
		fmt.Fprintf(w, "Dataset display name: %v\n", dataset.GetDisplayName())
		fmt.Fprintf(w, "Dataset create time:\n")
		fmt.Fprintf(w, "\tseconds: %v\n", dataset.GetCreateTime().GetSeconds())
		fmt.Fprintf(w, "\tnanos: %v\n", dataset.GetCreateTime().GetNanos())

		// Language entity extraction
		if metadata := dataset.GetTextExtractionDatasetMetadata(); metadata != nil {
			fmt.Fprintf(w, "Text extraction dataset metadata: %v\n", metadata)
		}

		// Language sentiment analysis
		if metadata := dataset.GetTextSentimentDatasetMetadata(); metadata != nil {
			fmt.Fprintf(w, "Text sentiment dataset metadata: %v\n", metadata)
		}

		// Language text classification
		if metadata := dataset.GetTextClassificationDatasetMetadata(); metadata != nil {
			fmt.Fprintf(w, "Text classification dataset metadata: %v\n", metadata)
		}

		// Translate
		if metadata := dataset.GetTranslationDatasetMetadata(); metadata != nil {
			fmt.Fprintf(w, "Translation dataset metadata:\n")
			fmt.Fprintf(w, "\tsource_language_code: %v\n", metadata.GetSourceLanguageCode())
			fmt.Fprintf(w, "\ttarget_language_code: %v\n", metadata.GetTargetLanguageCode())
		}

		// Vision classification
		if metadata := dataset.GetImageClassificationDatasetMetadata(); metadata != nil {
			fmt.Fprintf(w, "Image classification dataset metadata: %v\n", metadata)
		}

		// Vision object detection
		if metadata := dataset.GetImageObjectDetectionDatasetMetadata(); metadata != nil {
			fmt.Fprintf(w, "Image object detection dataset metadata: %v\n", metadata)
		}

	}

	return nil
}

// [END automl_list_datasets]
