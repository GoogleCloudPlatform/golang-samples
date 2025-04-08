// Copyright 2020 Google LLC
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

package policytagmanager

// [START data_catalog_ptm_create_policytag]
import (
	"context"
	"fmt"
	"io"

	datacatalog "cloud.google.com/go/datacatalog/apiv1beta1"
	"cloud.google.com/go/datacatalog/apiv1beta1/datacatalogpb"
)

// createPolicyTag creates a policy tag resource under a given parent taxonomy.
//
// It optionally accepts a parent ID, which can be used to create a hierarchical
// relationship between tags.
func createPolicyTag(w io.Writer, parent, displayName, parentPolicyTag string) (string, error) {
	// parent := "projects/myproject/locations/us/taxonomies/1234"
	// displayName := "Example Policy Tag"
	// parentPolicyTag := "projects/myproject/locations/us/taxonomies/1234/policyTags/5678"
	ctx := context.Background()
	policyClient, err := datacatalog.NewPolicyTagManagerClient(ctx)
	if err != nil {
		return "", fmt.Errorf("datacatalog.NewPolicyTagManagerClient: %w", err)
	}
	defer policyClient.Close()

	req := &datacatalogpb.CreatePolicyTagRequest{
		Parent: parent,
		PolicyTag: &datacatalogpb.PolicyTag{
			DisplayName: displayName,
			Description: "Example description for the policy tag",
		},
	}
	if parentPolicyTag != "" {
		req.PolicyTag.ParentPolicyTag = parentPolicyTag
	}
	resp, err := policyClient.CreatePolicyTag(ctx, req)
	if err != nil {
		return "", fmt.Errorf("CreatePolicyTag: %w", err)
	}

	fmt.Fprintf(w, "PolicyTag %s was created.\n", resp.Name)
	return resp.Name, nil
}

// [END data_catalog_ptm_create_policytag]
