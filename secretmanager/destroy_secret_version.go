// Copyright 2019 Google LLC
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

// [START secretmanager_destroy_secret_version]
import (
	"context"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

// destroySecretVersion destroys the given secret version, making the payload
// irrecoverable. Other secrets versions are unaffected.
func destroySecretVersion(name string) error {
	// name := "projects/my-project/secrets/my-secret/versions/5"

	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create secretmanager client: %w", err)
	}
	defer client.Close()

	// Build the request.
	req := &secretmanagerpb.DestroySecretVersionRequest{
		Name: name,
	}

	// Call the API.
	if _, err := client.DestroySecretVersion(ctx, req); err != nil {
		return fmt.Errorf("failed to destroy secret version: %w", err)
	}
	return nil
}

// [END secretmanager_destroy_secret_version]
