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

// [START healthcare_export_dicom_instance_gcs]
import (
	"context"
	"fmt"
	"io"

	healthcare "google.golang.org/api/healthcare/v1beta1"
)

// exportDICOMInstance exports DICOM objects to GCS.
func exportDICOMInstance(w io.Writer, projectID, location, datasetID, dicomStoreID, destination string) error {
	ctx := context.Background()

	healthcareService, err := healthcare.NewService(ctx)
	if err != nil {
		return fmt.Errorf("healthcare.NewService: %v", err)
	}

	storesService := healthcareService.Projects.Locations.Datasets.DicomStores

	req := &healthcare.ExportDicomDataRequest{
		GcsDestination: &healthcare.GoogleCloudHealthcareV1beta1DicomGcsDestination{
			UriPrefix: destination, // "gs://my-bucket/path/to/prefix/"
		},
	}
	name := fmt.Sprintf("projects/%s/locations/%s/datasets/%s/dicomStores/%s", projectID, location, datasetID, dicomStoreID)

	lro, err := storesService.Export(name, req).Do()
	if err != nil {
		return fmt.Errorf("Export: %v", err)
	}

	fmt.Fprintf(w, "Export to DICOM store started. Operation: %q\n", lro.Name)
	return nil
}

// [END healthcare_export_dicom_instance_gcs]
