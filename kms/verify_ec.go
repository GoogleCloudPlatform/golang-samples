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

// [START kms_verify_signature_ec]

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"

	cloudkms "cloud.google.com/go/kms/apiv1"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

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
