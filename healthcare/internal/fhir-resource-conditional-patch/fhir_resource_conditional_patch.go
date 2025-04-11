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

// Package conditionalpatch shows how to conditionallly patch.
// It's in a separate package so each sample can define queryParamOpt.
package conditionalpatch

// [START healthcare_conditional_patch_resource]

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	healthcare "google.golang.org/api/healthcare/v1beta1"
)

// queryParamOpt is a googleapi.Option (https://godoc.org/google.golang.org/api/googleapi#CallOption)
// that adds query parameters to an API call.
type queryParamOpt struct {
	key, value string
}

func (qp queryParamOpt) Get() (string, string) { return qp.key, qp.value }

// ConditionalPatchFHIRResource conditionally patches an FHIR resource.
func ConditionalPatchFHIRResource(w io.Writer, projectID, location, datasetID, fhirStoreID, resourceType string, active bool) error {
	// projectID := "my-project"
	// location := "us-central1"
	// datasetID := "my-dataset"
	// fhirStoreID := "my-fhir-store"
	// resourceType := "Patient"
	// active := true

	ctx := context.Background()

	healthcareService, err := healthcare.NewService(ctx)
	if err != nil {
		return fmt.Errorf("healthcare.NewService: %w", err)
	}

	fhirService := healthcareService.Projects.Locations.Datasets.FhirStores.Fhir

	parent := fmt.Sprintf("projects/%s/locations/%s/datasets/%s/fhirStores/%s", projectID, location, datasetID, fhirStoreID)

	payload := []map[string]interface{}{
		{
			"op":    "replace",
			"path":  "/active",
			"value": active,
		},
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("json.Encode: %w", err)
	}

	call := fhirService.ConditionalPatch(parent, resourceType, bytes.NewReader(jsonPayload))

	call.Header().Set("Content-Type", "application/json-patch+json")

	// Refine your search by appending tags to the request in the form of query
	// parameters. This searches for resources updated in the last 48 hours.
	twoDaysAgo := time.Now().Add(-48 * time.Hour).Format("2006-01-02")
	lastUpdated := queryParamOpt{key: "_lastUpdated", value: "gt" + twoDaysAgo}

	resp, err := call.Do(lastUpdated)
	if err != nil {
		return fmt.Errorf("ConditionalPatch: %w", err)
	}

	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read response: %w", err)
	}

	if resp.StatusCode > 299 {
		return fmt.Errorf("ConditionalPatch: status %d %s: %s", resp.StatusCode, resp.Status, respBytes)
	}
	fmt.Fprintf(w, "%s", respBytes)

	return nil
}

// [END healthcare_conditional_patch_resource]
