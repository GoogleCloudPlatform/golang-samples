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

package kms

// [START kms_get_keyring_policy]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/iam"
	cloudkms "cloud.google.com/go/kms/apiv1"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

// getRingPolicy retrieves and prints the IAM policy associated with the key ring.
func getRingPolicy(w io.Writer, name string) (*iam.Policy, error) {
	// name: "projects/PROJECT_ID/locations/global/keyRings/RING_ID"
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("cloudkms.NewKeyManagementClient: %v", err)
	}
	// Get the KeyRing.
	keyRingObj, err := client.GetKeyRing(ctx, &kmspb.GetKeyRingRequest{Name: name})
	if err != nil {
		return nil, fmt.Errorf("GetKeyRing: %v", err)
	}
	// Get IAM Policy.
	handle := client.KeyRingIAM(keyRingObj)
	policy, err := handle.Policy(ctx)
	if err != nil {
		return nil, fmt.Errorf("Policy: %v", err)
	}
	for _, role := range policy.Roles() {
		for _, member := range policy.Members(role) {
			fmt.Fprintf(w, "Role: %s Member: %s\n", role, member)
		}
	}
	return policy, nil
}

// [END kms_get_keyring_policy]
