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
	"google.golang.org/protobuf/types/known/timestamppb"
)

// [START secretmanager_create_secret_with_rotation]

// createSecretWithRotation creates a new secret with rotation configured.
func createSecretWithRotation(w io.Writer, projectID, secretID, topicName string) error {
	// projectID := "my-project"
	// secretID := "my-secret-with-rotation"
	// topicName := "projects/my-project/topics/my-topic"
	rotationPeriod := 24 * time.Hour

	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create secretmanager client: %w", err)
	}
	defer client.Close()

	req := &secretmanagerpb.CreateSecretRequest{
		Parent:   fmt.Sprintf("projects/%s", projectID),
		SecretId: secretID,
		Secret: &secretmanagerpb.Secret{
			Replication: &secretmanagerpb.Replication{
				Replication: &secretmanagerpb.Replication_Automatic_{
					Automatic: &secretmanagerpb.Replication_Automatic{},
				},
			},
			Topics: []*secretmanagerpb.Topic{
				{
					Name: topicName,
				},
			},
			Rotation: &secretmanagerpb.Rotation{
				NextRotationTime: timestamppb.New(time.Now().Add(time.Hour * 24)),
				RotationPeriod:   durationpb.New(rotationPeriod),
			},
		},
	}

	secret, err := client.CreateSecret(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create secret: %w", err)
	}

	fmt.Fprintf(w, "Created secret %s with rotation period %v and topic %s\n", secret.Name, rotationPeriod, topicName)
	return nil
}

// [END secretmanager_create_secret_with_rotation]
