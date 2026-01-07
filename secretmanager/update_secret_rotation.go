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

package secretmanager

import (
	"context"
	"fmt"
	"io"
	"time"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// [START secretmanager_update_secret_rotation_period]

// updateSecretRotationPeriod updates the rotation period of a secret.
func updateSecretRotationPeriod(w io.Writer, projectID, secretID string, rotationPeriod time.Duration) error {
	// projectID := "my-project"
	// secretID := "my-secret"
	// rotationPeriod := time.Hour * 24 * 7

	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create secretmanager client: %w", err)
	}
	defer client.Close()

	// Build the request.
	req := &secretmanagerpb.UpdateSecretRequest{
		Secret: &secretmanagerpb.Secret{
			Name: fmt.Sprintf("projects/%s/secrets/%s", projectID, secretID),
			Rotation: &secretmanagerpb.Rotation{
				RotationPeriod: durationpb.New(rotationPeriod),
			},
		},
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: []string{"rotation.rotation_period"},
		},
	}

	// Call the API.
	result, err := client.UpdateSecret(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to update secret: %w", err)
	}

	fmt.Fprintf(w, "Updated secret %s rotation period to %v\n", result.Name, result.Rotation.RotationPeriod.AsDuration())
	return nil
}

// [END secretmanager_update_secret_rotation_period]
