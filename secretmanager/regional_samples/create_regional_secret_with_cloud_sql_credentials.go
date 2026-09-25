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

// [START secretmanager_create_regional_secret_with_cloud_sql_credentials]
import (
	"context"
	"fmt"
	"io"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"google.golang.org/api/option"
)

// CreateRegionalSecretWithCloudSQLCredentials creates a new secret with the
// Cloud SQL DB credentials secret type. This type is required to enable
// Secret Manager's automatic rotation of Cloud SQL passwords. It can only be
// set when the secret is created, and the secret's location must match the
// region of the target Cloud SQL instance.
func CreateRegionalSecretWithCloudSQLCredentials(w io.Writer, projectId, locationId, secretId string) error {
	// parent := "projects/my-project/locations/my-location"
	// secretId := "my-secret"

	// Create the client.
	ctx := context.Background()

	// Endpoint to send the request to regional server
	endpoint := fmt.Sprintf("secretmanager.%s.rep.googleapis.com:443", locationId)
	client, err := secretmanager.NewClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		return fmt.Errorf("failed to create regional secretmanager client: %w", err)
	}
	defer client.Close()

	parent := fmt.Sprintf("projects/%s/locations/%s", projectId, locationId)

	// Build the request.
	req := &secretmanagerpb.CreateSecretRequest{
		Parent:   parent,
		SecretId: secretId,
		Secret: &secretmanagerpb.Secret{
			SecretType: secretmanagerpb.Secret_CLOUD_SQL_DB_CREDENTIALS,
		},
	}

	// Call the API.
	result, err := client.CreateSecret(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create regional secret: %w", err)
	}
	fmt.Fprintf(w, "Created secret: %s\n", result.Name)

	// This built-in identity is what you grant Cloud SQL IAM permissions to,
	// so that Secret Manager can rotate the database password on its behalf.
	fmt.Fprintf(w, "Grant this identity Cloud SQL IAM permissions to enable rotation: %s\n",
		result.GetPolicyMember().GetIamPolicyUidPrincipal())

	return nil
}

// [END secretmanager_create_regional_secret_with_cloud_sql_credentials]
