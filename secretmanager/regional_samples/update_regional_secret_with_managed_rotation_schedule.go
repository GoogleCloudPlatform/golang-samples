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

// [START secretmanager_update_regional_secret_with_managed_rotation_schedule]
import (
	"context"
	"fmt"
	"io"
	"time"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"google.golang.org/api/option"
	"google.golang.org/genproto/protobuf/field_mask"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// UpdateRegionalSecretWithManagedRotationSchedule updates the rotation
// schedule of a CLOUD_SQL_DB_CREDENTIALS typed secret.
func UpdateRegionalSecretWithManagedRotationSchedule(w io.Writer, projectId, locationId, secretId string, rotationPeriod time.Duration) error {
	// name := "projects/my-project/locations/my-location/secrets/my-secret"
	// rotationPeriod := 24 * time.Hour

	// Create the client.
	ctx := context.Background()
	// Endpoint to send the request to regional server
	endpoint := fmt.Sprintf("secretmanager.%s.rep.googleapis.com:443", locationId)
	client, err := secretmanager.NewClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		return fmt.Errorf("failed to create regional secretmanager client: %w", err)
	}
	defer client.Close()

	name := fmt.Sprintf("projects/%s/locations/%s/secrets/%s", projectId, locationId, secretId)

	// Build the request.
	// The rotation schedule of a CLOUD_SQL_DB_CREDENTIALS secret can be set
	// before or after enabling managed rotation; EnableManagedRotation does
	// not need to be called first. Other secret types also support a
	// rotation schedule, but only when Pub/Sub topics are configured.
	// Pub/Sub topics are not required for CLOUD_SQL_DB_CREDENTIALS.
	// next_rotation_time and rotation_period must be set together.
	req := &secretmanagerpb.UpdateSecretRequest{
		Secret: &secretmanagerpb.Secret{
			Name: name,
			Rotation: &secretmanagerpb.Rotation{
				NextRotationTime: timestamppb.New(time.Now().Add(rotationPeriod)),
				RotationPeriod:   durationpb.New(rotationPeriod),
			},
		},
		UpdateMask: &field_mask.FieldMask{
			// Mask only the rotation subfields being set.
			Paths: []string{"rotation.next_rotation_time", "rotation.rotation_period"},
		},
	}

	// Call the API.
	result, err := client.UpdateSecret(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to update regional secret's rotation schedule: %w", err)
	}
	fmt.Fprintf(w, "Updated regional secret rotation schedule: %s\n", result.Name)
	return nil
}

// [END secretmanager_update_regional_secret_with_managed_rotation_schedule]
