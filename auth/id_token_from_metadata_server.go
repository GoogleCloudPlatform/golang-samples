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

// [START auth_cloud_idtoken_metadata_server]
import (
	"context"
	"fmt"
	"io"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"
	"google.golang.org/api/option"
)

// getIdTokenFromMetadataServer uses the Google Cloud metadata server environment
// to create an identity token and add it to the HTTP request as part of an Authorization header.
func getIdTokenFromMetadataServer(w io.Writer, url string) error {
	// url := "http://www.example.com"

	ctx := context.Background()

	// Construct the GoogleCredentials object which obtains the default configuration from your
	// working environment.
	credentials, err := google.FindDefaultCredentials(ctx)
	if err != nil {
		return fmt.Errorf("failed to generate default credentials: %w", err)
	}

	ts, err := idtoken.NewTokenSource(ctx, url, option.WithCredentials(credentials))
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
	fmt.Fprintf(w, "Generated ID token.\n")

	return nil
}

// [END auth_cloud_idtoken_metadata_server]
