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

// [START kms_create_key_symmetric_encrypt_decrypt]
import (
	"context"
	"fmt"
	"io"

	kms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"
)

// createKeySymmetricEncryptDecrypt creates a new symmetric encrypt/decrypt key
// on Cloud KMS.
func createKeySymmetricEncryptDecrypt(w io.Writer, parent, id string) error {
	// parent := "projects/my-project/locations/us-east1/keyRings/my-key-ring"
	// id := "my-symmetric-encryption-key"

	// Create the client.
	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create kms client: %w", err)
	}
	defer client.Close()

	// Build the request.
	req := &kmspb.CreateCryptoKeyRequest{
		Parent:      parent,
		CryptoKeyId: id,
		CryptoKey: &kmspb.CryptoKey{
			Purpose: kmspb.CryptoKey_ENCRYPT_DECRYPT,
			VersionTemplate: &kmspb.CryptoKeyVersionTemplate{
				Algorithm: kmspb.CryptoKeyVersion_GOOGLE_SYMMETRIC_ENCRYPTION,
			},
		},
	}

	// Call the API.
	result, err := client.CreateCryptoKey(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create key: %w", err)
	}
	fmt.Fprintf(w, "Created key: %s\n", result.Name)
	return nil
}

// [END kms_create_key_symmetric_encrypt_decrypt]
