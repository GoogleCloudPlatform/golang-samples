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

// [START healthcare_fhir_store_set_iam_policy]
import (
	"context"
	"fmt"
	"io"

	healthcare "google.golang.org/api/healthcare/v1beta1"
)

// setFHIRIAMPolicy sets the FHIR store's IAM policy.
func setFHIRIAMPolicy(w io.Writer, projectID, location, datasetID, fhirStoreID string) error {
	ctx := context.Background()

	healthcareService, err := healthcare.NewService(ctx)
	if err != nil {
		return fmt.Errorf("healthcare.NewService: %v", err)
	}

	fhirService := healthcareService.Projects.Locations.Datasets.FhirStores

	name := fmt.Sprintf("projects/%s/locations/%s/datasets/%s/fhirStores/%s", projectID, location, datasetID, fhirStoreID)

	policy, err := fhirService.GetIamPolicy(name).Do()
	if err != nil {
		return fmt.Errorf("GetIamPolicy: %v", err)
	}

	policy.Bindings = append(policy.Bindings, &healthcare.Binding{
		Members: []string{"user:example@example.com"},
		Role:    "roles/viewer",
	})

	req := &healthcare.SetIamPolicyRequest{
		Policy: policy,
	}

	policy, err = fhirService.SetIamPolicy(name, req).Do()
	if err != nil {
		return fmt.Errorf("SetIamPolicy: %v", err)
	}

	fmt.Fprintf(w, "IAM Policy version: %v\n", policy.Version)
	return nil
}

// [END healthcare_fhir_store_set_iam_policy]
