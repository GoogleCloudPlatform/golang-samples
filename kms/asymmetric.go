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

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"math/big"

	cloudkms "cloud.google.com/go/kms/apiv1"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

// [START kms_create_asymmetric_key]

// createAsymmetricKey creates a new RSA encrypt/decrypt key pair on KMS.
// example keyRingName: "projects/PROJECT_ID/locations/global/keyRings/RING_ID"
func createAsymmetricKey(keyRingName, keyId string) error {
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return err
	}

	// Build the request.
	req := &kmspb.CreateCryptoKeyRequest{
		Parent:      keyRingName,
		CryptoKeyId: keyId,
		CryptoKey: &kmspb.CryptoKey{
			Purpose: kmspb.CryptoKey_ASYMMETRIC_DECRYPT,
			VersionTemplate: &kmspb.CryptoKeyVersionTemplate{
				Algorithm: kmspb.CryptoKeyVersion_RSA_DECRYPT_OAEP_2048_SHA256,
			},
		},
	}
	// Call the API.
	result, err := client.CreateCryptoKey(ctx, req)
	if err != nil {
		return err
	}
	log.Printf("Created crypto key. %s", result)
	return nil
}

// [END kms_create_asymmetric_key]

// [START kms_get_asymmetric_public]

// getAsymmetricPublicKey retrieves the public key from a saved asymmetric key pair on KMS.
// example keyName: "projects/PROJECT_ID/locations/global/keyRings/RING_ID/cryptoKeys/KEY_ID/cryptoKeyVersions/1"
func getAsymmetricPublicKey(keyName string) (interface{}, error) {
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, err
	}

	// Build the request.
	req := &kmspb.GetPublicKeyRequest{
		Name: keyName,
	}
	// Call the API.
	response, err := client.GetPublicKey(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch public key: %+v", err)
	}
	// Parse the key.
	keyBytes := []byte(response.Pem)
	block, _ := pem.Decode(keyBytes)
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %+v", err)
	}
	return publicKey, nil
}

// [END kms_get_asymmetric_public]

// [START kms_decrypt_rsa]

// decryptRSA will attempt to decrypt a given ciphertext with an 'RSA_DECRYPT_OAEP_2048_SHA256' private key.stored on Cloud KMS
// example keyName: "projects/PROJECT_ID/locations/global/keyRings/RING_ID/cryptoKeys/KEY_ID/cryptoKeyVersions/1"
func decryptRSA(keyName string, ciphertext []byte) ([]byte, error) {
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, err
	}

	// Build the request.
	req := &kmspb.AsymmetricDecryptRequest{
		Name:       keyName,
		Ciphertext: ciphertext,
	}
	// Call the API.
	response, err := client.AsymmetricDecrypt(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("decryption request failed: %+v", err)
	}
	return response.Plaintext, nil
}

// [END kms_decrypt_rsa]

// [START kms_encrypt_rsa]

// encryptRSA will encrypt data locally using an 'RSA_DECRYPT_OAEP_2048_SHA256' public key retrieved from Cloud KMS
// example keyName: "projects/PROJECT_ID/locations/global/keyRings/RING_ID/cryptoKeys/KEY_ID/cryptoKeyVersions/1"
func encryptRSA(keyName string, plaintext []byte) ([]byte, error) {
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, err
	}

	// Retrieve the public key from KMS.
	response, err := client.GetPublicKey(ctx, &kmspb.GetPublicKeyRequest{Name: keyName})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch public key: %+v", err)
	}
	// Parse the key.
	block, _ := pem.Decode([]byte(response.Pem))
	abstractKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %+v", err)
	}
	rsaKey, ok := abstractKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("key '%s' is not RSA", keyName)
	}
	// Encrypt data using the RSA public key.
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaKey, plaintext, nil)
	if err != nil {
		return nil, fmt.Errorf("encryption failed: %+v", err)
	}
	return ciphertext, nil
}

// [END kms_encrypt_rsa]

// [START kms_sign_asymmetric]

// signAsymmetric will sign a plaintext message using a saved asymmetric private key.
// example keyName: "projects/PROJECT_ID/locations/global/keyRings/RING_ID/cryptoKeys/KEY_ID/cryptoKeyVersions/1"
func signAsymmetric(keyName string, message []byte) ([]byte, error) {
	// Note: some key algorithms will require a different hash function.
	// For example, EC_SIGN_P384_SHA384 requires SHA-384.
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, err
	}
	// Find the digest of the message.
	digest := sha256.New()
	digest.Write(message)
	// Build the signing request.
	req := &kmspb.AsymmetricSignRequest{
		Name: keyName,
		Digest: &kmspb.Digest{
			Digest: &kmspb.Digest_Sha256{
				Sha256: digest.Sum(nil),
			},
		},
	}
	// Call the API.
	response, err := client.AsymmetricSign(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("asymmetric sign request failed: %+v", err)
	}
	return response.Signature, nil
}

// [END kms_sign_asymmetric]

// [START kms_verify_signature_rsa]

// verifySignatureRSA will verify that an 'RSA_SIGN_PSS_2048_SHA256' signature is valid for a given message.
// example keyName: "projects/PROJECT_ID/locations/global/keyRings/RING_ID/cryptoKeys/KEY_ID/cryptoKeyVersions/1"
func verifySignatureRSA(keyName string, signature, message []byte) error {
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return err
	}

	// Retrieve the public key from KMS.
	response, err := client.GetPublicKey(ctx, &kmspb.GetPublicKeyRequest{Name: keyName})
	if err != nil {
		return fmt.Errorf("failed to fetch public key: %+v", err)
	}
	// Parse the key.
	block, _ := pem.Decode([]byte(response.Pem))
	abstractKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse public key: %+v", err)
	}
	rsaKey, ok := abstractKey.(*rsa.PublicKey)
	if !ok {
		return fmt.Errorf("key '%s' is not RSA", keyName)
	}
	// Verify RSA signature.
	hash := sha256.New()
	hash.Write(message)
	digest := hash.Sum(nil)
	pssOptions := rsa.PSSOptions{SaltLength: len(digest), Hash: crypto.SHA256}
	err = rsa.VerifyPSS(rsaKey, crypto.SHA256, digest, signature, &pssOptions)
	if err != nil {
		return fmt.Errorf("signature verification failed: %+v", err)
	}
	return nil
}

// [END kms_verify_signature_rsa]

// [START kms_verify_signature_ec]

// verifySignatureEC will verify that an 'EC_SIGN_P256_SHA256' signature is valid for a given message.
// example keyName: "projects/PROJECT_ID/locations/global/keyRings/RING_ID/cryptoKeys/KEY_ID/cryptoKeyVersions/1"
func verifySignatureEC(keyName string, signature, message []byte) error {
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return err
	}

	// Retrieve the public key from KMS.
	response, err := client.GetPublicKey(ctx, &kmspb.GetPublicKeyRequest{Name: keyName})
	if err != nil {
		return fmt.Errorf("failed to fetch public key: %+v", err)
	}
	// Parse the key.
	block, _ := pem.Decode([]byte(response.Pem))
	abstractKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse public key: %+v", err)
	}
	ecKey, ok := abstractKey.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("key '%s' is not EC", keyName)
	}
	// Verify Elliptic Curve signature.
	var parsedSig struct{ R, S *big.Int }
	_, err = asn1.Unmarshal(signature, &parsedSig)
	if err != nil {
		return fmt.Errorf("failed to parse signature bytes: %+v", err)
	}
	hash := sha256.New()
	hash.Write(message)
	digest := hash.Sum(nil)
	if !ecdsa.Verify(ecKey, digest, parsedSig.R, parsedSig.S) {
		return errors.New("signature verification failed")
	}
	return nil
}

// [END kms_verify_signature_ec]
