// Copyright 2022 Google LLC
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

// [START spanner_postgresql_interleaved_table]

import (
	"context"
	"fmt"
	"io"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
	adminpb "cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
)

// pgInterleavedTable shows how to create an interleaved table on a Spanner
// PostgreSQL database. The Spanner PostgreSQL dialect extends the standard
// PostgreSQL dialect to allow the creation of interleaved tables.
func pgInterleavedTable(w io.Writer, db string) error {
	// db := "projects/my-project/instances/my-instance/databases/my-database"
	ctx := context.Background()

	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return err
	}
	defer adminClient.Close()

	// The Spanner PostgreSQL dialect extends the PostgreSQL dialect with certain Spanner
	// specific features, such as interleaved tables.
	// See https://cloud.google.com/spanner/docs/postgresql/data-definition-language#create_table
	// for the full CREATE TABLE syntax.
	req := &adminpb.UpdateDatabaseDdlRequest{
		Database: db,
		Statements: []string{
			`CREATE TABLE Singers (
				SingerId  bigint NOT NULL PRIMARY KEY,
				FirstName varchar(1024) NOT NULL,
				LastName  varchar(1024) NOT NULL
			)`,
			`CREATE TABLE Albums (
				SingerId bigint NOT NULL,
				AlbumId  bigint NOT NULL,
				Title    varchar(1024) NOT NULL,
				PRIMARY KEY (SingerId, AlbumId)
			) INTERLEAVE IN PARENT Singers ON DELETE CASCADE`},
	}
	op, err := adminClient.UpdateDatabaseDdl(ctx, req)
	if err != nil {
		return err
	}
	if err := op.Wait(ctx); err != nil {
		return err
	}
	fmt.Fprint(w, "Created interleaved table hierarchy using PostgreSQL dialect\n")

	return nil
}

// [END spanner_postgresql_interleaved_table]
