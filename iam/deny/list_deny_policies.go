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

// [START iam_list_deny_policy]
import (
	"context"
	"fmt"
	"io"

	iam "cloud.google.com/go/iam/apiv2"
	"cloud.google.com/go/iam/apiv2/iampb"
	"google.golang.org/api/iterator"
)

// listDenyPolicies lists all the deny policies that are attached to a resource.
// A resource can have up to 5 deny policies.
func listDenyPolicies(w io.Writer, projectID string) error {
	// projectID := "your_project_id"

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

	req := &iampb.ListPoliciesRequest{
		// Construct the full path of the resource's deny policies.
		// Its format is: "policies/ATTACHMENT_POINT/denypolicies"
		Parent: fmt.Sprintf("policies/%s/denypolicies", attachmentPoint),
	}
	it := policiesClient.ListPolicies(ctx, req)
	fmt.Fprintf(w, "Policies found in project %s:\n", projectID)

	for {
		policy, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "- %s\n", policy.GetName())
	}
	return nil
}

// [END iam_list_deny_policy]
