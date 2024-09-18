// Copyright 2023 Google LLC
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

// [START auth_cloud_accesstoken_impersonated_credentials]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/auth/credentials"
	"cloud.google.com/go/auth/credentials/impersonate"
)

// getAccessTokenFromImpersonatedCredentials uses a service account (SA1) to impersonate
// another service account (SA2) and obtain a token for the impersonated account.
// To obtain a token for SA2, SA1 should have the "roles/iam.serviceAccountTokenCreator" permission on SA2.
func getAccessTokenFromImpersonatedCredentials(w io.Writer, impersonatedServiceAccount, scope string) error {
	// impersonatedServiceAccount := "name@project.service.gserviceaccount.com"
	// scope := "https://www.googleapis.com/auth/cloud-platform"

	ctx := context.Background()

	// Construct the Credentials object which obtains the default configuration
	// from your working environment.
	creds, err := credentials.DetectDefault(&credentials.DetectOptions{
		Scopes: []string{scope},
	})
	if err != nil {
		fmt.Fprintf(w, "failed to generate default credentials: %v", err)
		return fmt.Errorf("failed to generate default credentials: %w", err)
	}

	impCreds, err := impersonate.NewCredentials(&impersonate.CredentialsOptions{
		TargetPrincipal: impersonatedServiceAccount,
		Scopes:          []string{scope},
		Lifetime:        300 * time.Second,
		// delegates: The chained list of delegates required to grant the final accessToken.
		// For more information, see:
		// https://cloud.google.com/iam/docs/create-short-lived-credentials-direct#sa-credentials-permissions
		// Delegates is NOT USED here.
		Delegates:   []string{},
		Credentials: creds,
	})
	if err != nil {
		fmt.Fprintf(w, "NewCredentials error: %v", err)
		return fmt.Errorf("NewCredentials error: %w", err)
	}

	// Get the token. Once you've obtained the token, you can use it to make a
	// authenticated call.
	t, err := impCreds.Token(ctx)
	if err != nil {
		fmt.Fprintf(w, "failed to receive token: %v", err)
		return fmt.Errorf("failed to receive token: %w", err)
	}
	fmt.Fprintf(w, "Generated OAuth2 token with length %d.\n", len(t.Value))

	return nil
}

// [END auth_cloud_accesstoken_impersonated_credentials]
