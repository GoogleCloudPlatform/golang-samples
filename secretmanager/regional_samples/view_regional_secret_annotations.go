// Copyright 2024 Google LLC
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

// [START secretmanager_view_regional_secret_annotations]
import (
	"context"
	"fmt"
	"io"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"google.golang.org/api/option"
)

// viewRegionalSecretAnnotations gets annotations with the given secret.
func ViewRegionalSecretAnnotations(w io.Writer, projectId, locationId, id string) error {
	// name := "projects/my-project/secrets/my-secret"

	// Create the client.
	ctx := context.Background()

	//Endpoint to send the request to regional server
	endpoint := fmt.Sprintf("secretmanager.%s.rep.googleapis.com:443", locationId)
	client, err := secretmanager.NewClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		return fmt.Errorf("failed to create secretmanager client: %w", err)
	}
	defer client.Close()

	name := fmt.Sprintf("projects/%s/locations/%s/secrets/%s", projectId, locationId, id)

	// Build the request.
	req := &secretmanagerpb.GetSecretRequest{
		Name: name,
	}

	// Call the API.
	result, err := client.GetSecret(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to get secret: %w", err)
	}

	annotations := result.Annotations
	fmt.Fprintf(w, "Found secret %s\n", result.Name)

	for key, value := range annotations {
		fmt.Fprintf(w, "Annotations key %s : Annotations Value %s", key, value)
	}
	return nil
}

// [END secretmanager_view_regional_secret_annotations]
