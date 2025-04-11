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

// [START spanner_add_and_drop_database_role]

import (
	"context"
	"io"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
	adminpb "cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
)

func addAndDropDatabaseRole(w io.Writer, db string) error {
	ctx := context.Background()
	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return err
	}
	defer adminClient.Close()

	// Set up database roles and membership. After database roles are created,
	// users can be granted roles by setting IAM policies.
	op, err := adminClient.UpdateDatabaseDdl(ctx, &adminpb.UpdateDatabaseDdlRequest{
		Database: db,
		Statements: []string{
			"CREATE ROLE parent",
			"GRANT SELECT ON TABLE Albums TO ROLE parent",
			"CREATE ROLE child",
			"GRANT ROLE parent TO ROLE child",
		},
	})
	if err != nil {
		return err
	}
	if err := op.Wait(ctx); err != nil {
		return err
	}

	// Delete role and membership.
	op, err = adminClient.UpdateDatabaseDdl(ctx, &adminpb.UpdateDatabaseDdlRequest{
		Database: db,
		Statements: []string{
			"REVOKE ROLE parent FROM ROLE child",
			"DROP ROLE child",
		},
	})
	if err != nil {
		return err
	}
	if err := op.Wait(ctx); err != nil {
		return err
	}
	return nil
}

// [END spanner_add_and_drop_database_role]
