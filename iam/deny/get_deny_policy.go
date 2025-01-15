// Copyright 2022 Google LLC
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

// [START iam_get_deny_policy]
import (
	"context"
	"fmt"
	"io"

	iam "cloud.google.com/go/iam/apiv2"
	"cloud.google.com/go/iam/apiv2/iampb"
)

// getDenyPolicy retrieves the deny policy given the project ID and policy ID.
func getDenyPolicy(w io.Writer, projectID, policyID string) error {
	// projectID := "your_project_id"
	// policyID := "your_policy_id"

	ctx := context.Background()
	policiesClient, err := iam.NewPoliciesClient(ctx)
	if err != nil {
		return fmt.Errorf("NewPoliciesClient: %w", err)
	}
	defer policiesClient.Close()

	// Each deny policy is attached to an organization, folder, or project.
	// To work with deny policies, specify the attachment point.
	//
	// Its format can be one of the following:
	// 1. cloudresourcemanager.googleapis.com/organizations/ORG_ID
	// 2. cloudresourcemanager.googleapis.com/folders/FOLDER_ID
	// 3. cloudresourcemanager.googleapis.com/projects/PROJECT_ID
	//
	// The attachment point is identified by its URL-encoded resource name. Hence, replace
	// the "/" with "%%2F".
	attachmentPoint := fmt.Sprintf(
		"cloudresourcemanager.googleapis.com%%2Fprojects%%2F%s",
		projectID,
	)

	req := &iampb.GetPolicyRequest{
		// Construct the full path of the policy.
		// Its format is: "policies/ATTACHMENT_POINT/denypolicies/POLICY_ID"
		Name: fmt.Sprintf("policies/%s/denypolicies/%s", attachmentPoint, policyID),
	}
	policy, err := policiesClient.GetPolicy(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to get policy: %w", err)
	}

	fmt.Fprintf(w, "Policy %s retrieved\n", policy.GetName())

	return nil
}

// [END iam_get_deny_policy]
