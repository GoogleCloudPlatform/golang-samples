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

// [START auth_cloud_idtoken_impersonated_credentials]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/auth/credentials"
	"cloud.google.com/go/auth/credentials/impersonate"
)

// getIdTokenFromImpersonatedCredentials uses a service account (SA1) to impersonate as
// another service account (SA2) and obtain id token for the impersonated account.
// To obtain token for SA2, SA1 should have the "roles/iam.serviceAccountTokenCreator" permission on SA2.
func getIdTokenFromImpersonatedCredentials(w io.Writer, scope, targetAudience, impersonatedServiceAccount string) error {
	// scope := "https://www.googleapis.com/auth/cloud-platform"
	// targetAudience := "http://www.example.com"
	// impersonatedServiceAccount := "name@project.service.gserviceaccount.com"
	ctx := context.Background()

	// Construct Credentials which obtains the default configuration from your
	// working environment.
	creds, err := credentials.DetectDefault(&credentials.DetectOptions{
		Scopes: []string{scope},
	})
	if err != nil {
		return fmt.Errorf("failed to generate default credentials: %w", err)
	}

	impCreds, err := impersonate.NewIDTokenCredentials(&impersonate.IDTokenOptions{
		Audience:        targetAudience,
		TargetPrincipal: impersonatedServiceAccount,
		IncludeEmail:    true,
		// delegates: The chained list of delegates required to grant the final accessToken.
		// For more information, see:
		// https://cloud.google.com/iam/docs/create-short-lived-credentials-direct#sa-credentials-permissions
		// Delegates is NOT USED here.
		Delegates:   []string{},
		Credentials: creds,
	})
	if err != nil {
		return fmt.Errorf("NewIDTokenCredentials error: %w", err)
	}

	// Get the ID token. Once you've obtained the ID token, you can use it to
	// make an authenticated call to the target audience.
	_, err = impCreds.Token(ctx)
	if err != nil {
		return fmt.Errorf("failed to receive token: %w", err)
	}
	fmt.Fprintf(w, "Generated ID token.\n")

	return nil
}

// [END auth_cloud_idtoken_impersonated_credentials]
