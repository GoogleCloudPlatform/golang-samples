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
	"hash/crc32"
	"io"

	kms "cloud.google.com/go/kms/apiv1"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
	"google.golang.org/protobuf/types/known/wrapperspb"
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

	// Optional, but recommended: Compute ciphertext's CRC32C.
	crc32c := func(data []byte) uint32 {
		t := crc32.MakeTable(crc32.Castagnoli)
		return crc32.Checksum(data, t)
	}
	ciphertextCRC32C := crc32c(ciphertext)

	// Build the request.
	req := &kmspb.DecryptRequest{
		Name:             name,
		Ciphertext:       ciphertext,
		CiphertextCrc32C: wrapperspb.Int64(int64(ciphertextCRC32C)),
	}

	// Call the API.
	result, err := client.Decrypt(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to decrypt ciphertext: %v", err)
	}

	// Optional, but recommended: perform integrity verification on result.
	// For more details on ensuring E2E in-transit integrity to and from Cloud KMS visit:
	// https://cloud.google.com/kms/docs/data-integrity-guidelines
	if int64(crc32c(result.Plaintext)) != result.PlaintextCrc32C.Value {
		return fmt.Errorf("Decrypt: response corrupted in-transit")
	}

	fmt.Fprintf(w, "Decrypted plaintext: %s", result.Plaintext)
	return nil
}

// [END kms_decrypt_symmetric]
