// Copyright 2019 Google LLC
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

package overview

// [START auth_overview_service_account]
import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
)

// serviceAccount shows how to use a service account to authenticate.
func serviceAccount() error {
	// Download service account key per https://cloud.google.com/docs/authentication/production.
	// Set environment variable GOOGLE_APPLICATION_CREDENTIALS=/path/to/service-account-key.json
	// This environment variable will be automatically picked up by the client.
	client, err := pubsub.NewClient(context.Background(), "your-project-id")
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %w", err)
	}
	defer client.Close()
	// Use the authenticated client.
	_ = client

	return nil
}

// [END auth_overview_service_account]
