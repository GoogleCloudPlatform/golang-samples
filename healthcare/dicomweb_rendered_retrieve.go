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

// [START healthcare_dicomweb_retrieve_rendered]
import (
	"context"
	"fmt"
	"io"
	"os"

	healthcare "google.golang.org/api/healthcare/v1"
)

// dicomWebRetrieveRendered retrieves a consumer imaging format like JPEG or PNG.
func dicomWebRetrieveRendered(w io.Writer, projectID, location, datasetID, dicomStoreID, dicomWebPath string, outputFile string) error {
	// projectID := "my-project"
	// location := "us-central1"
	// datasetID := "my-dataset"
	// dicomStoreID := "my-dicom-store"
	// dicomWebPath := "studies/1.3.6.1.4.1.11129.5.5.1113639985/series/1.3.6.1.4.1.11129.5.5.1953511724/instances/1.3.6.1.4.1.11129.5.5.9562821369/rendered"
	// outputFile := "rendered_image.png"
	ctx := context.Background()

	healthcareService, err := healthcare.NewService(ctx)
	if err != nil {
		return fmt.Errorf("healthcare.NewService: %v", err)
	}

	storesService := healthcareService.Projects.Locations.Datasets.DicomStores.Studies.Series.Instances

	parent := fmt.Sprintf("projects/%s/locations/%s/datasets/%s/dicomStores/%s", projectID, location, datasetID, dicomStoreID)

	call := storesService.RetrieveRendered(parent, dicomWebPath)
	call.Header().Set("Accept", "image/png")
	resp, err := call.Do()
	if err != nil {
		return fmt.Errorf("RetrieveRendered: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return fmt.Errorf("RetrieveRendered: status %d %s: %s", resp.StatusCode, resp.Status, resp.Body)
	}

	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("os.Create: %v", err)
	}
	defer file.Close()
	if _, err := io.Copy(file, resp.Body); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}

	fmt.Fprintf(w, "Rendered PNG image retrieved and downloaded to file: %v\n", outputFile)

	return nil
}

// [END healthcare_dicomweb_retrieve_rendered]
