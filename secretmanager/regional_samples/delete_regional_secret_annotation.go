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

// [START secretmanager_delete_regional_secret_annotation]

import (
	"context"
	"fmt"
	"io"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// deleteRegionalSecretAnnotation deletes an annotation on the given secret.
func deleteRegionalSecretAnnotation(w io.Writer, secretName, locationID, annotationKey string) error {
	// secretName := "projects/my-project/locations/my-location/secrets/my-secret"
	// locationID := "us-central1"
	// annotationKey := "annotationkey"

	// Create the client.
	ctx := context.Background()
	//Endpoint to send the request to regional server
	endpoint := fmt.Sprintf("secretmanager.%s.rep.googleapis.com:443", locationID)
	client, err := secretmanager.NewClient(ctx, option.WithEndpoint(endpoint))

	if err != nil {
		return fmt.Errorf("failed to create secretmanager client: %w", err)
	}
	defer client.Close()

	// Get the secret to access annotations.
	getRequest := &secretmanagerpb.GetSecretRequest{
		Name: secretName,
	}

	result, err := client.GetSecret(ctx, getRequest)
	if err != nil {
		return fmt.Errorf("failed to get secret: %w", err)
	}

	// Return if annotation to delete does not exist.
	if _, ok := result.Annotations[annotationKey]; !ok {
		return fmt.Errorf("annotation %s not found on secret %s", annotationKey, secretName)
	}

	// Remove annotation.
	delete(result.Annotations, annotationKey)

	// Build request to update secret.
	updateRequest := &secretmanagerpb.UpdateSecretRequest{
		Secret: &secretmanagerpb.Secret{
			Name:        secretName,
			Annotations: result.Annotations,
		},
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: []string{"annotations"},
		},
	}

	if _, err := client.UpdateSecret(ctx, updateRequest); err != nil {
		return fmt.Errorf("failed to update secret: %w", err)
	}
	fmt.Fprintf(w, "Deleted annotation %s from secret %s\n", annotationKey, secretName)
	return nil
}

// [END secretmanager_delete_regional_secret_annotation]
