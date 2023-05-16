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

// package downscopedoverview contains Google Cloud auth snippets showing how to
// downscope credentials with Credential Access Boundaries.
// https://cloud.google.com/iam/docs/downscoping-short-lived-credentials
package downscopedoverview

// [START auth_downscoping_initialize_downscoped_cred]

import (
	"context"
	"fmt"

	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/google/downscope"
)

// initializeCredentials will generate a downscoped token using the provided Access Boundary Rules.
func initializeCredentials(accessBoundary []downscope.AccessBoundaryRule) error {
	ctx := context.Background()

	// You must provide the "https://www.googleapis.com/auth/cloud-platform" scope.
	// This Source can be initialized in multiple ways; the following example uses
	// Application Default Credentials.
	rootSource, err := google.DefaultTokenSource(ctx, "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		return fmt.Errorf("failed to generate rootSource: %w", err)
	}

	// downscope.NewTokenSource constructs the token source with the configuration provided.
	dts, err := downscope.NewTokenSource(ctx, downscope.DownscopingConfig{RootSource: rootSource, Rules: accessBoundary})
	if err != nil {
		return fmt.Errorf("failed to generate downscoped token source: %w", err)
	}
	_ = dts
	// You can now use dts to access Google Storage resources.
	return nil
}

// [END auth_downscoping_initialize_downscoped_cred]
