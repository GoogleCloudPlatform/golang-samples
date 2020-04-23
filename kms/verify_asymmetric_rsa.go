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

// [START kms_verify_asymmetric_signature_rsa]
import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"

	kms "cloud.google.com/go/kms/apiv1"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

// verifyAsymmetricSignatureRSA will verify that an 'RSA_SIGN_PSS_2048_SHA256' signature
// is valid for a given message.
func verifyAsymmetricSignatureRSA(w io.Writer, name string, message, signature []byte) error {
	// name := "projects/my-project/locations/us-east1/keyRings/my-key-ring/cryptoKeys/my-key/cryptoKeyVersions/123"
	// message := "my message"
	// signature := []byte("...")  // Response from a sign request

	// Create the client.
	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create kms client: %v", err)
	}

	// Retrieve the public key from KMS.
	response, err := client.GetPublicKey(ctx, &kmspb.GetPublicKeyRequest{Name: name})
	if err != nil {
		return fmt.Errorf("failed to get public key: %v", err)
	}

	// Parse the public key. Note, this example assumes the public key is in the
	// RSA format.
	block, _ := pem.Decode([]byte(response.Pem))
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse public key: %v", err)
	}
	rsaKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return fmt.Errorf("public key is not rsa")
	}

	// Verify the RSA signature.
	digest := sha256.Sum256(message)
	if err := rsa.VerifyPSS(rsaKey, crypto.SHA256, digest[:], signature, &rsa.PSSOptions{
		SaltLength: len(digest),
		Hash:       crypto.SHA256,
	}); err != nil {
		return fmt.Errorf("failed to verify signature: %v", err)
	}

	fmt.Fprint(w, "Verified signature!\n")
	return nil
}

// [END kms_verify_asymmetric_signature_rsa]
