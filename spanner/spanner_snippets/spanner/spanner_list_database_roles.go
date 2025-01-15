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

// [START spanner_list_database_roles]

import (
	"context"
	"fmt"
	"io"
	"strings"

	"google.golang.org/api/iterator"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
	adminpb "cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
)

func listDatabaseRoles(w io.Writer, db string) error {
	ctx := context.Background()
	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return err
	}
	defer adminClient.Close()

	iter := adminClient.ListDatabaseRoles(ctx, &adminpb.ListDatabaseRolesRequest{
		Parent: db,
	})
	rolePrefix := db + "/databaseRoles/"
	for {
		role, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		if !strings.HasPrefix(role.Name, rolePrefix) {
			return fmt.Errorf("Role %v does not have prefix %v", role.Name, rolePrefix)
		}
		fmt.Fprintf(w, "%s\n", strings.TrimPrefix(role.Name, rolePrefix))
	}
	return nil
}

// [END spanner_list_database_roles]
