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

import (
	"context"
	"fmt"
	"io"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

// createUserManagedReplicationSecret creates a new secret with the given name with user managed replication.
// A secret is a logical wrapper around a collection of secret versions. Secret
// versions hold the actual secret material.
func createUserManagedReplicationSecret(w io.Writer, parent, id string, locations []string) error {
	// parent := "projects/my-project"
	// id := "my-secret"
	// locations := []string{"us-east1", "us-east4", "us-west1"}

	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create secretmanager client: %w", err)
	}
	defer client.Close()

	var replicaLocations []*secretmanagerpb.Replication_UserManaged_Replica
	for _, location := range locations {
		replicaLocations = append(replicaLocations, &secretmanagerpb.Replication_UserManaged_Replica{Location: location})
	}

	// Create the request to create the secret.
	req := &secretmanagerpb.CreateSecretRequest{
		Parent:   parent,
		SecretId: id,
		Secret: &secretmanagerpb.Secret{
			Replication: &secretmanagerpb.Replication{
				Replication: &secretmanagerpb.Replication_UserManaged_{
					UserManaged: &secretmanagerpb.Replication_UserManaged{
						Replicas: replicaLocations,
					},
				},
			},
		},
	}

	// Call the API.
	result, err := client.CreateSecret(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create secret with user managed replication: %w", err)
	}
	fmt.Fprintf(w, "Created secret with user managed replication: %s\n", result.Name)
	return nil
}
