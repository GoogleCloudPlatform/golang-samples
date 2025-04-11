// Copyright 2021 Google LLC
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

// [START secretmanager_list_secrets_with_filter]
import (
	"context"
	"fmt"
	"io"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"google.golang.org/api/iterator"
)

// listSecretsWithFilter lists all filter-matching secrets in the given project.
func listSecretsWithFilter(w io.Writer, parent string, filter string) error {
	// parent := "projects/my-project"
	// Follow https://cloud.google.com/secret-manager/docs/filtering
	// for filter syntax and examples.
	// filter := "name:name-substring"

	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create secretmanager client: %w", err)
	}
	defer client.Close()

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
			return fmt.Errorf("failed to list secrets: %w", err)
		}

		fmt.Fprintf(w, "Found secret %s\n", resp.Name)
	}

	return nil
}

// [END secretmanager_list_secrets_with_filter]
