// Copyright 2024 Google LLC
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

package secretmanager

// [START secretmanager_create_update_secret_label]
import (
	"context"
	"fmt"
	"io"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"google.golang.org/genproto/protobuf/field_mask"
)

// createUpdateSecretLabel updates the labels about an existing secret.
// If the label key exists, it updates the label, otherwise it creates a new one.
func createUpdateSecretLabel(w io.Writer, name string) error {
	// name := "projects/my-project/secrets/my-secret"

	labelKey := "labelkey"
	labelValue := "updatedlabelvalue"

	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create secretmanager client: %w", err)
	}
	defer client.Close()

	// Build the request to get the secret.
	req := &secretmanagerpb.GetSecretRequest{
		Name: name,
	}

	// Call the API.
	result, err := client.GetSecret(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to get secret: %w", err)
	}

	labels := result.Labels

	labels[labelKey] = labelValue

	// Build the request to update the secret.
	update_req := &secretmanagerpb.UpdateSecretRequest{
		Secret: &secretmanagerpb.Secret{
			Name:   name,
			Labels: labels,
		},
		UpdateMask: &field_mask.FieldMask{
			Paths: []string{"labels"},
		},
	}

	// Call the API.
	update_result, err := client.UpdateSecret(ctx, update_req)
	if err != nil {
		return fmt.Errorf("failed to update secret: %w", err)
	}
	fmt.Fprintf(w, "Updated secret: %s\n", update_result.Name)
	return nil
}

// [END secretmanager_create_update_secret_label]
