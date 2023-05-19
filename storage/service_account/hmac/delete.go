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

package hmac

// [START storage_delete_hmac_key]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
)

// deleteHMACKey deletes the HMAC key with the given access ID. Key must have state
// INACTIVE in order to succeed.
func deleteHMACKey(w io.Writer, accessID string, projectID string) error {
	ctx := context.Background()

	// Initialize client.
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close() // Closing the client safely cleans up background resources.

	handle := client.HMACKeyHandle(projectID, accessID)
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	if err = handle.Delete(ctx); err != nil {
		return fmt.Errorf("Delete: %w", err)
	}

	fmt.Fprintln(w, "The key is deleted, though it may still appear in ListHMACKeys results.")

	return nil
}

// [END storage_delete_hmac_key]
