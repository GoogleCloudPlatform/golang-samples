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

// [START kms_iam_add_member]
import (
	"context"
	"fmt"
	"io"

	kms "cloud.google.com/go/kms/apiv1"
)

// iamAddMember adds a new IAM member to the Cloud KMS key
func iamAddMember(w io.Writer, name, member string) error {
	// NOTE: The resource name can be either a key or a key ring. If IAM
	// permissions are granted on the key ring, the permissions apply to all keys
	// in the key ring.
	//
	// name := "projects/my-project/locations/us-east1/keyRings/my-key-ring/cryptoKeys/my-key"
	// member := "user:foo@example.com"

	// Create the client.
	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create kms client: %v", err)
	}

	// Get the current IAM policy.
	handle := client.ResourceIAM(name)
	policy, err := handle.Policy(ctx)
	if err != nil {
		return fmt.Errorf("failed to get IAM policy: %v", err)
	}

	// Grant the member permissions. This example grants permission to use the key
	// to encrypt data.
	policy.Add(member, "roles/cloudkms.cryptoKeyEncrypterDecrypter")
	if err := handle.SetPolicy(ctx, policy); err != nil {
		return fmt.Errorf("failed to save policy: %v", err)
	}

	fmt.Fprintf(w, "Updated IAM policy for %s\n", name)
	return nil
}

// [END kms_iam_add_member]
