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

// [START healthcare_fhir_execute_bundle]
import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	healthcare "google.golang.org/api/healthcare/v1"
)

// fhirExecuteBundle executes an FHIR bundle.
func fhirExecuteBundle(w io.Writer, projectID, location, datasetID, fhirStoreID string) error {
	ctx := context.Background()

	healthcareService, err := healthcare.NewService(ctx)
	if err != nil {
		return fmt.Errorf("healthcare.NewService: %w", err)
	}

	fhirService := healthcareService.Projects.Locations.Datasets.FhirStores.Fhir

	payload := map[string]interface{}{
		"resourceType": "Bundle",
		"type":         "transaction",
		"entry": []map[string]interface{}{
			{
				"resource": map[string]interface{}{
					"resourceType": "Patient",
					"active":       true,
				},
				"request": map[string]interface{}{
					"method": "POST",
					"url":    "Patient",
				},
			},
		},
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("json.Encode: %w", err)
	}

	parent := fmt.Sprintf("projects/%s/locations/%s/datasets/%s/fhirStores/%s", projectID, location, datasetID, fhirStoreID)

	call := fhirService.ExecuteBundle(parent, bytes.NewReader(jsonPayload))
	call.Header().Set("Content-Type", "application/fhir+json;charset=utf-8")
	resp, err := call.Do()
	if err != nil {
		return fmt.Errorf("executeBundle: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read response: %w", err)
	}

	if resp.StatusCode > 299 {
		return fmt.Errorf("create: status %d %s: %s", resp.StatusCode, resp.Status, respBytes)
	}
	fmt.Fprintf(w, "%s", respBytes)

	return nil
}

// [END healthcare_fhir_execute_bundle]
