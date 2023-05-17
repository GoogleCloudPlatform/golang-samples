// Copyright 2020 Google LLC
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

package serviceaccount

// [START storage_get_service_account]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
)

// getServiceAccount gets the default Cloud Storage service account email address.
func getServiceAccount(w io.Writer, projectID string) error {
	// projectID := "my-project-id"

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	serviceAccount, err := client.ServiceAccount(ctx, projectID)
	if err != nil {
		return fmt.Errorf("ServiceAccount: %w", err)
	}

	fmt.Fprintf(w, "The GCS service account for project %v is: %v\n", projectID, serviceAccount)
	return nil
}

// [END storage_get_service_account]
