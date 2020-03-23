// Copyright 2020 Google LLC
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

package secretmanager

// [START secretmanager_iam_grant_access]
import (
	"context"
	"fmt"
	"io"

	iampb "google3/google/iam/v1/iam_policy_go_proto"
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
)

// iamGrantAccess grants the given member access to the secret.
func iamGrantAccess(w io.Writer, name, member string) error {
	// name := "projects/my-project/secrets/my-secret"
	// member := "user:foo@example.com"

	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create secretmanager client: %v", err)
	}

	// Get the current IAM policy.
	getReq :=  &iampb.GetIamPolicyRequest{Resource: name}
	policy, err := client.GetIamPolicy(ctx, getReq)
	if err != nil {
		return fmt.Errorf("failed to get policy: %v", err)
	}

	// Grant the member access permissions.
	policy.Add(member, "roles/secretmanager.secretAccessor")
	setReq := &iampb.SetIamPolicyRequest{
		Resource: name,
		Policy: policy,
	}
	if err = client.SetIamPolicy(ctx, setReq); err != nil {
		return fmt.Errorf("failed to save policy: %v", err)
	}

	fmt.Fprintf(w, "Updated IAM policy for %s\n", name)
	return nil
}

// [END secretmanager_iam_grant_access]
