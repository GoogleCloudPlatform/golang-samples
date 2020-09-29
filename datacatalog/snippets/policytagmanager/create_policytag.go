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

import (
	"context"
	"fmt"
	"io"

	datacatalog "cloud.google.com/go/datacatalog/apiv1beta1"
	datacatalogpb "google.golang.org/genproto/googleapis/cloud/datacatalog/v1beta1"
)

// createPolicyTag creates a sample policy tag.
func createPolicyTag(parent, displayName, parentPolicyTag string, w io.Writer) (string, error) {
	// parent := "projects/myproject/locations/us/taxonomies/1234"
	// displayName := "PII Policy"
	ctx := context.Background()
	policyClient, err := datacatalog.NewPolicyTagManagerClient(ctx)
	if err != nil {
		return "", fmt.Errorf("datacatalog.NewPolicyTagManagerClient: %v", err)
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
		return "", fmt.Errorf("CreatePolicyTag: %v", err)
	}

	fmt.Fprintf(w, "PolicyTag %s was created.\n", resp.Name)
	return resp.Name, nil
}
