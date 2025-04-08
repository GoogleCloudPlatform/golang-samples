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

// [START kms_verify_asymmetric_signature_ec]
import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"

	kms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"
)

// verifyAsymmetricSignatureEC will verify that an 'EC_SIGN_P256_SHA256' signature is
// valid for a given message.
func verifyAsymmetricSignatureEC(w io.Writer, name string, message, signature []byte) error {
	// name := "projects/my-project/locations/us-east1/keyRings/my-key-ring/cryptoKeys/my-key/cryptoKeyVersions/123"
	// message := "my message"
	// signature := []byte("...")  // Response from a sign request

	// Create the client.
	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create kms client: %w", err)
	}
	defer client.Close()

	// Retrieve the public key from KMS.
	response, err := client.GetPublicKey(ctx, &kmspb.GetPublicKeyRequest{Name: name})
	if err != nil {
		return fmt.Errorf("failed to get public key: %w", err)
	}

	// Parse the public key. Note, this example assumes the public key is in the
	// ECDSA format.
	block, _ := pem.Decode([]byte(response.Pem))
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse public key: %w", err)
	}
	ecKey, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("public key is not elliptic curve")
	}

	// Verify Elliptic Curve signature.
	var parsedSig struct{ R, S *big.Int }
	if _, err = asn1.Unmarshal(signature, &parsedSig); err != nil {
		return fmt.Errorf("asn1.Unmarshal: %w", err)
	}

	digest := sha256.Sum256(message)
	if !ecdsa.Verify(ecKey, digest[:], parsedSig.R, parsedSig.S) {
		return fmt.Errorf("failed to verify signature")
	}
	fmt.Fprintf(w, "Verified signature!")
	return nil
}

// [END kms_verify_asymmetric_signature_ec]
