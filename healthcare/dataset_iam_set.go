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

// [START healthcare_dataset_set_iam_policy]
import (
	"context"
	"fmt"
	"io"

	healthcare "google.golang.org/api/healthcare/v1"
)

// setDatasetIAMPolicy sets an IAM policy for the dataset.
func setDatasetIAMPolicy(w io.Writer, projectID, location, datasetID string) error {
	ctx := context.Background()

	healthcareService, err := healthcare.NewService(ctx)
	if err != nil {
		return fmt.Errorf("healthcare.NewService: %w", err)
	}

	datasetsService := healthcareService.Projects.Locations.Datasets

	name := fmt.Sprintf("projects/%s/locations/%s/datasets/%s", projectID, location, datasetID)

	policy, err := datasetsService.GetIamPolicy(name).Do()
	if err != nil {
		return fmt.Errorf("GetIamPolicy: %w", err)
	}

	policy.Bindings = append(policy.Bindings, &healthcare.Binding{
		Members: []string{"user:example@example.com"},
		Role:    "roles/viewer",
	})

	req := &healthcare.SetIamPolicyRequest{
		Policy: policy,
	}

	policy, err = datasetsService.SetIamPolicy(name, req).Do()
	if err != nil {
		return fmt.Errorf("SetIamPolicy: %w", err)
	}

	fmt.Fprintf(w, "IAM Policy etag: %v", policy.Etag)
	return nil
}

// [END healthcare_dataset_set_iam_policy]
