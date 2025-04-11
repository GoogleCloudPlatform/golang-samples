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

package kms

// [START kms_iam_get_policy]
import (
	"context"
	"fmt"
	"io"

	kms "cloud.google.com/go/kms/apiv1"
)

// iamGetPolicy retrieves and prints the Cloud IAM policy associated with the
// Cloud KMS key.
func iamGetPolicy(w io.Writer, name string) error {
	// NOTE: The resource name can be either a key or a key ring.
	//
	// name := "projects/my-project/locations/us-east1/keyRings/my-key-ring/cryptoKeys/my-key"
	// name := "projects/my-project/locations/us-east1/keyRings/my-key-ring"

	// Create the client.
	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create kms client: %w", err)
	}
	defer client.Close()

	// Get the current policy.
	policy, err := client.ResourceIAM(name).Policy(ctx)
	if err != nil {
		return fmt.Errorf("failed to get IAM policy: %w", err)
	}

	// Print the policy members.
	for _, role := range policy.Roles() {
		fmt.Fprintf(w, "%s\n", role)
		for _, member := range policy.Members(role) {
			fmt.Fprintf(w, "- %s\n", member)
		}
		fmt.Fprintf(w, "\n")
	}
	return nil
}

// [END kms_iam_get_policy]
