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

	"golang.org/x/oauth2/google"
	"google.golang.org/api/impersonate"
	"google.golang.org/api/option"
)

// getAccessTokenFromImpersonatedCredentials uses a service account (SA1) to impersonate
// another service account (SA2) and obtain OAuth2 token for the impersonated account.
// To obtain a token for SA2, SA1 should have the "roles/iam.serviceAccountTokenCreator" permission on SA2.
func getAccessTokenFromImpersonatedCredentials(w io.Writer, impersonatedServiceAccount, scope string) error {
	// impersonatedServiceAccount := "name@project.service.gserviceaccount.com"
	// scope := "https://www.googleapis.com/auth/cloud-platform"

	ctx := context.Background()

	// Construct the GoogleCredentials object which obtains the default configuration from your
	// working environment.
	credentials, err := google.FindDefaultCredentials(ctx, scope)
	if err != nil {
		fmt.Fprintf(w, "failed to generate default credentials: %v", err)
		return fmt.Errorf("failed to generate default credentials: %w", err)
	}

	ts, err := impersonate.CredentialsTokenSource(ctx, impersonate.CredentialsConfig{
		TargetPrincipal: impersonatedServiceAccount,
		Scopes:          []string{scope},
		Lifetime:        300 * time.Second,
		// delegates: The chained list of delegates required to grant the final accessToken.
		// For more information, see:
		// https://cloud.google.com/iam/docs/create-short-lived-credentials-direct#sa-credentials-permissions
		// Delegates is NOT USED here.
		Delegates: []string{},
	}, option.WithCredentials(credentials))
	if err != nil {
		fmt.Fprintf(w, "CredentialsTokenSource error: %v", err)
		return fmt.Errorf("CredentialsTokenSource error: %w", err)
	}

	// Get the OAuth2 token.
	// Once you've obtained the OAuth2 token, you can use it to make an authenticated call.
	t, err := ts.Token()
	if err != nil {
		fmt.Fprintf(w, "failed to receive token: %v", err)
		return fmt.Errorf("failed to receive token: %w", err)
	}
	fmt.Fprintf(w, "Generated OAuth2 token with length %d.\n", len(t.AccessToken))

	return nil
}

// [END auth_cloud_accesstoken_impersonated_credentials]
