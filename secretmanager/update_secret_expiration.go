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

// [START secretmanager_update_secret_expiration]

import (
	"context"
	"fmt"
	"io"
	"time"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// updateSecretExpiration updates the expiration time of a secret.
func updateSecretExpiration(w io.Writer, secretName string) error {
	// secretName := "projects/my-project/secrets/my-secret"
	newExpire := time.Now().Add(2 * time.Hour)

	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create secretmanager client: %w", err)
	}
	defer client.Close()

	req := &secretmanagerpb.UpdateSecretRequest{
		Secret: &secretmanagerpb.Secret{
			Name: secretName,
			Expiration: &secretmanagerpb.Secret_ExpireTime{
				ExpireTime: timestamppb.New(newExpire),
			},
		},
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: []string{"expire_time"},
		},
	}

	secret, err := client.UpdateSecret(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to update secret: %w", err)
	}

	fmt.Fprintf(w, "Updated secret %s expiration time to %v\n", secret.Name, newExpire)
	return nil
}

// [END secretmanager_update_secret_expiration]
