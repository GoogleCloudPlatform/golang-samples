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

// Package authsnippets contains Google Cloud authentication snippets.
package authsnippets

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/pubsub/v1"
	"log"
	"os"
)

// [START auth_overview_service_account]

// Use service account to authenticate.
func service_account() {
	// Download service account key per https://cloud.google.com/docs/authentication/production.
	// Set environment variable GOOGLE_APPLICATION_CREDENTIALS=/path/to/service-account-key.json
	service, err := pubsub.NewService(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	// Use the authenticated client
	_ = service
}

// [END auth_overview_service_account]

// [START auth_overview_env_service_account]

// Use environment-provided service account to authenticate.
func env_service_account() {
	// If your application runs in a GCP environment, such as Compute Engine,
	// you don't need to provide any application credentials. The client
	// library will find the credentials by itself.
	service, err := pubsub.NewService(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	// Use the authenticated client
	_ = service
}

// [END auth_overview_env_service_account]

// [START auth_overview_api_key]

// Use API key to authenticate.
func api_key() {
	service, err := pubsub.NewService(context.Background(),
		option.WithAPIKey("api-key-string"))
	if err != nil {
		log.Fatal(err)
	}
	// Use the authenticated client
	_ = service
}

// [END auth_overview_api_key]

// [START auth_overview_oauth_client]

// Use OAuth client ID to authenticate.
func oauth_client() {
	ctx := context.Background()

	// Please make sure the redirect URL is the same as the one you specified when you
	// created the client ID.
	redirectURL := os.Getenv("OAUTH2_CALLBACK")
	if redirectURL == "" {
		redirectURL = "your redirect url"
	}
	config := &oauth2.Config{
		ClientID:     "your-client-id",
		ClientSecret: "your-client-secret",
		RedirectURL:  redirectURL,
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}

	// Dummy authorization flow to read auth code from stdin
	authorizeUrl := config.AuthCodeURL("your state")
	fmt.Printf("Follow the link in your browser to obtain auth code: %s", authorizeUrl)

	// Read the authentication code from the command line
	var code string
	fmt.Scanln(&code)

	// Exchange auth code for OAuth token.
	token, err := config.Exchange(ctx, code)
	if err != nil {
		log.Fatal(err)
	}
	service, err := pubsub.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, token)))

	// Use the authenticated client.
	_ = service
}

// [END auth_overview_oauth_client]
