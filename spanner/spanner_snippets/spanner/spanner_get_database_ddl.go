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

package spanner

// [START spanner_get_database_ddl]
import (
	"context"
	"fmt"
	"io"
	"regexp"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
	adminpb "cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
)

// getDatabaseDdl gets the DDL for the database
func getDatabaseDdl(w io.Writer, db string) error {
	// db = `projects/<project>/instances/<instance-id>/database/<database-id>`
	matches := regexp.MustCompile("^(.*)/databases/(.*)$").FindStringSubmatch(db)
	if matches == nil || len(matches) != 3 {
		return fmt.Errorf("getDatabaseDdl: invalid database id %s", db)
	}

	ctx := context.Background()
	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return err
	}
	defer adminClient.Close()

	op, err := adminClient.GetDatabaseDdl(ctx, &adminpb.GetDatabaseDdlRequest{
		Database: db,
	})

	if err != nil {
		return err
	}

	fmt.Fprintf(w, "Database DDL is as follows: \n [%v]", op.GetStatements())

	return nil

}

// [END spanner_get_database_ddl]
