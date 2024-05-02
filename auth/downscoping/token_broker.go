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

package downscopedoverview

// [START auth_downscoping_token_broker]

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/google/downscope"
)

// createDownscopedToken would be run on the token broker in order to generate
// a downscoped access token that only grants access to objects whose name begins with prefix.
// The token broker would then pass the newly created token to the requesting token consumer for use.
func createDownscopedToken(bucketName string, prefix string) error {
	// bucketName := "foo"
	// prefix := "profile-picture-"

	ctx := context.Background()
	// A condition can optionally be provided to further restrict access permissions.
	condition := downscope.AvailabilityCondition{
		Expression:  "resource.name.startsWith('projects/_/buckets/" + bucketName + "/objects/" + prefix + "')",
		Title:       prefix + " Only",
		Description: "Restricts a token to only be able to access objects that start with `" + prefix + "`",
	}
	// Initializes an accessBoundary with one Rule which restricts the downscoped
	// token to only be able to access the bucket "bucketName" and only grants it the
	// permission "storage.objectViewer".
	accessBoundary := []downscope.AccessBoundaryRule{
		{
			AvailableResource:    "//storage.googleapis.com/projects/_/buckets/" + bucketName,
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
		return fmt.Errorf("failed to generate rootSource: %w", err)
	}

	// downscope.NewTokenSource constructs the token source with the configuration provided.
	dts, err := downscope.NewTokenSource(ctx, downscope.DownscopingConfig{RootSource: rootSource, Rules: accessBoundary})
	if err != nil {
		return fmt.Errorf("failed to generate downscoped token source: %w", err)
	}
	// Token() uses the previously declared TokenSource to generate a downscoped token.
	tok, err := dts.Token()
	if err != nil {
		return fmt.Errorf("failed to generate token: %w", err)
	}
	// Pass this token back to the token consumer.
	_ = tok
	return nil
}

// [END auth_downscoping_token_broker]
