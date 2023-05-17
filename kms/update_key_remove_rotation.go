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

// [START kms_update_key_remove_rotation_schedule]
import (
	"context"
	"fmt"
	"io"

	kms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"
	fieldmask "google.golang.org/genproto/protobuf/field_mask"
)

// removeRotationSchedule updates a key to remove a rotation schedule, if one
// exists.
func removeRotationSchedule(w io.Writer, name string) error {
	// name := "projects/my-project/locations/us-east1/keyRings/my-key-ring/cryptoKeys/my-key"

	// Create the client.
	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create kms client: %w", err)
	}
	defer client.Close()

	// Build the request.
	req := &kmspb.UpdateCryptoKeyRequest{
		CryptoKey: &kmspb.CryptoKey{
			// Provide the name of the key to update
			Name: name,

			// Remove any rotation fields.
			RotationSchedule: nil,
			NextRotationTime: nil,
		},
		UpdateMask: &fieldmask.FieldMask{
			Paths: []string{"rotation_period", "next_rotation_time"},
		},
	}

	// Call the API.
	result, err := client.UpdateCryptoKey(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to update key: %w", err)
	}
	fmt.Fprintf(w, "Updated key: %s\n", result.Name)
	return nil
}

// [END kms_update_key_remove_rotation_schedule]
