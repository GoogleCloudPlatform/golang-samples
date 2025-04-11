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

// [START auth_cloud_verify_google_idtoken]
import (
	"context"
	"fmt"
	"io"

	"google.golang.org/api/idtoken"
)

// verifyGoogleIdToken verifies the obtained Google id token.
// This is done at the receiving end of the OIDC endpoint.
// The most common use case for verifying the ID token is when you are protecting
// your own APIs with IAP. Google services already verify credentials as a platform,
// so verifying ID tokens before making Google API calls is usually unnecessary.
func verifyGoogleIdToken(w io.Writer, token, expectedAudience string) error {
	// url := "id-token"
	// targetAudience := "https://example.com"

	ctx := context.Background()

	validator, err := idtoken.NewValidator(ctx)
	if err != nil {
		return fmt.Errorf("failed to create NewValidator: %w", err)
	}

	payload, err := validator.Validate(ctx, token, expectedAudience)
	if err != nil {
		return fmt.Errorf("failed to validate ID token: %w", err)
	}

	// Verify that the token contains subject.
	// Get the User id.
	if payload.Subject != "" {
		fmt.Fprintf(w, "User id: %s", payload.Subject)
	}
	fmt.Fprintf(w, "ID token verified.\n")

	return nil
}

// [END auth_cloud_verify_google_idtoken]
