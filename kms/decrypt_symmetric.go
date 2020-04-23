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

// [START kms_decrypt_symmetric]
import (
	"context"
	"fmt"
	"io"

	kms "cloud.google.com/go/kms/apiv1"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

// decryptSymmetric will decrypt the input ciphertext bytes using the specified symmetric key.
func decryptSymmetric(w io.Writer, name string, ciphertext []byte) error {
	// name := "projects/my-project/locations/us-east1/keyRings/my-key-ring/cryptoKeys/my-key"
	// ciphertext := []byte("...")  // result of a symmetric encryption call

	// Create the client.
	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create kms client: %v", err)
	}

	// Build the request.
	req := &kmspb.DecryptRequest{
		Name:       name,
		Ciphertext: ciphertext,
	}

	// Call the API.
	result, err := client.Decrypt(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to decrypt ciphertext: %v", err)
	}
	fmt.Fprintf(w, "Decrypted plaintext: %s", result.Plaintext)
	return nil
}

// [END kms_decrypt_symmetric]
