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

package spanner

// [START spanner_create_database]
import (
	"context"
	"fmt"
	"io"
	"regexp"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
	adminpb "cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
)

func createDatabase(ctx context.Context, w io.Writer, db string) error {
	matches := regexp.MustCompile("^(.*)/databases/(.*)$").FindStringSubmatch(db)
	if matches == nil || len(matches) != 3 {
		return fmt.Errorf("Invalid database id %s", db)
	}

	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return err
	}
	defer adminClient.Close()

	op, err := adminClient.CreateDatabase(ctx, &adminpb.CreateDatabaseRequest{
		Parent:          matches[1],
		CreateStatement: "CREATE DATABASE `" + matches[2] + "`",
		ExtraStatements: []string{
			`CREATE TABLE Singers (
				SingerId   INT64 NOT NULL,
				FirstName  STRING(1024),
				LastName   STRING(1024),
				SingerInfo BYTES(MAX),
				FullName   STRING(2048) AS (
					ARRAY_TO_STRING([FirstName, LastName], " ")
				) STORED
			) PRIMARY KEY (SingerId)`,
			`CREATE TABLE Albums (
				SingerId     INT64 NOT NULL,
				AlbumId      INT64 NOT NULL,
				AlbumTitle   STRING(MAX)
			) PRIMARY KEY (SingerId, AlbumId),
			INTERLEAVE IN PARENT Singers ON DELETE CASCADE`,
		},
	})
	if err != nil {
		return err
	}
	if _, err := op.Wait(ctx); err != nil {
		return err
	}
	fmt.Fprintf(w, "Created database [%s]\n", db)
	return nil
}

// [END spanner_create_database]
