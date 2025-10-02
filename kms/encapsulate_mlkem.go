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

// [START kms_encapsulate_mlkem]
import (
	"context"
	"crypto/mlkem"
	"fmt"
	"hash/crc32"
	"io"

	kms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"
)

// encapsulateMLKEM demonstrates how to encapsulate a shared secret using an ML-KEM-768 public key
// from Cloud KMS.
func encapsulateMLKEM(w io.Writer, keyVersionName string) error {
	// keyVersionName := "projects/my-project/locations/us-east1/keyRings/my-key-ring/cryptoKeys/my-key/cryptoKeyVersions/1"

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

	// Build the request to get the public key in NIST PQC format.
	req := &kmspb.GetPublicKeyRequest{
		Name:            keyVersionName,
		PublicKeyFormat: kmspb.PublicKey_NIST_PQC,
	}

	// Call the API to get the public key.
	response, err := client.GetPublicKey(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to get public key: %w", err)
	}

	// Optional, but recommended: perform integrity verification on the response.
	// For more details on ensuring E2E in-transit integrity to and from Cloud KMS visit:
	// https://cloud.google.com/kms/docs/data-integrity-guidelines
	if response.GetName() != req.GetName() {
		return fmt.Errorf("GetPublicKey: request corrupted in-transit")
	}
	if response.GetPublicKeyFormat() != req.GetPublicKeyFormat() {
		return fmt.Errorf("GetPublicKey: request corrupted in-transit")
	}
	if int64(crc32c(response.GetPublicKey().GetData())) != response.GetPublicKey().GetCrc32CChecksum().GetValue() {
		return fmt.Errorf("GetPublicKey: response corrupted in-transit")
	}

	// Use the public key with crypto/mlkem to encapsulate a shared secret.
	ek, err := mlkem.NewEncapsulationKey768(response.GetPublicKey().GetData())
	if err != nil {
		return fmt.Errorf("NewEncapsulationKey768: %w", err)
	}
	sharedSecret, ciphertext := ek.Encapsulate()

	fmt.Fprintf(w, "Encapsulated ciphertext: %x\n", ciphertext)
	fmt.Fprintf(w, "Shared secret: %x\n", sharedSecret)
	return nil
}

// [END kms_encapsulate_mlkem]
