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

package regional_secretmanager

import (
	"context"
	"fmt"
	"io"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"google.golang.org/api/option"
)

// [START secretmanager_create_regional_secret_with_cmek]

// createSecretWithCMEK creates a new secret encrypted with a customer-managed key.
func createRegionalSecretWithCMEK(w io.Writer, projectID, secretID, locationID, kmsKeyName string) error {
	// projectID := "my-project"
	// secretID := "my-secret-with-cmek"
	// kmsKeyName := "projects/my-project/locations/{locationID}/keyRings/{keyringname}/cryptoKeys/{keyname}"

	ctx := context.Background()
	endpoint := fmt.Sprintf("secretmanager.%s.rep.googleapis.com:443", locationID)
	client, err := secretmanager.NewClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		return fmt.Errorf("failed to create secretmanager client: %w", err)
	}
	defer client.Close()

	parent := fmt.Sprintf("projects/%s/locations/%s", projectID, locationID)

	req := &secretmanagerpb.CreateSecretRequest{
		Parent:   parent,
		SecretId: secretID,
		Secret: &secretmanagerpb.Secret{
			CustomerManagedEncryption: &secretmanagerpb.CustomerManagedEncryption{
				KmsKeyName: kmsKeyName,
			},
		},
	}

	secret, err := client.CreateSecret(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create secret: %w", err)
	}

	fmt.Fprintf(w, "Created secret %s with CMEK key %s\n", secret.Name, kmsKeyName)
	return nil
}

// [END secretmanager_create_regional_secret_with_cmek]
