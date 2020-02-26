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

// [START storage_create_hmac_key]
import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"io"
	"time"
)

// createHMACKey creates a new HMAC key using the given project and service account.
func createHMACKey(w io.Writer, projectID string, serviceAccountEmail string) (*storage.HMACKey, error) {
	ctx := context.Background()

	// Initialize client.
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close() // Closing the client safely cleans up background resources.

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	key, err := client.CreateHMACKey(ctx, projectID, serviceAccountEmail)
	if err != nil {
		return nil, fmt.Errorf("CreateHMACKey: %v", err)
	}

	fmt.Fprintf(w, "%s\n", key)
	fmt.Fprintf(w, "The base64 encoded secret is %s\n", key.Secret)
	fmt.Fprintln(w, "Do not miss that secret, there is no API to recover it.")
	fmt.Fprintln(w, "The HMAC key metadata is")
	fmt.Fprintf(w, "%+v", key)

	return key, nil
}

// [END storage_create_hmac_key]
