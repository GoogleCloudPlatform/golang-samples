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

// [START secretmanager_bind_tags_to_secret]

import (
	"context"
	"fmt"
	"io"

	resourcemanager "cloud.google.com/go/resourcemanager/apiv3"
	"cloud.google.com/go/resourcemanager/apiv3/resourcemanagerpb"
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

// bindTagsToSecret creates a secret and binds a tag to it.
func bindTagsToSecret(w io.Writer, projectID, secretID, tagValue string) error {
	// projectID := "my-project"
	// secretID := "my-secret"
	// tagValue := "tagValues/281476592621530"

	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create secretmanager client: %w", err)
	}
	defer client.Close()

	parent := fmt.Sprintf("projects/%s", projectID)

	createReq := &secretmanagerpb.CreateSecretRequest{
		Parent:   parent,
		SecretId: secretID,
		Secret: &secretmanagerpb.Secret{
			Replication: &secretmanagerpb.Replication{
				Replication: &secretmanagerpb.Replication_Automatic_{
					Automatic: &secretmanagerpb.Replication_Automatic{},
				},
			},
		},
	}

	secret, err := client.CreateSecret(ctx, createReq)
	if err != nil {
		return fmt.Errorf("failed to create secret: %w", err)
	}
	fmt.Fprintf(w, "Created secret %s\n", secret.Name)

	tagBindingsClient, err := resourcemanager.NewTagBindingsClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create tagbindings client: %w", err)
	}
	defer tagBindingsClient.Close()

	bindingReq := &resourcemanagerpb.CreateTagBindingRequest{
		TagBinding: &resourcemanagerpb.TagBinding{
			Parent:   fmt.Sprintf("//secretmanager.googleapis.com/%s", secret.Name),
			TagValue: tagValue,
		},
	}

	_, err = tagBindingsClient.CreateTagBinding(ctx, bindingReq)
	if err != nil {
		return fmt.Errorf("failed to start create tag binding operation: %w", err)
	}

	fmt.Fprintf(w, "Tag binding created for secret %s with tag value %s\n", secret.Name, tagValue)
	return nil
}

// [END secretmanager_bind_tags_to_secret]
