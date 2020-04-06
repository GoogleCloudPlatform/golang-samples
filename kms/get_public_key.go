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
	"io"

	kms "cloud.google.com/go/kms/apiv1"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

// getPublicKey retrieves the public key from an asymmetric key pair on
// Cloud KMS.
func getPublicKey(w io.Writer, name string) error {
	// parent := "projects/my-project/locations/us-east1/keyRings/my-key-ring/cryptoKeys/my-key/cryptoKeyVersions/123"

	// Create the client.
	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create kms client: %v", err)
	}

	// Build the request.
	req := &kmspb.GetPublicKeyRequest{
		Name: name,
	}

	// Call the API.
	result, err := client.GetPublicKey(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to get public key: %v", err)
	}

	// The 'Pem' field is the raw string representation of the public key.
	key := result.Pem

	//
	// Optional - parse the public key. This transforms the string key into a Go
	// PublicKey.
	//

	block, _ := pem.Decode([]byte(key))
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse public key: %v", err)
	}
	fmt.Fprintf(w, "Retrieved public key: %v\n", publicKey)
	return nil
}

// [END kms_get_public_key]
