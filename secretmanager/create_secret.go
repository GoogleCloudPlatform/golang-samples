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

package secretmanager

// [START secretmanager_create_secret]
import (
	"context"
	"fmt"
	"io"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

// createSecret creates a new secret with the given name. A secret is a logical
// wrapper around a collection of secret versions. Secret versions hold the
// actual secret material.
func createSecret(w io.Writer, parent, id string, ttl *secretmanagerpb.Secret_Ttl) (*secretmanagerpb.Secret, error) {
	// parent := "projects/my-project"
	// id := "my-secret"

	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create secretmanager client: %w", err)
	}
	defer client.Close()

	// Build the request.
	req := &secretmanagerpb.CreateSecretRequest{
		Parent:   parent,
		SecretId: id,
		Secret: &secretmanagerpb.Secret{
			Replication: &secretmanagerpb.Replication{
				Replication: &secretmanagerpb.Replication_Automatic_{
					Automatic: &secretmanagerpb.Replication_Automatic{},
				},
			},
			Expiration: ttl,
		},
	}

	// Call the API.
	result, err := client.CreateSecret(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create secret: %w", err)
	}
	fmt.Fprintf(w, "Created secret: %s\n", result.Name)
	return result, nil
}

// [END secretmanager_create_secret]
