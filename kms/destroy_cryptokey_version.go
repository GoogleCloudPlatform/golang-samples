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

// [START kms_destroy_cryptokey_version]
import (
	"context"
	"fmt"
	"io"

	cloudkms "cloud.google.com/go/kms/apiv1"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

// destroyCryptoKeyVersion marks a specified key version for deletion.
// The key can be restored if requested within 24 hours.
func destroyCryptoKeyVersion(w io.Writer, name string) error {
	// name := "projects/PROJECT_ID/locations/global/keyRings/RING_ID/cryptoKeys/KEY_ID/cryptoKeyVersions/1"
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return fmt.Errorf("cloudkms.NewKeyManagementClient: %v", err)
	}
	// Build the request.
	req := &kmspb.DestroyCryptoKeyVersionRequest{
		Name: name,
	}
	// Call the API.
	result, err := client.DestroyCryptoKeyVersion(ctx, req)
	if err != nil {
		return fmt.Errorf("DestroyCryptoKeyVersion: %v", err)
	}
	fmt.Fprintf(w, "Destroyed crypto key version: %s", result)
	return nil
}

// [END kms_destroy_cryptokey_version]
