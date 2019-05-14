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

// [START healthcare_hl7v2_store_get_iam_policy]
import (
	"context"
	"fmt"
	"io"

	healthcare "google.golang.org/api/healthcare/v1beta1"
)

// hl7V2IAMPolicy gets the IAM policy.
func hl7V2IAMPolicy(w io.Writer, projectID, location, datasetID, hl7V2StoreID string) error {
	ctx := context.Background()

	healthcareService, err := healthcare.NewService(ctx)
	if err != nil {
		return fmt.Errorf("healthcare.NewService: %v", err)
	}

	storesService := healthcareService.Projects.Locations.Datasets.Hl7V2Stores

	name := fmt.Sprintf("projects/%s/locations/%s/datasets/%s/hl7v2Stores/%s", projectID, location, datasetID, hl7V2StoreID)

	policy, err := storesService.GetIamPolicy(name).Do()
	if err != nil {
		return fmt.Errorf("GetIamPolicy: %v", err)
	}

	fmt.Fprintf(w, "IAM policy etag: %q\n", policy.Etag)
	return nil
}

// [END healthcare_hl7v2_store_get_iam_policy]
