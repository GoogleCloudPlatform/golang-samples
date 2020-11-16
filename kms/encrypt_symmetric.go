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

// [START kms_encrypt_symmetric]
import (
	"context"
	"fmt"
	"hash/crc32"
	"io"

	kms "cloud.google.com/go/kms/apiv1"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

// encryptSymmetric encrypts the input plaintext with the specified symmetric
// Cloud KMS key.
func encryptSymmetric(w io.Writer, name string, message string) error {
	// name := "projects/my-project/locations/us-east1/keyRings/my-key-ring/cryptoKeys/my-key"
	// message := "Sample message"

	// Create the client.
	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create kms client: %v", err)
	}

	// Convert the message into bytes. Cryptographic plaintexts and
	// ciphertexts are always byte arrays.
	plaintext := []byte(message)

	// Optional, but recommended: Compute plaintext's CRC32C.
	// See crc32c() function below.
	plaintextCRC32C = crc32c(plaintext)

	// Build the request.
	req := &kmspb.EncryptRequest{
		Name:            name,
		Plaintext:       plaintext,
		PlaintextCrc32C: plaintextCrc32c,
	}

	// Call the API.
	result, err := client.Encrypt(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to encrypt: %v", err)
	}

	// Optional, but recommended: perform integrity verification on result.
	// For more details on ensuring E2E in-transit integrity to and from Cloud KMS visit:
	// https://cloud.google.com/kms/docs/data-integrity-guidelines
	if result.VerifiedPlaintextCrc32C == false {
		return fmt.Errorf("Encrypt: request corrupted in-transit")
	}
	if crc32c(result.Ciphertext) != result.CiphertextCrc32C {
		return fmt.Errorf("The response received from the server was corrupted in-transit.")
	}
	// End integrity verification

	fmt.Fprintf(w, "Encrypted ciphertext: %s", result.Ciphertext)
	return nil
}

func crc32c(data []byte) uint32 {
	t := crc32.MakeTable(crc32.Castagnoli)
	return crc32.Checksum(data, t)
}

// [END kms_encrypt_symmetric]
