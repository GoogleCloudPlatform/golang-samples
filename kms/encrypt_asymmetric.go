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

// [START kms_encrypt_asymmetric]
import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"

	kms "cloud.google.com/go/kms/apiv1"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

// encryptAsymmetric encrypts data on your local machine using an
// 'RSA_DECRYPT_OAEP_2048_SHA256' public key retrieved from Cloud KMS.
func encryptAsymmetric(w io.Writer, name string, message string) error {
	// name := "projects/my-project/locations/us-east1/keyRings/my-key-ring/cryptoKeys/my-key/cryptoKeyVersions/123"
	// message := "Sample message"

	// Create the client.
	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create kms client: %v", err)
	}

	// Retrieve the public key from Cloud KMS. This is the only operation that
	// involves Cloud KMS. The remaining operations take place on your local
	// machine.
	response, err := client.GetPublicKey(ctx, &kmspb.GetPublicKeyRequest{
		Name: name,
	})
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

	// Convert the message into bytes. Cryptographic plaintexts and
	// ciphertexts are always byte arrays.
	plaintext := []byte(message)

	// Encrypt data using the RSA public key.
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaKey, plaintext, nil)
	if err != nil {
		return fmt.Errorf("rsa.EncryptOAEP: %v", err)
	}
	fmt.Fprintf(w, "Encrypted ciphertext: %s", ciphertext)
	return nil
}

// [END kms_encrypt_asymmetric]
