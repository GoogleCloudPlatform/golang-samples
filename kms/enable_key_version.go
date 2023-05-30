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

// [START kms_enable_key_version]
import (
	"context"
	"fmt"
	"io"

	kms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"
	fieldmask "google.golang.org/genproto/protobuf/field_mask"
)

// enableKeyVersion disables the specified key version on Cloud KMS.
func enableKeyVersion(w io.Writer, name string) error {
	// name := "projects/my-project/locations/us-east1/keyRings/my-key-ring/cryptoKeys/my-key/cryptoKeyVersions/123"

	// Create the client.
	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create kms client: %w", err)
	}
	defer client.Close()

	// Build the request.
	req := &kmspb.UpdateCryptoKeyVersionRequest{
		CryptoKeyVersion: &kmspb.CryptoKeyVersion{
			Name:  name,
			State: kmspb.CryptoKeyVersion_ENABLED,
		},
		UpdateMask: &fieldmask.FieldMask{
			Paths: []string{"state"},
		},
	}

	// Call the API.
	result, err := client.UpdateCryptoKeyVersion(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to update key version: %w", err)
	}
	fmt.Fprintf(w, "Enabled key version: %s\n", result)
	return nil
}

// [END kms_enable_key_version]
