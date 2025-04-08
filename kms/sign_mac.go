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

// [START kms_sign_mac]
import (
	"context"
	"encoding/base64"
	"fmt"
	"io"

	kms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"
)

// signMac will sign a plaintext message using the HMAC algorithm.
func signMac(w io.Writer, name string, data string) error {
	// name := "projects/my-project/locations/us-east1/keyRings/my-key-ring/cryptoKeys/my-key/cryptoKeyVersions/123"
	// data := "my data to sign"

	// Create the client.
	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create kms client: %w", err)
	}
	defer client.Close()

	// Build the request.
	req := &kmspb.MacSignRequest{
		Name: name,
		Data: []byte(data),
	}

	// Generate HMAC of data.
	result, err := client.MacSign(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to hmac sign: %w", err)
	}

	// The data comes back as raw bytes, which may include non-printable
	// characters. This base64-encodes the result so it can be printed below.
	encodedSignature := base64.StdEncoding.EncodeToString(result.Mac)

	fmt.Fprintf(w, "Signature: %s", encodedSignature)
	return nil
}

// [END kms_sign_mac]
