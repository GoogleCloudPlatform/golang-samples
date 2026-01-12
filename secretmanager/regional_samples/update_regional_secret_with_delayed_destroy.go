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

import (
	"context"
	"fmt"
	"io"
	"time"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/types/known/durationpb"
	fieldmaskpb "google.golang.org/protobuf/types/known/fieldmaskpb"
)

// [START secretmanager_update_regional_secret_with_delayed_destroy]

// updateRegionalSecretWithDelayedDestroy updates the version destroy TTL of an existing regional secret.
func updateRegionalSecretWithDelayedDestroy(w io.Writer, projectID, secretID, locationID string) error {
	// projectID := "my-project"
	// secretID := "my-secret-with-ttl"
	// locationID := "us-central1"
	newDelayedDestroy := time.Hour * 48

	ctx := context.Background()
	endpoint := fmt.Sprintf("secretmanager.%s.rep.googleapis.com:443", locationID)
	client, err := secretmanager.NewClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		return fmt.Errorf("failed to create secretmanager client: %w", err)
	}
	defer client.Close()

	req := &secretmanagerpb.UpdateSecretRequest{
		Secret: &secretmanagerpb.Secret{
			Name:              fmt.Sprintf("projects/%s/locations/%s/secrets/%s", projectID, locationID, secretID),
			VersionDestroyTtl: durationpb.New(newDelayedDestroy),
		},
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: []string{"version_destroy_ttl"},
		},
	}

	secret, err := client.UpdateSecret(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to update secret: %w", err)
	}

	fmt.Fprintf(w, "Updated secret %s with version destroy TTL %v\n", secret.Name, newDelayedDestroy)
	return nil
}

// [END secretmanager_update_regional_secret_with_delayed_destroy]
