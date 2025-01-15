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

// [START data_catalog_ptm_set_iam_policy]
import (
	"context"
	"fmt"
	"io"

	datacatalog "cloud.google.com/go/datacatalog/apiv1beta1"
	"cloud.google.com/go/iam/apiv1/iampb"
)

// setIAMPolicy demonstrates altering the policy of a given taxonomy or policy
// tag resource.  In this example, we append a binding to the existing policy
// to add the fine grained reader role to a specific member.
func setIAMPolicy(w io.Writer, resourceID, member string) error {
	// resourceID := "projects/myproject/locations/us/taxonomies/1234/policyTags/5678"
	// member := "group:my-trusted-group@example.com"
	ctx := context.Background()
	policyClient, err := datacatalog.NewPolicyTagManagerClient(ctx)
	if err != nil {
		return fmt.Errorf("datacatalog.NewPolicyTagManagerClient: %w", err)
	}
	defer policyClient.Close()

	// First, retrieve the existing policy.
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

	// Alter the policy to add an additional binding.
	newPolicy := policy
	newPolicy.Bindings = append(newPolicy.Bindings, &iampb.Binding{
		Role:    "roles/datacatalog.categoryFineGrainedReader",
		Members: []string{member},
	})

	sReq := &iampb.SetIamPolicyRequest{
		Resource: resourceID,
		Policy:   newPolicy,
	}
	updatedPolicy, err := policyClient.SetIamPolicy(ctx, sReq)
	if err != nil {
		return fmt.Errorf("SetIamPolicy: %w", err)
	}
	fmt.Fprintf(w, "set policy on resource %s with Etag %x\n", resourceID, updatedPolicy.Etag)
	return nil
}

// [END data_catalog_ptm_set_iam_policy]
