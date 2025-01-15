// Copyright 2023 Google LLC
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

package spanner

// [START spanner_enable_fine_grained_access]

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/iam/apiv1/iampb"
	database "cloud.google.com/go/spanner/admin/database/apiv1"
	expr "google.golang.org/genproto/googleapis/type/expr"
)

func enableFineGrainedAccess(w io.Writer, db string, iamMember string, databaseRole string, title string) error {
	// iamMember = "user:alice@example.com"
	// databaseRole = "parent"
	// title = "condition title"
	ctx := context.Background()
	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return err
	}
	defer adminClient.Close()

	policy, err := adminClient.GetIamPolicy(ctx, &iampb.GetIamPolicyRequest{
		Resource: db,
		Options: &iampb.GetPolicyOptions{
			// IAM conditions need at least version 3
			RequestedPolicyVersion: 3,
		},
	})
	if err != nil {
		return err
	}

	// IAM conditions need at least version 3
	if policy.Version < 3 {
		policy.Version = 3
	}
	policy.Bindings = append(policy.Bindings, []*iampb.Binding{
		{
			Role:    "roles/spanner.fineGrainedAccessUser",
			Members: []string{iamMember},
		},
		{
			Role:    "roles/spanner.databaseRoleUser",
			Members: []string{iamMember},
			Condition: &expr.Expr{
				Expression: fmt.Sprintf(`resource.name.endsWith("/databaseRoles/%s")`, databaseRole),
				Title:      title,
			},
		},
	}...)
	_, err = adminClient.SetIamPolicy(ctx, &iampb.SetIamPolicyRequest{
		Resource: db,
		Policy:   policy,
	})
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "Enabled fine-grained access in IAM.\n")
	return nil
}

// [END spanner_enable_fine_grained_access]
