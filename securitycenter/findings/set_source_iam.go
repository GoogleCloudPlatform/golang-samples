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

package findings

// [START securitycenter_set_source_iam]
import (
	"context"
	"fmt"
	"io"

	iam "cloud.google.com/go/iam/apiv1/iampb"
	securitycenter "cloud.google.com/go/securitycenter/apiv1"
)

// setSourceIamPolicy grants user roles/securitycenter.findingsEditor permision
// for a source. sourceName is the full resource name of the source to be
// updated. user is an email address that IAM can grant permissions to.
func setSourceIamPolicy(w io.Writer, sourceName string, user string) error {
	// sourceName := "organizations/111122222444/sources/1234"
	// user := "someuser@some_domain.com
	// Instantiate a context and a security service client to make API calls.
	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("securitycenter.NewClient: %w", err)
	}
	defer client.Close() // Closing the client safely cleans up background resources.

	// Retrieve the existing policy so we can update only a specific
	// field.
	existing, err := client.GetIamPolicy(ctx, &iam.GetIamPolicyRequest{
		Resource: sourceName,
	})
	if err != nil {
		return fmt.Errorf("GetIamPolicy(%s): %w", sourceName, err)
	}

	req := &iam.SetIamPolicyRequest{
		Resource: sourceName,
		Policy: &iam.Policy{
			// Enables partial update of existing policy
			Etag: existing.Etag,
			Bindings: []*iam.Binding{{
				Role: "roles/securitycenter.findingsEditor",
				// New IAM Binding for the user.
				Members: []string{fmt.Sprintf("user:%s", user)},
			},
			},
		},
	}
	policy, err := client.SetIamPolicy(ctx, req)
	if err != nil {
		return fmt.Errorf("SetIamPolicy(%s, %v): %w", sourceName, req.Policy, err)
	}

	fmt.Fprint(w, "Bindings:\n")
	for _, binding := range policy.Bindings {
		for _, member := range binding.Members {
			fmt.Fprintf(w, "Principal: %s Role: %s\n", member, binding.Role)
		}
	}
	return nil
}

// [END securitycenter_set_source_iam]
