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

// [START storage_deactivate_hmac_key]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
)

// deactivateHMACKey deactivates the HMAC key with the given access ID.
func deactivateHMACKey(w io.Writer, accessID string, projectID string) (*storage.HMACKey, error) {
	ctx := context.Background()

	// Initialize client.
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close() // Closing the client safely cleans up background resources.

	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	handle := client.HMACKeyHandle(projectID, accessID)
	key, err := handle.Update(ctx, storage.HMACKeyAttrsToUpdate{State: "INACTIVE"})
	if err != nil {
		return nil, fmt.Errorf("Update: %w", err)
	}

	fmt.Fprintln(w, "The HMAC key metadata is:")
	fmt.Fprintf(w, "%+v", key)

	return key, nil
}

// [END storage_deactivate_hmac_key]
