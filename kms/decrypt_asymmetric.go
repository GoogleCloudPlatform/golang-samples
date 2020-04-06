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

// [START kms_decrypt_rsa]
import (
	"context"
	"fmt"

	cloudkms "cloud.google.com/go/kms/apiv1"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

// decryptRSA will attempt to decrypt a given ciphertext with an
// 'RSA_DECRYPT_OAEP_2048_SHA256' private key stored on Cloud KMS.
func decryptRSA(name string, ciphertext []byte) ([]byte, error) {
	// name := "projects/PROJECT_ID/locations/global/keyRings/RING_ID/cryptoKeys/KEY_ID/cryptoKeyVersions/1"
	// ciphertext, err := encryptRSA(rsaDecryptPath, []byte("Sample message"))
	// if err != nil {
	//   return nil, fmt.Errorf("encryptRSA: %v", err)
	// }
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("cloudkms.NewKeyManagementClient: %v", err)
	}

	// Build the request.
	req := &kmspb.AsymmetricDecryptRequest{
		Name:       name,
		Ciphertext: ciphertext,
	}
	// Call the API.
	response, err := client.AsymmetricDecrypt(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("AsymmetricDecrypt: %v", err)
	}
	return response.Plaintext, nil
}

// [END kms_decrypt_rsa]
