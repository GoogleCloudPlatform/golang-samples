// Copyright 2021 Google LLC
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

package downscopingoverview

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/google/downscope"
)

// [START auth_overview_downscoping_token_broker]

// createDownscopedToken would be run on the token broker in order to generate
// a downscoped access token.  The token broker would then pass the newly created
// token to the requesting token consumer for use.
// TO REVIEWER: should I change the signature of this function?  Passing in
// arguments directly might make it a bit confusing, since I was asked to
// include the initialization of the boundary and the creation of the token in the same function.
func createDownscopedToken(ctx context.Context) (*oauth2.Token, error) {
	// A condition can optionally be provided to further restrict access permissions.
	condition := downscope.AvailabilityCondition{
		Expression:  "resource.name.startsWith('projects/_/buckets/foo/objects/profile-picture-'",
		Title:       "Pictures Only",
		Description: "Restricts a token to only be able to access objects that start with `profile-picture-`",
	}
	// Initializes an accessBoundary with one Rule which restricts the downscoped
	// token to only be able to access the bucket "foo" and only grants it the
	// permission "storage.objectViewer".
	accessBoundary := []downscope.AccessBoundaryRule{
		{
			AvailableResource:    "//storage.googleapis.com/projects/_/buckets/foo",
			AvailablePermissions: []string{"inRole:roles/storage.objectViewer"},
			Condition:            &condition, // Optional
		},
	}

	// This Source can be initialized in multiple ways; the following example uses
	// Application Default Credentials.
	var rootSource oauth2.TokenSource

	// You must provide the "https://www.googleapis.com/auth/cloud-platform" scope.
	rootSource, err := google.DefaultTokenSource(ctx, "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		return nil, fmt.Errorf("failed to generate rootSource: %v", err)
	}

	// downscope.NewTokenSource constructs the token source with the configuration provided.
	dts, err := downscope.NewTokenSource(ctx, downscope.DownscopingConfig{RootSource: rootSource, Rules: accessBoundary})
	if err != nil {
		return nil, fmt.Errorf("failed to generate downscoped token source: %v", err)
	}
	// Token() uses the previously declared TokenSource to generate a downscoped token.
	tok, err := dts.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %v", err)
	}
	// Pass this token back to the token consumer.
	return tok, nil
}

// [END downscoping_overview_token_broker]
