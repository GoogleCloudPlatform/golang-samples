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

package secretmanager

// [START secretmanager_delete_secret_with_etag]
import (
	"context"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

// deleteSecretWithEtag deletes the secret with the given name and all of its versions.
func deleteSecretWithEtag(name, etag string) error {
	// name := "projects/my-project/secrets/my-secret"
	// etag := `"123"`

	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create secretmanager client: %w", err)
	}
	defer client.Close()

	// Build the request.
	req := &secretmanagerpb.DeleteSecretRequest{
		Name: name,
		Etag: etag,
	}

	// Call the API.
	if err := client.DeleteSecret(ctx, req); err != nil {
		return fmt.Errorf("failed to delete secret: %w", err)
	}
	return nil
}

// [END secretmanager_delete_secret_with_etag]
