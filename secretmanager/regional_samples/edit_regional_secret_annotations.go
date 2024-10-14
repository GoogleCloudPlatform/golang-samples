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

// [START secretmanager_edit_regional_secret_annotations]
import (
	"context"
	"fmt"
	"io"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"google.golang.org/api/option"
	"google.golang.org/genproto/protobuf/field_mask"
)

// editRegionalSecretAnnotations updates the annotations about an existing secret.
// If the annotation key exists, it updates the annotation, otherwise it creates a new one.
func EditRegionalSecretAnnotation(w io.Writer, projectId, locationId, secretId string) error {
	// name := "projects/my-project/secrets/my-secret"

	annotationKey := "annotationkey"
	annotationValue := "updatedannotationvalue"

	ctx := context.Background()

	//Endpoint to send the request to regional server
	endpoint := fmt.Sprintf("secretmanager.%s.rep.googleapis.com:443", locationId)
	client, err := secretmanager.NewClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		return fmt.Errorf("failed to create secretmanager client: %w", err)
	}
	defer client.Close()

	name := fmt.Sprintf("projects/%s/locations/%s/secrets/%s", projectId, locationId, secretId)

	// Build the request to get the secret.
	req := &secretmanagerpb.GetSecretRequest{
		Name: name,
	}

	result, err := client.GetSecret(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to get secret: %w", err)
	}

	annotations := result.Annotations

	annotations[annotationKey] = annotationValue

	// Build the request to update the secret.
	update_req := &secretmanagerpb.UpdateSecretRequest{
		Secret: &secretmanagerpb.Secret{
			Name:        name,
			Annotations: annotations,
		},
		UpdateMask: &field_mask.FieldMask{
			Paths: []string{"annotations"},
		},
	}

	// Call the API.
	update_result, err := client.UpdateSecret(ctx, update_req)
	if err != nil {
		return fmt.Errorf("failed to update secret: %w", err)
	}
	fmt.Fprintf(w, "Updated secret: %s\n", update_result.Name)
	return nil
}

// [END secretmanager_edit_regional_secret_annotations]
