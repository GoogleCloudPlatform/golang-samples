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

// [START iam_update_deny_policy]
import (
	"context"
	"fmt"
	"io"

	iam "cloud.google.com/go/iam/apiv2"

	"cloud.google.com/go/iam/apiv2/iampb"
	"google.golang.org/genproto/googleapis/type/expr"
)

// updateDenyPolicy updates the deny rules and/ or its display name after policy creation.
func updateDenyPolicy(w io.Writer, projectID, policyID, etag string) error {
	// projectID := "your_project_id"
	// policyID := "your_policy_id"
	// etag := "your_etag"

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

	denyRule := &iampb.DenyRule{
		// Add one or more principals who should be denied the permissions specified in this rule.
		// For more information on allowed values,
		// see: https://cloud.google.com/iam/help/deny/principal-identifiers
		DeniedPrincipals: []string{"principalSet://goog/public:all"},
		// Optionally, set the principals who should be exempted from the
		// list of denied principals. For example, if you want to deny certain permissions
		// to a group but exempt a few principals, then add those here.
		// ExceptionPrincipals: []string{"principalSet://goog/group/project-admins@example.com"},
		//
		// Set the permissions to deny.
		// The permission value is of the format: service_fqdn/resource.action
		// For the list of supported permissions,
		// see: https://cloud.google.com/iam/help/deny/supported-permissions
		DeniedPermissions: []string{"cloudresourcemanager.googleapis.com/projects.delete"},
		// Optionally, add the permissions to be exempted from this rule.
		// Meaning, the deny rule will not be applicable to these permissions.
		// ExceptionPermissions: []string{"cloudresourcemanager.googleapis.com/projects.create"},
		//
		// Set the condition which will enforce the deny rule.
		// If this condition is true, the deny rule will be applicable.
		// Else, the rule will not be enforced.
		// The expression uses Common Expression Language syntax (CEL).
		// Here we block access based on tags.
		//
		// Here, we create a deny rule that denies the
		// cloudresourcemanager.googleapis.com/projects.delete permission
		// to everyone except project-admins@example.com for resources that are tagged prod.
		// A tag is a key-value pair that can be attached to an organization, folder, or project.
		// For more info, see: https://cloud.google.com/iam/docs/deny-access#create-deny-policy
		DenialCondition: &expr.Expr{
			Expression: "!resource.matchTag('12345678/env', 'prod')",
		},
	}

	// Set the rule description and deny rule to update.
	policyRule := &iampb.PolicyRule{
		Description: "block all principals from deleting projects, unless the principal is a member of project-admins@example.com and the project being deleted has a tag with the value prod",
		Kind: &iampb.PolicyRule_DenyRule{
			DenyRule: denyRule,
		},
	}

	// Set the policy resource path, version (etag) and the updated deny rules.
	policy := &iampb.Policy{
		// Construct the full path of the policy.
		// Its format is: "policies/ATTACHMENT_POINT/denypolicies/POLICY_ID"
		Name:  fmt.Sprintf("policies/%s/denypolicies/%s", attachmentPoint, policyID),
		Etag:  etag,
		Rules: [](*iampb.PolicyRule){policyRule},
	}

	// Create the update policy request.
	req := &iampb.UpdatePolicyRequest{
		Policy: policy,
	}
	op, err := policiesClient.UpdatePolicy(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to update policy: %w", err)
	}

	policy, err = op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("unable to wait for the operation: %w", err)
	}

	fmt.Fprintf(w, "Policy %s updated\n", policy.GetName())

	return nil
}

// [END iam_update_deny_policy]
