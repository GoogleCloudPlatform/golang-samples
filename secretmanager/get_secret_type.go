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

// [START secretmanager_get_secret_type]
import (
	"context"
	"fmt"
	"io"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

// getSecretType gets and prints the secret type (e.g. CLOUD_SQL_DB_CREDENTIALS,
// ACCESS_KEY, CERTIFICATE, OTHER_DB_CREDENTIALS, OTHER, or
// SECRET_TYPE_UNSPECIFIED for a secret with no type restriction) of the given
// secret.
func getSecretType(w io.Writer, name string) error {
	// name := "projects/my-project/secrets/my-secret"

	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create secretmanager client: %w", err)
	}
	defer client.Close()

	// Build the request.
	req := &secretmanagerpb.GetSecretRequest{
		Name: name,
	}

	// Call the API.
	result, err := client.GetSecret(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to get secret: %w", err)
	}

	fmt.Fprintf(w, "Found secret %s with secret type %s\n", result.Name, result.SecretType)
	return nil
}

// [END secretmanager_get_secret_type]
