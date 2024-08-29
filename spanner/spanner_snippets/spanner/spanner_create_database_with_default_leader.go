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

// [START spanner_create_database_with_default_leader]
import (
	"context"
	"fmt"
	"io"
	"regexp"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
	adminpb "cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
)

// createDatabaseWithDefaultLeader creates a database with a default leader
func createDatabaseWithDefaultLeader(w io.Writer, db string, defaultLeader string) error {
	// db = `projects/<project>/instances/<instance-id>/database/<database-id>`
	// defaultLeader = `my-default-leader`
	matches := regexp.MustCompile("^(.*)/databases/(.*)$").FindStringSubmatch(db)
	if matches == nil || len(matches) != 3 {
		return fmt.Errorf("createDatabaseWithDefaultLeader: invalid database id %s", db)
	}

	ctx := context.Background()
	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return err
	}
	defer adminClient.Close()

	alterDatabase := fmt.Sprintf(
		"ALTER DATABASE `%s` SET OPTIONS (default_leader = '%s')",
		matches[2], defaultLeader,
	)

	req := adminpb.CreateDatabaseRequest{
		Parent:          matches[1],
		CreateStatement: "CREATE DATABASE `" + matches[2] + "`",
		ExtraStatements: []string{
			`CREATE TABLE Singers (
				SingerId   INT64 NOT NULL,
				FirstName  STRING(1024),
				LastName   STRING(1024),
				SingerInfo BYTES(MAX)
			) PRIMARY KEY (SingerId)`,
			`CREATE TABLE Albums (
				SingerId     INT64 NOT NULL,
				AlbumId      INT64 NOT NULL,
				AlbumTitle   STRING(MAX)
			) PRIMARY KEY (SingerId, AlbumId),
			INTERLEAVE IN PARENT Singers ON DELETE CASCADE`,
			alterDatabase,
		},
	}
	op, err := adminClient.CreateDatabase(ctx, &req)
	if err != nil {
		return fmt.Errorf("createDatabaseWithDefaultLeader.CreateDatabase: %w", err)
	}
	dbObj, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("createDatabaseWithDefaultLeader.Wait: %w", err)
	}
	fmt.Fprintf(w, "Created database [%s] with default leader%q\n", dbObj.Name, dbObj.DefaultLeader)
	return nil

}

// [END spanner_create_database_with_default_leader]
