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

package snippets

// [START healthcare_dicom_keeplist_deidentify_dataset]
import (
	"context"
	"fmt"
	"io"
	"time"

	healthcare "google.golang.org/api/healthcare/v1beta1"
)

// deidentifyDataset creates a new dataset containing de-identified data from the source dataset.
func deidentifyDataset(w io.Writer, projectID, location, sourceDatasetID, destinationDatasetID string) error {
	ctx := context.Background()

	healthcareService, err := healthcare.NewService(ctx)
	if err != nil {
		return fmt.Errorf("healthcare.NewService: %v", err)
	}

	datasetsService := healthcareService.Projects.Locations.Datasets

	parent := fmt.Sprintf("projects/%s/locations/%s", projectID, location)

	req := &healthcare.DeidentifyDatasetRequest{
		DestinationDataset: fmt.Sprintf("%s/datasets/%s", parent, destinationDatasetID),
		Config: &healthcare.DeidentifyConfig{
			Dicom: &healthcare.DicomConfig{
				KeepList: &healthcare.TagFilterList{
					Tags: []string{
						"PatientID",
					},
				},
			},
		},
	}

	sourceName := fmt.Sprintf("%s/datasets/%s", parent, sourceDatasetID)
	resp, err := datasetsService.Deidentify(sourceName, req).Do()
	if err != nil {
		return fmt.Errorf("Deidentify: %v", err)
	}

	// Wait for the deidentification operation to finish.
	operationService := healthcareService.Projects.Locations.Datasets.Operations
	for {
		op, err := operationService.Get(resp.Name).Do()
		if err != nil {
			return fmt.Errorf("operationService.Get: %v", err)
		}
		if !op.Done {
			time.Sleep(1 * time.Second)
			continue
		}
		if op.Error != nil {
			return fmt.Errorf("deidentify operation error: %v", *op.Error)
		}
		fmt.Fprintf(w, "Created de-identified dataset %s from %s\n", resp.Name, sourceName)
		return nil
	}
}

// [END healthcare_dicom_keeplist_deidentify_dataset]
