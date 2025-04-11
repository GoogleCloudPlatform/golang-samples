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

// [START storage_list_hmac_keys]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

// listHMACKeys lists all HMAC keys associated with the project.
func listHMACKeys(w io.Writer, projectID string) ([]*storage.HMACKey, error) {
	ctx := context.Background()

	// Initialize client.
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close() // Closing the client safely cleans up background resources.

	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	iter := client.ListHMACKeys(ctx, projectID)
	var keys []*storage.HMACKey
	for {
		key, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("ListHMACKeys: %w", err)
		}
		fmt.Fprintf(w, "Service Account Email: %s\n", key.ServiceAccountEmail)
		fmt.Fprintf(w, "Access ID: %s\n", key.AccessID)

		keys = append(keys, key)
	}

	return keys, nil
}

// [END storage_list_hmac_keys]
