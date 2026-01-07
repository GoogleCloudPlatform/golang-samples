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
	"google.golang.org/protobuf/types/known/timestamppb"
)

// [START secretmanager_create_secret_with_expire_time]

// createSecretWithExpireTime creates a new secret with an expiration time.
func createSecretWithExpireTime(w io.Writer, projectID, secretID string, expireTime time.Time) error {
	// projectID := "my-project"
	// secretID := "my-secret-with-expiry"
	// expireTime := time.Now().Add(time.Hour * 24)

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
			Expiration: &secretmanagerpb.Secret_ExpireTime{
				ExpireTime: timestamppb.New(expireTime),
			},
		},
	}

	secret, err := client.CreateSecret(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create secret: %w", err)
	}

	fmt.Fprintf(w, "Created secret %s with expiration time %v\n", secret.Name, expireTime)
	return nil
}

// [END secretmanager_create_secret_with_expire_time]
