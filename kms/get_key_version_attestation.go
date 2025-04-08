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

// [START kms_get_key_version_attestation]
import (
	"context"
	"fmt"
	"io"

	kms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"
)

// getKeyVersionAttestation gets the attestation on a key version, if one
// exists.
func getKeyVersionAttestation(w io.Writer, name string) error {
	// name := "projects/my-project/locations/us-east1/keyRings/my-key-ring/cryptoKeys/my-key/cryptoKeyVersions/123"

	// Create the client.
	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create kms client: %w", err)
	}

	// Build the request.
	req := &kmspb.GetCryptoKeyVersionRequest{
		Name: name,
	}

	// Call the API.
	result, err := client.GetCryptoKeyVersion(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to get key: %w", err)
	}

	// Only HSM keys have an attestation. For other key types, the attestion will
	// be nil.
	attestation := result.Attestation
	if attestation == nil {
		return fmt.Errorf("no attestation for %s", name)
	}

	// Print the attestation, hex-encoded.
	fmt.Fprintf(w, "%s: %x", attestation.Format, attestation.Content)
	return nil
}

// [END kms_get_key_version_attestation]
