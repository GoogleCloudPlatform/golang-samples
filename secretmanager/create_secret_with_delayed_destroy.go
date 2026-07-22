// Copyright 2025 Google LLC
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

// [START secretmanager_create_secret_with_delayed_destroy]
import (
	"context"
	"fmt"
	"io"
	"time"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"google.golang.org/protobuf/types/known/durationpb"
)

// createSecretWithDelayedDestroy creates a new secret with the given name
// and version destroy ttl. A secret is a logical wrapper around a collection
// of secret versions. Secret versions hold the actual secret material.
func createSecretWithDelayedDestroy(w io.Writer, parent, secretID string, versionDestroyTtl int) error {
	// parent := "projects/my-project"
	// secretID := "my-secret"
	// versionDestroyTtl := 86400

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
		SecretId: secretID,
		Secret: &secretmanagerpb.Secret{
			Replication: &secretmanagerpb.Replication{
				Replication: &secretmanagerpb.Replication_Automatic_{
					Automatic: &secretmanagerpb.Replication_Automatic{},
				},
			},
			VersionDestroyTtl: durationpb.New(time.Duration(versionDestroyTtl) * time.Second),
		},
	}

	// Call the API.
	result, err := client.CreateSecret(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create secret: %w", err)
	}
	fmt.Fprintf(w, "Created secret with version destroy ttl: %s\n", result.Name)
	return nil
}

// [END secretmanager_create_secret_with_delayed_destroy]
