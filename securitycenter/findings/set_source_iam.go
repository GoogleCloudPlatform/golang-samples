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

// Pacakge assets contains example snippets for working with findings
// and there parent resource "sources".
package findings

// [START set_iam_policy_source]
import (
	"context"
	"fmt"

	securitycenter "cloud.google.com/go/securitycenter/apiv1"
	iam "google.golang.org/genproto/googleapis/iam/v1"
)

// setIamPolicySource grants user roles/securitycenter.findingsEditor permision
// for a source.  sourceName is the full resource name of the source to be
// updated.  user is an email address. Returns the updated policy for source.
func setIamPolicySource(sourceName string, user string) (*iam.Policy, error) {
	// sourceName := "organizations/111122222444/sources/1234"
	// user := "someuser@some_domain.com
	// Instantiate a context and a security service client to make API calls.
	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("Error instantiating client %v\n", err)
	}
	defer client.Close() // Closing the client safely cleans up background resources.

	// Retrieve the existing policy so we can update only a specific
	// field.
	existing, err := client.GetIamPolicy(ctx, &iam.GetIamPolicyRequest{
		Resource: sourceName,
	})
	if err != nil {
		return nil, fmt.Errorf("Couldn't get IAM policy for %s: %v", sourceName, err)
	}

	// New IAM Binding for the user.
	newBinding := &iam.Binding{
		Role:    "roles/securitycenter.findingsEditor",
		Members: []string{fmt.Sprintf("user:%s", user)},
	}

	req := &iam.SetIamPolicyRequest{
		Resource: sourceName,
		Policy: &iam.Policy{
			// Enables partial update of existing policy
			Etag:     existing.Etag,
			Bindings: []*iam.Binding{newBinding},
		},
	}
	policy, err := client.SetIamPolicy(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("Error updating source %s with policy %v: %v",
			sourceName, req.Policy, err)
	}

	return policy, nil
}

// [END set_iam_policy_source]
