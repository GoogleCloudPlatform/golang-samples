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

// [START kms_update_key_update_labels]
import (
	"context"
	"fmt"
	"io"

	kms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"
	fieldmask "google.golang.org/genproto/protobuf/field_mask"
)

// updateKeyUpdateLabels updates an existing KMS key, adding a new label.
func updateKeyUpdateLabels(w io.Writer, name string) error {
	// name := "projects/my-project/locations/us-east1/keyRings/my-key-ring/cryptoKeys/my-key"

	// Create the client.
	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create kms client: %w", err)
	}
	defer client.Close()

	//
	// Step 1 - get the current set of labels on the key
	//

	// Build the request.
	getReq := &kmspb.GetCryptoKeyRequest{
		Name: name,
	}

	// Call the API.
	result, err := client.GetCryptoKey(ctx, getReq)
	if err != nil {
		return fmt.Errorf("failed to get key: %w", err)
	}

	//
	// Step 2 - add a label to the list of labels
	//

	labels := result.Labels
	labels["new_label"] = "new_value"

	// Build the request.
	updateReq := &kmspb.UpdateCryptoKeyRequest{
		CryptoKey: &kmspb.CryptoKey{
			Name:   name,
			Labels: labels,
		},
		UpdateMask: &fieldmask.FieldMask{
			Paths: []string{"labels"},
		},
	}

	// Call the API.
	result, err = client.UpdateCryptoKey(ctx, updateReq)
	if err != nil {
		return fmt.Errorf("failed to update key: %w", err)
	}

	// Print the labels.
	for k, v := range result.Labels {
		fmt.Fprintf(w, "%s=%s\n", k, v)
	}
	return nil
}

// [END kms_update_key_update_labels]
