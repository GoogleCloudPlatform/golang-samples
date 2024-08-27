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

// [START spanner_postgresql_case_sensitivity]

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/spanner"
	database "cloud.google.com/go/spanner/admin/database/apiv1"
	adminpb "cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
	"google.golang.org/api/iterator"
)

// pgCaseSensitivity shows the rules for case-sensitivity and case folding for
// a Spanner PostgreSQL database.
func pgCaseSensitivity(w io.Writer, db string) error {
	// db := "projects/my-project/instances/my-instance/databases/my-database"
	ctx := context.Background()

	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return err
	}
	defer adminClient.Close()

	// Spanner PostgreSQL follows the case sensitivity rules of PostgreSQL. This means that:
	// 1. Identifiers that are not double-quoted are folded to lower case.
	// 2. Identifiers that are double-quoted retain their case and are case-sensitive.
	// See https://www.postgresql.org/docs/current/sql-syntax-lexical.html#SQL-SYNTAX-IDENTIFIERS
	// for more information.
	req := &adminpb.UpdateDatabaseDdlRequest{
		Database: db,
		Statements: []string{
			`CREATE TABLE Singers (
				-- SingerId will be folded to "singerid"
				SingerId  bigint NOT NULL PRIMARY KEY,
				-- FirstName and LastName are double-quoted and will therefore retain their
				-- mixed case and are case-sensitive. This means that any statement that
				-- references any of these columns must use double quotes.
				"FirstName" varchar(1024) NOT NULL,
				"LastName"  varchar(1024) NOT NULL
			)`},
	}
	op, err := adminClient.UpdateDatabaseDdl(ctx, req)
	if err != nil {
		return err
	}
	if err := op.Wait(ctx); err != nil {
		return err
	}

	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	m := []*spanner.Mutation{
		spanner.InsertOrUpdateMap("Singers", map[string]interface{}{
			// Column names in mutations are always case-insensitive, regardless whether the
			// columns were double-quoted or not during creation.
			"singerid":  1,
			"firstname": "Bruce",
			"lastname":  "Allison",
		}),
	}
	_, err = client.Apply(context.Background(), m)
	if err != nil {
		return err
	}

	iter := client.Single().Query(ctx, spanner.Statement{SQL: "SELECT * FROM Singers"})
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		var singerID int64
		var firstName, lastName string
		// SingerId is automatically folded to lower case. Accessing the column by its name in
		// a result set must therefore use all lower-case letters.
		if err := row.ColumnByName("singerid", &singerID); err != nil {
			return err
		}
		// FirstName and LastName were double-quoted during creation, and retain their mixed
		// case when returned in a result set.
		if err := row.ColumnByName("FirstName", &firstName); err != nil {
			return err
		}
		if err := row.ColumnByName("LastName", &lastName); err != nil {
			return err
		}
		fmt.Fprintf(w, "SingerId: %d, FirstName: %s, LastName: %s\n", singerID, firstName, lastName)
	}

	// Aliases are also identifiers, and specifying an alias in double quotes will make the alias
	// retain its case.
	iterWithAliases := client.Single().Query(ctx, spanner.Statement{
		SQL: `SELECT singerid AS "SingerId",
				     concat("FirstName", ' '::varchar, "LastName") AS "FullName"
			  FROM Singers`})
	defer iterWithAliases.Stop()
	for {
		row, err := iterWithAliases.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		var singerID int64
		var fullName string
		// The aliases are double-quoted and therefore retains their mixed case.
		if err := row.ColumnByName("SingerId", &singerID); err != nil {
			return err
		}
		if err := row.ColumnByName("FullName", &fullName); err != nil {
			return err
		}
		fmt.Fprintf(w, "SingerId: %d, FullName: %s\n", singerID, fullName)
	}

	// DML statements must also follow the PostgreSQL case rules.
	stmt := spanner.Statement{
		SQL: `INSERT INTO Singers (SingerId, "FirstName", "LastName")
				  VALUES ($1, $2, $3)`,
		Params: map[string]interface{}{
			"p1": 2,
			"p2": "Alice",
			"p3": "Bruxelles",
		},
	}
	_, err = client.ReadWriteTransaction(ctx, func(ctx context.Context, transaction *spanner.ReadWriteTransaction) error {
		_, err := transaction.Update(ctx, stmt)
		return err
	})
	if err != nil {
		return err
	}

	return nil
}

// [END spanner_postgresql_case_sensitivity]
