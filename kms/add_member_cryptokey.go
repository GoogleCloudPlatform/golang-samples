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

// Package kms contains samples for asymmetric keys feature of Cloud Key Management Service
// https://cloud.google.com/kms/
package kms

// [START kms_add_member_to_cryptokey_policy]

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/iam"
	cloudkms "cloud.google.com/go/kms/apiv1"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

// addMemberCryptoKeyPolicy adds a new member to a specified IAM role for the key
// example keyName: "projects/PROJECT_ID/locations/global/keyRings/RING_ID/cryptoKeys/KEY_ID"
func addMemberCryptoKeyPolicy(w io.Writer, keyName, member string, role iam.RoleName) error {
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return err
	}

	// Get the desired CryptoKey.
	keyObj, err := client.GetCryptoKey(ctx, &kmspb.GetCryptoKeyRequest{Name: keyName})
	if err != nil {
		return err
	}
	// Get IAM Policy.
	handle := client.CryptoKeyIAM(keyObj)
	policy, err := handle.Policy(ctx)
	if err != nil {
		return err
	}
	// Add Member.
	policy.Add(member, role)
	err = handle.SetPolicy(ctx, policy)
	if err != nil {
		return err
	}
	fmt.Fprint(w, "Added member to cryptokey policy.")
	return nil
}

// [END kms_add_member_to_cryptokey_policy]
