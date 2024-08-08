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

// [START secretmanager_list_regional_secrets_with_filter]
import (
	"context"
	"fmt"
	"io"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// listSecretsWithFilter lists all filter-matching secrets in the given project.
func ListRegionalSecretsWithFilter(w io.Writer, projectId, locationId string, filter string) error {
	// parent := "projects/my-project/locations/my-location"
	// Follow https://cloud.google.com/secret-manager/docs/filtering
	// for filter syntax and examples.
	// filter := "name:name-substring"

	// Create the client.
	ctx := context.Background()
	//Endpoint to send the request to regional server
	endpoint := fmt.Sprintf("secretmanager.%s.rep.googleapis.com:443", locationId)
	client, err := secretmanager.NewClient(ctx, option.WithEndpoint(endpoint))

	if err != nil {
		return fmt.Errorf("failed to create regional secretmanager client: %w", err)
	}
	defer client.Close()

	parent := fmt.Sprintf("projects/%s/locations/%s", projectId, locationId)
	// Build the request.
	req := &secretmanagerpb.ListSecretsRequest{
		Parent: parent,
		Filter: filter,
	}

	// Call the API.
	it := client.ListSecrets(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			return fmt.Errorf("failed to list regional secrets: %w", err)
		}

		fmt.Fprintf(w, "Found regional secret %s\n", resp.Name)
	}

	return nil
}

// [END secretmanager_list_regional_secrets_with_filter]
