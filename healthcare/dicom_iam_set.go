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

// [START healthcare_dicom_store_set_iam_policy]
import (
	"context"
	"fmt"
	"io"

	healthcare "google.golang.org/api/healthcare/v1"
)

// setDICOMIAMPolicy sets the DICOM store's IAM policy.
func setDICOMIAMPolicy(w io.Writer, projectID, location, datasetID, dicomStoreID string) error {
	ctx := context.Background()

	healthcareService, err := healthcare.NewService(ctx)
	if err != nil {
		return fmt.Errorf("healthcare.NewService: %v", err)
	}

	dicomService := healthcareService.Projects.Locations.Datasets.DicomStores

	name := fmt.Sprintf("projects/%s/locations/%s/datasets/%s/dicomStores/%s", projectID, location, datasetID, dicomStoreID)

	policy, err := dicomService.GetIamPolicy(name).Do()
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

	policy, err = dicomService.SetIamPolicy(name, req).Do()
	if err != nil {
		return fmt.Errorf("SetIamPolicy: %v", err)
	}

	fmt.Fprintf(w, "IAM Policy etag: %v\n", policy.Etag)
	return nil
}

// [END healthcare_dicom_store_set_iam_policy]
