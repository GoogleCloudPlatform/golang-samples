// Copyright 2022 Google LLC
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

// [START secretmanager_update_secret_with_alias]
import (
	"context"
	"fmt"
	"io"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"google.golang.org/genproto/protobuf/field_mask"
)

// updateSecret updates the alias map on an existing secret.
func updateSecretWithAlias(w io.Writer, name string) error {
	// name := "projects/my-project/secrets/my-secret"

	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create secretmanager client: %w", err)
	}
	defer client.Close()

	// Build the request.
	req := &secretmanagerpb.UpdateSecretRequest{
		Secret: &secretmanagerpb.Secret{
			Name: name,
			VersionAliases: map[string]int64{
				"test": 1,
			},
		},
		UpdateMask: &field_mask.FieldMask{
			Paths: []string{"version_aliases"},
		},
	}

	// Call the API.
	result, err := client.UpdateSecret(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to update secret: %w", err)
	}
	fmt.Fprintf(w, "Updated secret: %s\n", result.Name)
	return nil
}

// [END secretmanager_update_secret_with_alias]
