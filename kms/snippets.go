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

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/iam"
	cloudkms "cloud.google.com/go/kms/apiv1"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
	fieldmask "google.golang.org/genproto/protobuf/field_mask"
)

// [START kms_create_keyring]

// createKeyRing creates a new ring to store keys on KMS.
// example parentName: "projects/PROJECT_ID/locations/global/"
func createKeyRing(w io.Writer, parentName, keyRingID string) error {
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return err
	}
	// Build the request.
	req := &kmspb.CreateKeyRingRequest{
		Parent:    parentName,
		KeyRingId: keyRingID,
	}
	// Call the API.
	result, err := client.CreateKeyRing(ctx, req)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Created key ring: %s", result)
	return nil
}

// [END kms_create_keyring]

// [START kms_create_cryptokey]

// createCryptoKey creates a new symmetric encrypt/decrypt key on KMS.
// example keyRingName: "projects/PROJECT_ID/locations/global/keyRings/RING_ID"
func createCryptoKey(w io.Writer, keyRingName, keyID string) error {
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return err
	}

	// Build the request.
	req := &kmspb.CreateCryptoKeyRequest{
		Parent:      keyRingName,
		CryptoKeyId: keyID,
		CryptoKey: &kmspb.CryptoKey{
			Purpose: kmspb.CryptoKey_ENCRYPT_DECRYPT,
			VersionTemplate: &kmspb.CryptoKeyVersionTemplate{
				Algorithm: kmspb.CryptoKeyVersion_GOOGLE_SYMMETRIC_ENCRYPTION,
			},
		},
	}
	// Call the API.
	result, err := client.CreateCryptoKey(ctx, req)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Created crypto key. %s", result)
	return nil
}

// [END kms_create_cryptokey]

// [START kms_disable_cryptokey_version]

// disableCryptoKeyVersion disables a specified key version on KMS.
// example keyVersionName: "projects/PROJECT_ID/locations/global/keyRings/RING_ID/cryptoKeys/KEY_ID/cryptoKeyVersions/1"
func disableCryptoKeyVersion(w io.Writer, keyVersionName string) error {
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return err
	}
	// Build the request.
	req := &kmspb.UpdateCryptoKeyVersionRequest{
		CryptoKeyVersion: &kmspb.CryptoKeyVersion{
			Name:  keyVersionName,
			State: kmspb.CryptoKeyVersion_DISABLED,
		},
		UpdateMask: &fieldmask.FieldMask{
			Paths: []string{"state"},
		},
	}
	// Call the API.
	result, err := client.UpdateCryptoKeyVersion(ctx, req)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Disabled crypto key version: %s", result)
	return nil
}

// [END kms_disable_cryptokey_version]

// [START kms_enable_cryptokey_version]

// enableCryptoKeyVersion enables a previously disabled key version on KMS.
// example keyVersionName: "projects/PROJECT_ID/locations/global/keyRings/RING_ID/cryptoKeys/KEY_ID/cryptoKeyVersions/1"
func enableCryptoKeyVersion(w io.Writer, keyVersionName string) error {
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return err
	}
	// Build the request.
	req := &kmspb.UpdateCryptoKeyVersionRequest{
		CryptoKeyVersion: &kmspb.CryptoKeyVersion{
			Name:  keyVersionName,
			State: kmspb.CryptoKeyVersion_ENABLED,
		},
		UpdateMask: &fieldmask.FieldMask{
			Paths: []string{"state"},
		},
	}
	// Call the API.
	result, err := client.UpdateCryptoKeyVersion(ctx, req)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Enabled crypto key version: %s", result)
	return nil
}

// [END kms_enable_cryptokey_version]

// [START kms_destroy_cryptokey_version]

// destroyCryptoKeyVersion marks a specified key version for deletion. The key can be restored if requested within 24 hours.
// example keyVersionName: "projects/PROJECT_ID/locations/global/keyRings/RING_ID/cryptoKeys/KEY_ID/cryptoKeyVersions/1"
func destroyCryptoKeyVersion(w io.Writer, keyVersionName string) error {
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return err
	}
	// Build the request.
	req := &kmspb.DestroyCryptoKeyVersionRequest{
		Name: keyVersionName,
	}
	// Call the API.
	result, err := client.DestroyCryptoKeyVersion(ctx, req)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Destroyed crypto key version: %s", result)
	return nil
}

// [END kms_destroy_cryptokey_version]

// [START kms_restore_cryptokey_version]

// restoreCryptoKeyVersion attempts to recover a key that has been marked for destruction within the last 24 hours.
// example keyVersionName: "projects/PROJECT_ID/locations/global/keyRings/RING_ID/cryptoKeys/KEY_ID/cryptoKeyVersions/1"
func restoreCryptoKeyVersion(w io.Writer, keyVersionName string) error {
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return err
	}
	// Build the request.
	req := &kmspb.RestoreCryptoKeyVersionRequest{
		Name: keyVersionName,
	}
	// Call the API.
	result, err := client.RestoreCryptoKeyVersion(ctx, req)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Restored crypto key version: %s", result)
	return nil
}

// [END kms_restore_cryptokey_version]

// [START kms_get_keyring_policy]

// getRingPolicy retrieves and prints the IAM policy associated with the key ring
// example keyRingName: "projects/PROJECT_ID/locations/global/keyRings/RING_ID"
func getRingPolicy(w io.Writer, keyRingName string) (*iam.Policy, error) {
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, err
	}
	// Get the KeyRing.
	keyRingObj, err := client.GetKeyRing(ctx, &kmspb.GetKeyRingRequest{Name: keyRingName})
	if err != nil {
		return nil, err
	}
	// Get IAM Policy.
	handle := client.KeyRingIAM(keyRingObj)
	policy, err := handle.Policy(ctx)
	if err != nil {
		return nil, err
	}
	for _, role := range policy.Roles() {
		for _, member := range policy.Members(role) {
			fmt.Fprintf(w, "Role: %s Member: %s\n", role, member)
		}
	}
	return policy, nil
}

// [END kms_get_keyring_policy]

// [START kms_get_cryptokey_policy]

// getCryptoKeyPolicy retrieves and prints the IAM policy associated with the key
// example keyName: "projects/PROJECT_ID/locations/global/keyRings/RING_ID/cryptoKeys/KEY_ID"
func getCryptoKeyPolicy(w io.Writer, keyName string) (*iam.Policy, error) {
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, err
	}
	// Get the KeyRing.
	keyObj, err := client.GetCryptoKey(ctx, &kmspb.GetCryptoKeyRequest{Name: keyName})
	if err != nil {
		return nil, err
	}
	// Get IAM Policy.
	handle := client.CryptoKeyIAM(keyObj)
	policy, err := handle.Policy(ctx)
	if err != nil {
		return nil, err
	}
	for _, role := range policy.Roles() {
		for _, member := range policy.Members(role) {
			fmt.Fprintf(w, "Role: %s Member: %s\n", role, member)
		}
	}
	return policy, nil
}

// [END kms_get_cryptokey_policy]

// [START kms_add_member_to_keyring_policy]

// addMemberRingPolicy adds a new member to a specified IAM role for the key ring
// example keyRingName: "projects/PROJECT_ID/locations/global/keyRings/RING_ID"
func addMemberRingPolicy(w io.Writer, keyRingName, member string, role iam.RoleName) error {
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return err
	}

	// Get the KeyRing.
	keyRingObj, err := client.GetKeyRing(ctx, &kmspb.GetKeyRingRequest{Name: keyRingName})
	if err != nil {
		return err
	}
	// Get IAM Policy.
	handle := client.KeyRingIAM(keyRingObj)
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
	fmt.Fprint(w, "Added member to keyring policy.")
	return nil
}

// [END kms_add_member_to_keyring_policy]

// [START kms_remove_member_from_keyring_policy]

// removeMemberRingPolicy removes a specified member from an IAM role for the key ring
// example keyRingName: "projects/PROJECT_ID/locations/global/keyRings/RING_ID"
func removeMemberRingPolicy(w io.Writer, keyRingName, member string, role iam.RoleName) error {
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return err
	}

	// Get the KeyRing.
	keyRingObj, err := client.GetKeyRing(ctx, &kmspb.GetKeyRingRequest{Name: keyRingName})
	if err != nil {
		return err
	}
	// Get IAM Policy.
	handle := client.KeyRingIAM(keyRingObj)
	policy, err := handle.Policy(ctx)
	if err != nil {
		return err
	}

	// Remove Member.
	policy.Remove(member, role)
	err = handle.SetPolicy(ctx, policy)
	if err != nil {
		return err
	}
	fmt.Fprint(w, "Removed member from keyring policy.")
	return nil
}

// [END kms_remove_member_from_keyring_policy]

// [START kms_add_member_to_cryptokey_policy]

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

// [START kms_remove_member_from_cryptokey_policy]

// removeMemberCryptoKeyPolicy removes a specified member from an IAM role for the key
// example keyName: "projects/PROJECT_ID/locations/global/keyRings/RING_ID/cryptoKeys/KEY_ID"
func removeMemberCryptoKeyPolicy(w io.Writer, keyName, member string, role iam.RoleName) error {
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
	// Remove Member.
	policy.Remove(member, role)
	err = handle.SetPolicy(ctx, policy)
	if err != nil {
		return err
	}
	fmt.Fprint(w, "Removed member from cryptokey policy.")
	return nil
}

// [END kms_remove_member_from_cryptokey_policy]

// [START kms_encrypt]

// encrypt will encrypt the input plaintext with the specified symmetric key
// example keyName: "projects/PROJECT_ID/locations/global/keyRings/RING_ID/cryptoKeys/KEY_ID"
func encryptSymmetric(keyName string, plaintext []byte) ([]byte, error) {
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, err
	}

	// Build the request.
	req := &kmspb.EncryptRequest{
		Name:      keyName,
		Plaintext: plaintext,
	}
	// Call the API.
	resp, err := client.Encrypt(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Ciphertext, nil
}

// [END kms_encrypt]

// [START kms_decrypt]

// decrypt will decrypt the input ciphertext bytes using the specified symmetric key
// example keyName: "projects/PROJECT_ID/locations/global/keyRings/RING_ID/cryptoKeys/KEY_ID"
func decryptSymmetric(keyName string, ciphertext []byte) ([]byte, error) {
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, err
	}

	// Build the request.
	req := &kmspb.DecryptRequest{
		Name:       keyName,
		Ciphertext: ciphertext,
	}
	// Call the API.
	resp, err := client.Decrypt(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Plaintext, nil
}

// [END kms_decrypt]
