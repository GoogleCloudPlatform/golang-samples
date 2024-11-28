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

// [START healthcare_dicomweb_store_instance]
import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	healthcare "google.golang.org/api/healthcare/v1"
)

// dicomWebStoreInstance stores the given dicomFile with the dicomWebPath.
func dicomWebStoreInstance(w io.Writer, projectID, location, datasetID, dicomStoreID, dicomWebPath, dicomFile string) error {
	ctx := context.Background()

	dicomData, err := os.ReadFile(dicomFile)
	if err != nil {
		return fmt.Errorf("os.ReadFile: %w", err)
	}

	healthcareService, err := healthcare.NewService(ctx)
	if err != nil {
		return fmt.Errorf("healthcare.NewService: %w", err)
	}

	storesService := healthcareService.Projects.Locations.Datasets.DicomStores

	parent := fmt.Sprintf("projects/%s/locations/%s/datasets/%s/dicomStores/%s", projectID, location, datasetID, dicomStoreID)

	call := storesService.StoreInstances(parent, dicomWebPath, bytes.NewReader(dicomData))
	call.Header().Set("Content-Type", "application/dicom")
	resp, err := call.Do()
	if err != nil {
		return fmt.Errorf("StoreInstances: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read response: %w", err)
	}

	if resp.StatusCode > 299 {
		return fmt.Errorf("StoreInstances: status %d %s: %s", resp.StatusCode, resp.Status, respBytes)
	}
	fmt.Fprintf(w, "%s", respBytes)
	return nil
}

// [END healthcare_dicomweb_store_instance]
