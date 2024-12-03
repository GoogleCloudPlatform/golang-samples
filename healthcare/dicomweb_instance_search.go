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

// [START healthcare_dicomweb_search_instances]
import (
	"context"
	"fmt"
	"io"

	healthcare "google.golang.org/api/healthcare/v1"
)

// dicomWebSearchInstances searches instances.
func dicomWebSearchInstances(w io.Writer, projectID, location, datasetID, dicomStoreID string) error {
	// projectID := "my-project"
	// location := "us-central1"
	// datasetID := "my-dataset"
	// dicomStoreID := "my-dicom-store"
	ctx := context.Background()

	healthcareService, err := healthcare.NewService(ctx)
	if err != nil {
		return fmt.Errorf("healthcare.NewService: %w", err)
	}

	storesService := healthcareService.Projects.Locations.Datasets.DicomStores

	parent := fmt.Sprintf("projects/%s/locations/%s/datasets/%s/dicomStores/%s", projectID, location, datasetID, dicomStoreID)

	resp, err := storesService.SearchForInstances(parent, "instances").Do()
	if err != nil {
		return fmt.Errorf("SearchForInstances: %w", err)
	}

	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("ioutil.ReadAll: %w", err)
	}

	if resp.StatusCode > 299 {
		return fmt.Errorf("SearchForInstances: status %d %s: %s", resp.StatusCode, resp.Status, respBytes)
	}

	respString := string(respBytes)
	fmt.Fprintf(w, "Found instances: %s\n", respString)
	return nil
}

// [END healthcare_dicomweb_search_instances]
