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

// [START secretmanager_enable_regional_secret_managed_rotation]
import (
	"context"
	"fmt"
	"io"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"google.golang.org/api/option"
)

// EnableRegionalSecretManagedRotation enables managed rotation for a Cloud
// SQL DB credentials secret. This links the secret to a Cloud SQL instance
// and database user, and can only be called once per secret. It adds the
// secret's first version and sets the matching password on the Cloud SQL
// user, taking the place of a manually added secret version, which this
// secret type doesn't support. Afterwards, use RotateRegionalSecret to
// trigger further rotations.
//
// instanceId is the bare Cloud SQL instance ID (e.g. "my-instance") -- not a
// connection name. Neither the project nor the region should be included:
// passing "PROJECT_ID:INSTANCE_ID" (as gcloud's own
// `enable-managed-rotation --help` examples misleadingly show) or the full
// "PROJECT_ID:LOCATION_ID:INSTANCE_ID" connection name both fail -- the
// service already knows the project from the secret's own path, and prepends
// it internally, so a qualified value ends up double-prefixed.
func EnableRegionalSecretManagedRotation(w io.Writer, projectId, locationId, secretId, instanceId, username string) error {
	// parent := "projects/my-project/locations/my-location/secrets/my-secret"
	// instanceId := "my-cloud-sql-instance"
	// username := "my-db-user"

	// Create the client.
	ctx := context.Background()

	// Endpoint to send the request to regional server
	endpoint := fmt.Sprintf("secretmanager.%s.rep.googleapis.com:443", locationId)
	client, err := secretmanager.NewClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		return fmt.Errorf("failed to create regional secretmanager client: %w", err)
	}
	defer client.Close()

	// Despite the field name, parent holds the full secret resource name,
	// not a collection parent.
	parent := fmt.Sprintf("projects/%s/locations/%s/secrets/%s", projectId, locationId, secretId)

	// Build the request. Leaving Password unset lets Secret Manager generate
	// a secure password itself.
	req := &secretmanagerpb.EnableManagedRotationRequest{
		Parent: parent,
		Credentials: &secretmanagerpb.EnableManagedRotationRequest_CloudSqlSingleUserCredentials{
			CloudSqlSingleUserCredentials: &secretmanagerpb.EnableManagedRotationRequest_CloudSQLSingleUserCredentials{
				InstanceId: instanceId,
				Username:   username,
			},
		},
	}

	// Call the API.
	result, err := client.EnableManagedRotation(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to enable managed rotation for regional secret: %w", err)
	}
	fmt.Fprintf(w, "Enabled managed rotation, created secret version: %s\n", result.Name)

	return nil
}

// [END secretmanager_enable_regional_secret_managed_rotation]
