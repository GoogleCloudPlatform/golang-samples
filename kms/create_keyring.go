// Copyright 2019 Google LLC
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

// [START kms_create_keyring]
import (
	"context"
	"fmt"
	"io"

	cloudkms "cloud.google.com/go/kms/apiv1"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

// createKeyRing creates a new ring to store keys on KMS.
func createKeyRing(w io.Writer, parent, keyRingID string) error {
	// parent := "projects/PROJECT_ID/locations/global/"
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return fmt.Errorf("cloudkms.NewKeyManagementClient: %v", err)
	}
	// Build the request.
	req := &kmspb.CreateKeyRingRequest{
		Parent:    parent,
		KeyRingId: keyRingID,
	}
	// Call the API.
	result, err := client.CreateKeyRing(ctx, req)
	if err != nil {
		return fmt.Errorf("CreateKeyRing: %v", err)
	}
	fmt.Fprintf(w, "Created key ring: %s", result)
	return nil
}

// [END kms_create_keyring]
