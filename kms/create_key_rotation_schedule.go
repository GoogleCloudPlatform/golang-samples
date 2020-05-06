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

// [START kms_create_key_rotation_schedule]
import (
	"context"
	"fmt"
	"io"
	"time"

	kms "cloud.google.com/go/kms/apiv1"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/timestamp"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

// createKeyRotationSchedule creates a key with a rotation schedule.
func createKeyRotationSchedule(w io.Writer, parent, id string) error {
	// name := "projects/my-project/locations/us-east1/keyRings/my-key-ring"
	// id := "my-key-with-rotation-schedule"

	// Create the client.
	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create kms client: %v", err)
	}

	// Build the request.
	req := &kmspb.CreateCryptoKeyRequest{
		Parent:      parent,
		CryptoKeyId: id,
		CryptoKey: &kmspb.CryptoKey{
			Purpose: kmspb.CryptoKey_ENCRYPT_DECRYPT,
			VersionTemplate: &kmspb.CryptoKeyVersionTemplate{
				Algorithm: kmspb.CryptoKeyVersion_GOOGLE_SYMMETRIC_ENCRYPTION,
			},

			// Rotate the key every 30 days
			RotationSchedule: &kmspb.CryptoKey_RotationPeriod{
				RotationPeriod: &duration.Duration{
					Seconds: int64(60 * 60 * 24 * 30), // 30 days
				},
			},

			// Start the first rotation in 24 hours
			NextRotationTime: &timestamp.Timestamp{
				Seconds: time.Now().Add(24 * time.Hour).Unix(),
			},
		},
	}

	// Call the API.
	result, err := client.CreateCryptoKey(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create key: %v", err)
	}
	fmt.Fprintf(w, "Created key: %s\n", result.Name)
	return nil
}

// [END kms_create_key_rotation_schedule]
