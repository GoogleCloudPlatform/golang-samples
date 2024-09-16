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

package overview

// [START auth_overview_oauth_client]
import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/auth"
	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

// oauthClient shows how to use an OAuth client ID to authenticate as an end-user.
func oauthClient() error {
	ctx := context.Background()

	// Please make sure the redirect URL is the same as the one you specified when you
	// created the client ID.
	redirectURL := os.Getenv("OAUTH2_CALLBACK")
	if redirectURL == "" {
		redirectURL = "your redirect url"
	}
	state := "state"
	opts := &auth.Options3LO{
		ClientID:     "your-client-id",
		ClientSecret: "your-client-secret",
		RedirectURL:  redirectURL,
		Scopes:       []string{"email", "profile"},
		AuthURL:      "https://accounts.google.com/o/oauth2/auth",
		TokenURL:     "https://oauth2.googleapis.com/token",
		AuthStyle:    auth.StyleInParams,
		AuthHandlerOpts: &auth.AuthorizationHandlerOptions{
			State: state,
			Handler: func(authCodeURL string) (string, string, error) {
				// Dummy authorization flow to read auth code from stdin.
				fmt.Printf("Follow the link in your browser to obtain auth code: %s", authCodeURL)
				var code string
				fmt.Scanln(&code)
				return code, state, nil
			},
		},
	}

	// Exchange auth code for OAuth token.
	tp, err := auth.New3LOTokenProvider(opts)
	if err != nil {
		return fmt.Errorf("auth.New3LOTokenProvider: %w", err)
	}
	creds := auth.NewCredentials(&auth.CredentialsOptions{
		TokenProvider: tp,
	})
	client, err := pubsub.NewClient(ctx, "your-project-id", option.WithAuthCredentials(creds))
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %w", err)
	}
	defer client.Close()

	// Use the authenticated client.
	_ = client

	return nil
}

// [END auth_overview_oauth_client]
