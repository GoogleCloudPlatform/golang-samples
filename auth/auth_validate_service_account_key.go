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

package snippets

// [START auth_validate_service_account_key]
import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"golang.org/x/oauth2/google"
)

// validateServiceAccountKey validates that a JSON file is a service account key.
//
// This sample uses a type-specific loader to ensure the credentials are
// specifically for a service account, which helps prevent the accidental
// use of other credential types such as user credentials.
func validateServiceAccountKey(w io.Writer, keyPath string) error {
	keyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		fmt.Fprintf(w, "failed to read service account key file: %v", err)
		return fmt.Errorf("failed to read service account key file %q: %w", keyPath, err)
	}

	scope := "https://www.googleapis.com/auth/cloud-platform"

	// Use a type-specific credential loader to validate the service account key.
	// google.CredentialsFromJSONWithType returns an error if the 'type' field in the JSON
	// is missing or is not 'service_account'.
	// Note: This validates the format and type locally; it does not verify
	// the key's status with Google Cloud's authentication server.
	ctx := context.Background()
	creds, err := google.CredentialsFromJSONWithType(ctx, keyBytes, google.ServiceAccount, scope)
	if err != nil {
		fmt.Fprintf(w, "invalid service account key: %v", err)
		return fmt.Errorf("invalid service account key: %w", err)
	}

	var f struct {
		ClientEmail string `json:"client_email"`
	}
	if err := json.Unmarshal(keyBytes, &f); err != nil {
		fmt.Fprintf(w, "failed to parse service account email: %v", err)
		return fmt.Errorf("failed to parse service account email: %w", err)
	}

	fmt.Fprintf(w, "Successfully validated service account key for: %s\n", f.ClientEmail)

	_ = creds // You can use creds.TokenSource to get a TokenSource for authenticated requests.

	return nil
}

// [END auth_validate_service_account_key]
