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
	iampb "google.golang.org/genproto/googleapis/iam/v1"
)

// setIAMPolicy defines the IAM policy for an associated taxonomy or policy tag.
func setIamPolicy(resourceID string, w io.Writer) error {
	ctx := context.Background()
	policyClient, err := datacatalog.NewPolicyTagManagerClient(ctx)
	if err != nil {
		return fmt.Errorf("datacatalog.NewPolicyTagManagerClient: %v", err)
	}
	defer policyClient.Close()

	req := &iampb.SetIamPolicyRequest{
		Resource: resourceID,
		Policy: &iampb.Policy{
			Version: 3,
			Bindings: []*iampb.Binding{
				{
					Role:    "roles/datacatalog.categoryFineGrainedReader",
					Members: []string{"allAuthenticatedUsers"},
				},
			},
		},
	}
	policy, err := policyClient.SetIamPolicy(ctx, req)
	if err != nil {
		return fmt.Errorf("SetIamPolicy: %v", err)
	}
	fmt.Fprintf(w, "set policy on resource %s with Etag %x\n", resourceID, policy.Etag)
	return nil
}
