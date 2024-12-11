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

// [START healthcare_get_patient_everything]
import (
	"context"
	"fmt"
	"io"

	healthcare "google.golang.org/api/healthcare/v1"
)

// fhirGetPatientEverything gets all resources associated with a particular
// patient compartment.
func fhirGetPatientEverything(w io.Writer, projectID, location, datasetID, fhirStoreID, fhirResourceID string) error {
	ctx := context.Background()

	healthcareService, err := healthcare.NewService(ctx)
	if err != nil {
		return fmt.Errorf("healthcare.NewService: %w", err)
	}
	fhirService := healthcareService.Projects.Locations.Datasets.FhirStores.Fhir
	name := fmt.Sprintf("projects/%s/locations/%s/datasets/%s/fhirStores/%s/fhir/Patient/%s", projectID, location, datasetID, fhirStoreID, fhirResourceID)

	resp, err := fhirService.PatientEverything(name).Do()
	if err != nil {
		return fmt.Errorf("PatientEverything: %w", err)
	}

	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read response: %w", err)
	}

	if resp.StatusCode > 299 {
		return fmt.Errorf("PatientEverything: status %d %s: %s", resp.StatusCode, resp.Status, respBytes)
	}
	fmt.Fprintf(w, "%s", respBytes)

	return nil
}

// [END healthcare_get_patient_everything]
