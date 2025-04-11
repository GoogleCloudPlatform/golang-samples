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

// [START kms_get_public_key]
import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"hash/crc32"
	"io"

	kms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"
)

// getPublicKey retrieves the public key from an asymmetric key pair on
// Cloud KMS.
func getPublicKey(w io.Writer, name string) error {
	// name := "projects/my-project/locations/us-east1/keyRings/my-key-ring/cryptoKeys/my-key/cryptoKeyVersions/123"

	// Create the client.
	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create kms client: %w", err)
	}
	defer client.Close()

	// Build the request.
	req := &kmspb.GetPublicKeyRequest{
		Name: name,
	}

	// Call the API.
	result, err := client.GetPublicKey(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to get public key: %w", err)
	}

	// The 'Pem' field is the raw string representation of the public key.
	// Convert 'Pem' into bytes for further processing.
	key := []byte(result.Pem)

	// Optional, but recommended: perform integrity verification on result.
	// For more details on ensuring E2E in-transit integrity to and from Cloud KMS visit:
	// https://cloud.google.com/kms/docs/data-integrity-guidelines
	crc32c := func(data []byte) uint32 {
		t := crc32.MakeTable(crc32.Castagnoli)
		return crc32.Checksum(data, t)
	}
	if int64(crc32c(key)) != result.PemCrc32C.Value {
		return fmt.Errorf("getPublicKey: response corrupted in-transit")
	}

	// Optional - parse the public key. This transforms the string key into a Go
	// PublicKey.
	block, _ := pem.Decode(key)
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse public key: %w", err)
	}
	fmt.Fprintf(w, "Retrieved public key: %v\n", publicKey)
	return nil
}

// [END kms_get_public_key]
