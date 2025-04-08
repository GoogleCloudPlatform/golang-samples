// Copyright 2023 Google LLC
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

// Prompt: You are a Go programmer with experience in Google  Cloud. Write a Go version of the ImportDataImageClassification.java example.

package snippets

// [START aiplatform_import_data_image_classification]
import (
	"context"
	"fmt"
	"io"
	"os"

	aiplatform "cloud.google.com/go/aiplatform/apiv1"
	aiplatformpb "cloud.google.com/go/aiplatform/apiv1/aiplatformpb"
	"google.golang.org/api/option"
)

// importDataImageClassification imports data to an existing dataset.
func importDataImageClassification(w io.Writer, projectID, location, datasetID, gcsSourceUri string) error {
	// projectID := "my-project-id"
	// location := "us-central1"
	// datasetID := "YOUR_DATASET_ID"
	// gcsSourceUri := "gs://YOUR_GCS_SOURCE_BUCKET/path_to_your_image_source/[file.csv/file.jsonl]"

	apiEndpoint := fmt.Sprintf("%s-aiplatform.googleapis.com:443", location)
	clientOption := option.WithEndpoint(apiEndpoint)

	ctx := context.Background()
	aiplatformClient, err := aiplatform.NewDatasetClient(ctx, clientOption)
	if err != nil {
		return fmt.Errorf("aiplatform.NewDatasetClient: %w", err)
	}
	defer aiplatformClient.Close()

	importConfigs := []*aiplatformpb.ImportDataConfig{
		{
			Source: &aiplatformpb.ImportDataConfig_GcsSource{
				GcsSource: &aiplatformpb.GcsSource{Uris: []string{gcsSourceUri}},
			},
			ImportSchemaUri: "gs://google-cloud-aiplatform/schema/dataset/ioformat/image_classification_io_format_1.0.0.yaml",
		},
	}
	name := fmt.Sprintf("projects/%s/locations/%s/datasets/%s", projectID, location, datasetID)
	op, err := aiplatformClient.ImportData(ctx, &aiplatformpb.ImportDataRequest{
		Name:          name,
		ImportConfigs: importConfigs,
	})
	if err != nil {
		return fmt.Errorf("ImportData: %w", err)
	}
	fmt.Fprintf(os.Stdout, "Processing operation name: %q\n", op.Name())

	// ImportDataReponse, the first return value from Wait(), is an empty struct
	_, err = op.Wait(ctx)
	if err != nil {
		return err
	}

	fmt.Fprint(w, "Import Data Image Classification Response successful\n")

	return nil
}

// [END aiplatform_import_data_image_classification]
