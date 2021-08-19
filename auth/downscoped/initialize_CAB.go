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

// [START auth_downscoping_rules]

import "golang.org/x/oauth2/google/downscope"

// constructCAB shows how to initialize a Credential Access Boundary for downscoping tokens.
// This sample uses on Access Boundary Rule; up to 10 may be included.
func constructCAB(bucketName string) {
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
			AvailableResource:    "//storage.googleapis.com/projects/_/buckets/" + bucketName,
			AvailablePermissions: []string{"inRole:roles/storage.objectViewer"},
			Condition:            &condition, // Optional
		},
	}

	// You can now use this accessBoundary to generate a downscoped token
	_ = accessBoundary
}

// [END auth_downscoping_rules]
