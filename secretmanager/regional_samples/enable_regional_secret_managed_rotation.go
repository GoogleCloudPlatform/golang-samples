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

// EnableRegionalSecretManagedRotation enables managed rotation of a
// CLOUD_SQL_DB_CREDENTIALS typed secret. It validates and enables the
// rotation, adding a version and sets the passed password.
// Note: AddSecretVersion is disabled on the CLOUD_SQL_DB_CREDENTIALS
// currently and for any necessary manual rotations please trigger
// RotateRegionalSecret.
func EnableRegionalSecretManagedRotation(w io.Writer, projectId, locationId, secretId, instanceId, username string) error {
	// parent := "projects/my-project/locations/my-location/secrets/my-secret"
	// instanceId := "my-instance"
	// username := "my-user"

	// Create the client.
	ctx := context.Background()

	// Endpoint to send the request to regional server
	endpoint := fmt.Sprintf("secretmanager.%s.rep.googleapis.com:443", locationId)
	client, err := secretmanager.NewClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		return fmt.Errorf("failed to create regional secretmanager client: %w", err)
	}
	defer client.Close()

	parent := fmt.Sprintf("projects/%s/locations/%s/secrets/%s", projectId, locationId, secretId)

	// Build the request.
	req := &secretmanagerpb.EnableManagedRotationRequest{
		Parent: parent,
		Credentials: &secretmanagerpb.EnableManagedRotationRequest_CloudSqlSingleUserCredentials{
			CloudSqlSingleUserCredentials: &secretmanagerpb.EnableManagedRotationRequest_CloudSQLSingleUserCredentials{
				InstanceId: instanceId,
				Username:   username,
			},
		},
	}

	// Enable managed rotation.
	result, err := client.EnableManagedRotation(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to enable managed rotation for regional secret: %w", err)
	}
	fmt.Fprintf(w, "Enabled managed rotation, created secret version: %s\n", result.Name)

	return nil
}

// [END secretmanager_enable_regional_secret_managed_rotation]
