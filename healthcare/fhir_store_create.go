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

// [START healthcare_create_fhir_store]
import (
	"context"
	"fmt"
	"io"

	healthcare "google.golang.org/api/healthcare/v1beta1"
)

// createFHIRStore creates an FHIR store.
func createFHIRStore(w io.Writer, projectID, location, datasetID, fhirStoreID string) error {
	ctx := context.Background()

	healthcareService, err := healthcare.NewService(ctx)
	if err != nil {
		return fmt.Errorf("healthcare.NewService: %v", err)
	}

	storesService := healthcareService.Projects.Locations.Datasets.FhirStores

	store := &healthcare.FhirStore{}
	parent := fmt.Sprintf("projects/%s/locations/%s/datasets/%s", projectID, location, datasetID)

	resp, err := storesService.Create(parent, store).FhirStoreId(fhirStoreID).Do()
	if err != nil {
		return fmt.Errorf("Create: %v", err)
	}

	fmt.Fprintf(w, "Created FHIR store: %q\n", resp.Name)
	return nil
}

// [END healthcare_create_fhir_store]
