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

// [START secretmanager_rotate_regional_secret]
import (
	"context"
	"fmt"
	"io"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"google.golang.org/api/option"
)

// RotateRegionalSecret triggers a managed rotation for a Cloud SQL DB
// credentials secret. Managed rotation must already be enabled on the secret
// (see EnableRegionalSecretManagedRotation). Each call generates a new
// password, updates the Cloud SQL user, and adds the result as a new secret
// version.
func RotateRegionalSecret(w io.Writer, projectId, locationId, secretId string) error {
	// parent := "projects/my-project/locations/my-location/secrets/my-secret"

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

	// Build the request.
	req := &secretmanagerpb.RotateSecretRequest{
		Parent: parent,
	}

	// Call the API.
	result, err := client.RotateSecret(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to rotate regional secret: %w", err)
	}
	fmt.Fprintf(w, "Rotated secret, created secret version: %s\n", result.Name)

	return nil
}

// [END secretmanager_rotate_regional_secret]
