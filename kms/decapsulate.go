// Copyright 2025 Google LLC
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

// [START kms_decapsulate]
import (
	"context"
	"fmt"
	"hash/crc32"
	"io"

	kms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// decapsulate decapsulates the given ciphertext using a saved private key of purpose
// KEY_ENCAPSULATION stored in KMS.
func decapsulate(w io.Writer, keyVersionName string, ciphertext []byte) error {
	// keyVersionName := "projects/my-project/locations/us-east1/keyRings/my-key-ring/cryptoKeys/my-key/cryptoKeyVersions/1"
	// ciphertext := []byte("...")

	// Create the client.
	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create kms client: %w", err)
	}
	defer client.Close()

	// crc32c calculates the CRC32C checksum of the given data.
	crc32c := func(data []byte) uint32 {
		t := crc32.MakeTable(crc32.Castagnoli)
		return crc32.Checksum(data, t)
	}

	// Optional but recommended: Compute ciphertext's CRC32C.
	ciphertextCRC32C := crc32c(ciphertext)

	// Build the request.
	req := &kmspb.DecapsulateRequest{
		Name:             keyVersionName,
		Ciphertext:       ciphertext,
		CiphertextCrc32C: wrapperspb.Int64(int64(ciphertextCRC32C)),
	}

	// Call the API.
	result, err := client.Decapsulate(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to decapsulate: %w", err)
	}

	// Optional, but recommended: perform integrity verification on the response.
	// For more details on ensuring E2E in-transit integrity to and from Cloud KMS visit:
	// https://cloud.google.com/kms/docs/data-integrity-guidelines
	if !result.GetVerifiedCiphertextCrc32C() {
		return fmt.Errorf("Decapsulate: request corrupted in-transit")
	}
	if result.GetName() != req.GetName() {
		return fmt.Errorf("Decapsulate: request corrupted in-transit")
	}
	if int64(crc32c(result.GetSharedSecret())) != result.GetSharedSecretCrc32C() {
		return fmt.Errorf("Decapsulate: response corrupted in-transit")
	}

	fmt.Fprintf(w, "Decapsulated plaintext: %x", result.GetSharedSecret())
	return nil
}

// [END kms_decapsulate]
