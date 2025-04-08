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

// [START kms_create_key_ring]
import (
	"context"
	"fmt"
	"io"

	kms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"
)

// createKeyRing creates a new ring to store keys on KMS.
func createKeyRing(w io.Writer, parent, id string) error {
	// parent := "projects/PROJECT_ID/locations/global"
	// id := "my-key-ring"

	// Create the client.
	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create kms client: %w", err)
	}
	defer client.Close()

	// Build the request.
	req := &kmspb.CreateKeyRingRequest{
		Parent:    parent,
		KeyRingId: id,
	}

	// Call the API.
	result, err := client.CreateKeyRing(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create key ring: %w", err)
	}
	fmt.Fprintf(w, "Created key ring: %s\n", result.Name)
	return nil
}

// [END kms_create_key_ring]
