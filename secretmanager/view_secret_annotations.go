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

package secretmanager

// [START secretmanager_view_secret_annotations]
import (
	"context"
	"fmt"
	"io"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

// viewSecretAnnotations gets annotations with the given secret.
func viewSecretAnnotations(w io.Writer, secretName string) error {
	// name := "projects/my-project/secrets/my-secret"

	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create secretmanager client: %w", err)
	}
	defer client.Close()

	// Build the request.
	req := &secretmanagerpb.GetSecretRequest{
		Name: secretName,
	}

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

// [END secretmanager_view_secret_annotations]
