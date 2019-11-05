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

// [START kms_verify_signature_rsa]
import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	cloudkms "cloud.google.com/go/kms/apiv1"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

// verifySignatureRSA will verify that an 'RSA_SIGN_PSS_2048_SHA256' signature is valid for a given message.
func verifySignatureRSA(name string, signature, message []byte) error {
	// name: "projects/PROJECT_ID/locations/global/keyRings/RING_ID/cryptoKeys/KEY_ID/cryptoKeyVersions/1"
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return fmt.Errorf("cloudkms.NewKeyManagementClient: %v", err)
	}
	// Retrieve the public key from KMS.
	response, err := client.GetPublicKey(ctx, &kmspb.GetPublicKeyRequest{Name: name})
	if err != nil {
		return fmt.Errorf("GetPublicKey: %v", err)
	}
	// Parse the key.
	block, _ := pem.Decode([]byte(response.Pem))
	abstractKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("x509.ParsePKIXPublicKey: %v", err)
	}
	rsaKey, ok := abstractKey.(*rsa.PublicKey)
	if !ok {
		return fmt.Errorf("key '%s' is not RSA", name)
	}
	// Verify RSA signature.
	hash := sha256.New()
	hash.Write(message)
	digest := hash.Sum(nil)
	pssOptions := rsa.PSSOptions{SaltLength: len(digest), Hash: crypto.SHA256}
	if err = rsa.VerifyPSS(rsaKey, crypto.SHA256, digest, signature, &pssOptions); err != nil {
		return fmt.Errorf("rsa.VerifyPSS: %v", err)
	}
	return nil
}

// [END kms_verify_signature_rsa]
