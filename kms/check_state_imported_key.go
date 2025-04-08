// Copyright 2023 Google LLC
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

// [START kms_check_state_imported_key]
import (
	"context"
	"fmt"
	"io"

	kms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"
)

// checkStateImportedKey checks the state of a CryptoKeyVersion in KMS.
func checkStateImportedKey(w io.Writer, name string) error {
	// name := "projects/PROJECT_ID/locations/global/keyRings/my-key-ring/cryptoKeys/my-imported-key/cryptoKeyVersions/1"

	// Create the client.
	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create kms client: %w", err)
	}
	defer client.Close()

	// Call the API.
	result, err := client.GetCryptoKeyVersion(ctx, &kmspb.GetCryptoKeyVersionRequest{
		Name: name,
	})
	if err != nil {
		return fmt.Errorf("failed to get crypto key version: %w", err)
	}
	fmt.Fprintf(w, "Current state of crypto key version %q: %s\n", result.Name, result.State)
	return nil
}

// [END kms_check_state_imported_key]
