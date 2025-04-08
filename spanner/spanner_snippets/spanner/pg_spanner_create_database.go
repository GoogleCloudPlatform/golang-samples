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

// [START spanner_postgresql_create_database]
import (
	"context"
	"fmt"
	"io"
	"regexp"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
	adminpb "cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
)

// pgCreateDatabase shows how to create a Spanner database that uses the
// PostgreSQL dialect.
func pgCreateDatabase(ctx context.Context, w io.Writer, db string) error {
	// db := "projects/my-project/instances/my-instance/databases/my-database"
	matches := regexp.MustCompile("^(.*)/databases/(.*)$").FindStringSubmatch(db)
	if matches == nil || len(matches) != 3 {
		return fmt.Errorf("invalid database id %v", db)
	}

	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return err
	}
	defer adminClient.Close()

	// Databases with PostgreSQL dialect do not support extra DDL statements in the `CreateDatabase` call.
	req := &adminpb.CreateDatabaseRequest{
		Parent:          matches[1],
		DatabaseDialect: adminpb.DatabaseDialect_POSTGRESQL,
		// Note that PostgreSQL uses double quotes for quoting identifiers. This also
		// includes database names in the CREATE DATABASE statement.
		CreateStatement: `CREATE DATABASE "` + matches[2] + `"`,
	}
	opCreate, err := adminClient.CreateDatabase(ctx, req)
	if err != nil {
		return err
	}
	if _, err := opCreate.Wait(ctx); err != nil {
		return err
	}
	updateReq := &adminpb.UpdateDatabaseDdlRequest{
		Database: db,
		Statements: []string{
			`CREATE TABLE Singers (
				SingerId   bigint NOT NULL PRIMARY KEY,
				FirstName  varchar(1024),
				LastName   varchar(1024),
				SingerInfo bytea
			)`,
			`CREATE TABLE Albums (
				AlbumId      bigint NOT NULL,
				SingerId     bigint NOT NULL REFERENCES Singers (SingerId),
				AlbumTitle   text,
                PRIMARY KEY(SingerId, AlbumId)
			)`,
			`CREATE TABLE Venues (
				VenueId  bigint NOT NULL PRIMARY KEY,
				Name     varchar(1024) NOT NULL
			)`,
		},
	}
	opUpdate, err := adminClient.UpdateDatabaseDdl(ctx, updateReq)
	if err != nil {
		return err
	}
	if err := opUpdate.Wait(ctx); err != nil {
		return err
	}
	fmt.Fprintf(w, "Created Spanner PostgreSQL database [%v]\n", db)
	return nil
}

// [END spanner_postgresql_create_database]
