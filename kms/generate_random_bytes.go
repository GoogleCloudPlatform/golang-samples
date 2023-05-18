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

// [START kms_generate_random_bytes]
import (
	"context"
	"encoding/base64"
	"fmt"
	"io"

	kms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"
)

// generateRandomBytes generates random bytes with entropy sourced from the
// given location.
func generateRandomBytes(w io.Writer, location string, numBytes int32) error {
	// name := "projects/my-project/locations/us-east1"
	// numBytes := 256

	// Create the client.
	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create kms client: %w", err)
	}
	defer client.Close()

	// Build the request.
	req := &kmspb.GenerateRandomBytesRequest{
		Location:        location,
		LengthBytes:     numBytes,
		ProtectionLevel: kmspb.ProtectionLevel_HSM,
	}

	// Generate random bytes.
	result, err := client.GenerateRandomBytes(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// The data comes back as raw bytes, which may include non-printable
	// characters. This base64-encodes the result so it can be printed below.
	encodedData := base64.StdEncoding.EncodeToString(result.Data)

	fmt.Fprintf(w, "Random bytes: %s", encodedData)
	return nil
}

// [END kms_generate_random_bytes]
