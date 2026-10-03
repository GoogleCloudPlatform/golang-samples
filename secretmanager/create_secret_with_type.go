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

// [START secretmanager_create_secret_with_type]
import (
	"context"
	"fmt"
	"io"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

// createSecretWithType creates a new secret with the given secret type
// restriction (e.g. ACCESS_KEY, CERTIFICATE, OTHER_DB_CREDENTIALS, or OTHER --
// use CLOUD_SQL_DB_CREDENTIALS only for a secret that will go through
// EnableManagedRotation, which additionally requires a regional secret; see
// the regional_samples package). Unlike CLOUD_SQL_DB_CREDENTIALS, these
// other secret types are plain metadata tags: they don't require any
// additional credentials payload at creation time.
func createSecretWithType(w io.Writer, parent, id string, secretType secretmanagerpb.Secret_SecretType) error {
	// parent := "projects/my-project"
	// id := "my-secret"
	// secretType := secretmanagerpb.Secret_ACCESS_KEY

	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create secretmanager client: %w", err)
	}
	defer client.Close()

	// Build the request.
	req := &secretmanagerpb.CreateSecretRequest{
		Parent:   parent,
		SecretId: id,
		Secret: &secretmanagerpb.Secret{
			SecretType: secretType,
			Replication: &secretmanagerpb.Replication{
				Replication: &secretmanagerpb.Replication_Automatic_{
					Automatic: &secretmanagerpb.Replication_Automatic{},
				},
			},
		},
	}

	// Call the API.
	result, err := client.CreateSecret(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create secret: %w", err)
	}
	fmt.Fprintf(w, "Created secret with secret type: %s\n", result.Name)
	return nil
}

// [END secretmanager_create_secret_with_type]
