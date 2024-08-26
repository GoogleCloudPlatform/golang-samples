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

// [START data_catalog_ptm_get_iam_policy]
import (
	"context"
	"fmt"
	"io"
	"strings"

	datacatalog "cloud.google.com/go/datacatalog/apiv1beta1"
	"cloud.google.com/go/iam/apiv1/iampb"
)

// getIAMPolicy prints information about the IAM policy associated with a
// given taxonomy or policy tag resource.
func getIAMPolicy(w io.Writer, resourceID string) error {
	// resourceID := "projects/myproject/locations/us/taxonomies/1234"
	ctx := context.Background()
	policyClient, err := datacatalog.NewPolicyTagManagerClient(ctx)
	if err != nil {
		return fmt.Errorf("datacatalog.NewPolicyTagManagerClient: %w", err)
	}
	defer policyClient.Close()

	req := &iampb.GetIamPolicyRequest{
		Resource: resourceID,
		Options: &iampb.GetPolicyOptions{
			RequestedPolicyVersion: 3,
		},
	}
	policy, err := policyClient.GetIamPolicy(ctx, req)
	if err != nil {
		return fmt.Errorf("GetIamPolicy: %w", err)
	}
	fmt.Fprintf(w, "Policy has version %d with Etag %x and %d bindings\n", policy.Version, policy.Etag, len(policy.Bindings))
	if len(policy.Bindings) > 0 {
		for _, binding := range policy.Bindings {
			fmt.Fprintf(w, "\trole %s with %d members: (%s)\n", binding.Role, len(binding.Members), strings.Join(binding.Members, ", "))
		}
	}
	return nil
}

// [END data_catalog_ptm_get_iam_policy]
