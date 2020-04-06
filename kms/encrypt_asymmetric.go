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

// [START kms_encrypt_rsa]
import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	cloudkms "cloud.google.com/go/kms/apiv1"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

// encryptRSA will encrypt data locally using an 'RSA_DECRYPT_OAEP_2048_SHA256'
// public key retrieved from Cloud KMS.
func encryptRSA(name string, plaintext []byte) ([]byte, error) {
	// name: "projects/PROJECT_ID/locations/global/keyRings/RING_ID/cryptoKeys/KEY_ID/cryptoKeyVersions/1"
	// plaintext := []byte("Sample message")
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("cloudkms.NewKeyManagementClient: %v", err)
	}

	// Retrieve the public key from KMS.
	response, err := client.GetPublicKey(ctx, &kmspb.GetPublicKeyRequest{Name: name})
	if err != nil {
		return nil, fmt.Errorf("GetPublicKey: %v", err)
	}
	// Parse the key.
	block, _ := pem.Decode([]byte(response.Pem))
	abstractKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("x509.ParsePKIXPublicKey: %+v", err)
	}
	rsaKey, ok := abstractKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("key %q is not RSA", name)
	}
	// Encrypt data using the RSA public key.
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaKey, plaintext, nil)
	if err != nil {
		return nil, fmt.Errorf("rsa.EncryptOAEP: %v", err)
	}
	return ciphertext, nil
}

// [END kms_encrypt_rsa]
