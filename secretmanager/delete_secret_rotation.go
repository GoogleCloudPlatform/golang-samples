// Copyright 2026 Google LLC
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

// [START secretmanager_delete_secret_rotation]

import (
	"context"
	"fmt"
	"io"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// deleteSecretRotation removes the rotation configuration from a secret.
func deleteSecretRotation(w io.Writer, secretName string) error {
	// secretName := "projects/my-project/secrets/my-secret"

	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create secretmanager client: %w", err)
	}
	defer client.Close()

	req := &secretmanagerpb.UpdateSecretRequest{
		Secret: &secretmanagerpb.Secret{
			Name: secretName,
		},
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: []string{"rotation"},
		},
	}

	result, err := client.UpdateSecret(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to update secret: %w", err)
	}

	fmt.Fprintf(w, "Removed rotation from secret %s\n", result.Name)
	return nil
}

// [END secretmanager_delete_secret_rotation]
