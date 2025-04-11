// Copyright 2022 Google LLC
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

// [START auth_cloud_idtoken_service_account]
import (
	"context"
	"fmt"
	"io"
	"os"

	"google.golang.org/api/idtoken"
	"google.golang.org/api/option"
)

// getIdTokenFromServiceAccount obtains the id token
// by providing the target audience using service account credentials.
func getIdTokenFromServiceAccount(w io.Writer, jsonCredentialsPath, url string) error {
	// jsonCredentialsPath := "/path/example"
	// targetAudience := "http://www.example.com"

	// *NOTE*:
	// Using service account keys introduces risk; they are long-lived, and can be used by anyone
	// that obtains the key. Proper rotation and storage reduce this risk but do not eliminate it.
	// For these reasons, you should consider an alternative approach that
	// does not use a service account key. Several alternatives to service account keys
	// are described here:
	// https://cloud.google.com/docs/authentication/external/set-up-adc

	ctx := context.Background()

	data, err := os.ReadFile(jsonCredentialsPath)
	if err != nil {
		return fmt.Errorf("failed to read json file: %w", err)
	}

	ts, err := idtoken.NewTokenSource(ctx, url, option.WithCredentialsJSON(data))
	if err != nil {
		return fmt.Errorf("failed to create NewTokenSource: %w", err)
	}

	// Get the ID token.
	// Once you've obtained the ID token, you can use it to make an authenticated call
	// to the target audience.
	_, err = ts.Token()
	if err != nil {
		return fmt.Errorf("failed to receive token: %w", err)
	}
	fmt.Fprintf(w, "Generated ID token. \n")

	return nil
}

// [END auth_cloud_idtoken_service_account]
