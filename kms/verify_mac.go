// Copyright 2021 Google LLC
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

// [START kms_verify_mac]
import (
	"context"
	"fmt"
	"io"

	kms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"
)

// verifyMac will verify a previous HMAC signature.
func verifyMac(w io.Writer, name string, data, signature []byte) error {
	// name := "projects/my-project/locations/us-east1/keyRings/my-key-ring/cryptoKeys/my-key/cryptoKeyVersions/123"
	// data := "my previous data"
	// signature := []byte("...")  // Response from a sign request

	// Create the client.
	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create kms client: %w", err)
	}
	defer client.Close()

	// Build the request.
	req := &kmspb.MacVerifyRequest{
		Name: name,
		Data: data,
		Mac:  signature,
	}

	// Verify the signature.
	result, err := client.MacVerify(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to verify signature: %w", err)
	}

	fmt.Fprintf(w, "Verified: %t", result.Success)
	return nil
}

// [END kms_verify_mac]
