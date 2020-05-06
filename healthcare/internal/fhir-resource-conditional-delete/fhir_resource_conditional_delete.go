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

// Package conditionaldelete shows how to conditionallly delete.
// It's in a separate package so each sample can define queryParamOpt.
package conditionaldelete

// [START healthcare_conditional_delete_resource]
import (
	"context"
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

// ConditionalDeleteFHIRResource conditionally deletes an FHIR resource.
func ConditionalDeleteFHIRResource(w io.Writer, projectID, location, datasetID, fhirStoreID, resourceType string) error {
	// projectID := "my-project"
	// location := "us-central1"
	// datasetID := "my-dataset"
	// fhirStoreID := "my-fhir-store"
	// resourceType := "Patient"

	ctx := context.Background()

	healthcareService, err := healthcare.NewService(ctx)
	if err != nil {
		return fmt.Errorf("healthcare.NewService: %v", err)
	}

	fhirService := healthcareService.Projects.Locations.Datasets.FhirStores.Fhir

	parent := fmt.Sprintf("projects/%s/locations/%s/datasets/%s/fhirStores/%s", projectID, location, datasetID, fhirStoreID)

	call := fhirService.ConditionalDelete(parent, resourceType)

	// Refine your search by appending tags to the request in the form of query
	// parameters. This searches for resources updated in the last 48 hours.
	twoDaysAgo := time.Now().Add(-48 * time.Hour).Format("2006-01-02")
	lastUpdated := queryParamOpt{key: "_lastUpdated", value: "gt" + twoDaysAgo}

	if _, err := call.Do(lastUpdated); err != nil {
		return fmt.Errorf("ConditionalDelete: %v", err)
	}

	fmt.Fprintf(w, "Deleted %q", parent)

	return nil
}

// [END healthcare_conditional_delete_resource]
