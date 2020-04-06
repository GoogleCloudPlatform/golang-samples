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

// [START kms_decrypt]
import (
	"context"
	"fmt"

	cloudkms "cloud.google.com/go/kms/apiv1"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

// decryptSymmetric will decrypt the input ciphertext bytes using the specified symmetric key.
func decryptSymmetric(name string, ciphertext []byte) ([]byte, error) {
	// name := "projects/PROJECT_ID/locations/global/keyRings/RING_ID/cryptoKeys/KEY_ID"
	// cipherBytes, err := encryptRSA(rsaDecryptPath, []byte("Sample message"))
	// if err != nil {
	//   return nil, fmt.Errorf("encryptRSA: %v", err)
	// }
	// ciphertext := base64.StdEncoding.EncodeToString(cipherBytes)
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("cloudkms.NewKeyManagementClient: %v", err)
	}

	// Build the request.
	req := &kmspb.DecryptRequest{
		Name:       name,
		Ciphertext: ciphertext,
	}
	// Call the API.
	resp, err := client.Decrypt(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("Decrypt: %v", err)
	}
	return resp.Plaintext, nil
}

// [END kms_decrypt]
