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

// [START kms_update_key_set_primary]
import (
	"context"
	"fmt"
	"io"

	kms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"
)

// updateKeySetPrimary updates the primary key version on a Cloud KMS key.
func updateKeySetPrimary(w io.Writer, name, version string) error {
	// name := "projects/my-project/locations/us-east1/keyRings/my-key-ring/cryptoKeys/my-key"
	// version := "123"

	// Create the client.
	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create kms client: %w", err)
	}
	defer client.Close()

	// Build the request.
	req := &kmspb.UpdateCryptoKeyPrimaryVersionRequest{
		Name:               name,
		CryptoKeyVersionId: version,
	}

	// Call the API.
	result, err := client.UpdateCryptoKeyPrimaryVersion(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to update key: %w", err)
	}
	fmt.Fprintf(w, "Updated key primary: %s\n", result.Name)

	return nil
}

// [END kms_update_key_set_primary]
