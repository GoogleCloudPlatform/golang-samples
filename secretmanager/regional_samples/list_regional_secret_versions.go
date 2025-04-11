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

// [START secretmanager_list_regional_secret_versions]
import (
	"context"
	"fmt"
	"io"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// listSecretVersions lists all secret versions in the given secret and their
// metadata.
func ListRegionalSecretVersions(w io.Writer, projectId, locationId, secretId string) error {
	// parent := "projects/my-project/locations/my-location/secrets/my-secret"

	// Create the client.
	ctx := context.Background()
	//Endpoint to send the request to regional server
	endpoint := fmt.Sprintf("secretmanager.%s.rep.googleapis.com:443", locationId)
	client, err := secretmanager.NewClient(ctx, option.WithEndpoint(endpoint))

	if err != nil {
		return fmt.Errorf("failed to create regional secretmanager client: %w", err)
	}
	defer client.Close()

	parent := fmt.Sprintf("projects/%s/locations/%s/secrets/%s", projectId, locationId, secretId)
	// Build the request.
	req := &secretmanagerpb.ListSecretVersionsRequest{
		Parent: parent,
	}

	// Call the API.
	it := client.ListSecretVersions(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			return fmt.Errorf("failed to list regional secret versions: %w", err)
		}

		fmt.Fprintf(w, "Found regional secret version %s with state %s\n",
			resp.Name, resp.State)
	}

	return nil
}

// [END secretmanager_list_regional_secret_versions]
