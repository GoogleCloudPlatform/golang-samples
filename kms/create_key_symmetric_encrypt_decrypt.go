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

// [START kms_create_cryptokey]
import (
	"context"
	"fmt"
	"io"

	cloudkms "cloud.google.com/go/kms/apiv1"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

// createCryptoKey creates a new symmetric encrypt/decrypt key on KMS.
func createCryptoKey(w io.Writer, keyRingName, keyID string) error {
	// keyRingName := "projects/PROJECT_ID/locations/global/keyRings/RING_ID"
	// keyID := "key-" + strconv.Itoa(int(time.Now().Unix()))
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return fmt.Errorf("cloudkms.NewKeyManagementClient: %v", err)
	}

	// Build the request.
	req := &kmspb.CreateCryptoKeyRequest{
		Parent:      keyRingName,
		CryptoKeyId: keyID,
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
		return fmt.Errorf("CreateCryptoKey: %v", err)
	}
	fmt.Fprintf(w, "Created crypto key. %s", result)
	return nil
}

// [END kms_create_cryptokey]
