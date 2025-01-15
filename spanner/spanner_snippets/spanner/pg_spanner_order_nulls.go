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

// [START spanner_postgresql_order_nulls]

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/spanner"
	database "cloud.google.com/go/spanner/admin/database/apiv1"
	adminpb "cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
	"google.golang.org/api/iterator"
)

// pgOrderNulls shows how a Spanner PostgreSQL database orders null values in a
// query, and how an application can change the default behavior by adding
// `NULLS FIRST` or `NULLS LAST` to an `ORDER BY` clause.
func pgOrderNulls(w io.Writer, db string) error {
	// db := "projects/my-project/instances/my-instance/databases/my-database"
	ctx := context.Background()

	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return err
	}
	defer adminClient.Close()

	// Spanner PostgreSQL follows the ORDER BY rules for NULL values of PostgreSQL. This means that:
	// 1. NULL values are ordered last by default when a query result is ordered in ascending order.
	// 2. NULL values are ordered first by default when a query result is ordered in descending order.
	// 3. NULL values can be order first or last by specifying NULLS FIRST or NULLS LAST in the ORDER BY clause.
	op, err := adminClient.UpdateDatabaseDdl(ctx, &adminpb.UpdateDatabaseDdlRequest{
		Database: db,
		Statements: []string{
			`CREATE TABLE Singers (
				SingerId  bigint NOT NULL PRIMARY KEY,
				Name varchar(1024)
			)`},
	})
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

	if _, err := client.Apply(ctx, []*spanner.Mutation{
		spanner.InsertOrUpdateMap("Singers", map[string]interface{}{
			"SingerId": 1,
			"Name":     "Bruce",
		}),
		spanner.InsertOrUpdateMap("Singers", map[string]interface{}{
			"SingerId": 2,
			"Name":     "Alice",
		}),
		spanner.InsertOrUpdateMap("Singers", map[string]interface{}{
			"SingerId": 3,
			"Name":     spanner.NullString{},
		}),
	}); err != nil {
		return err
	}

	// This returns the singers in order Alice, Bruce, null
	iterOrderByName := client.Single().Query(ctx, spanner.Statement{SQL: "SELECT Name FROM Singers ORDER BY Name"})
	fmt.Fprintln(w, "Singers ORDER BY Name")
	if err := printSingerNames(w, iterOrderByName); err != nil {
		return err
	}

	// This returns the singers in order null, Bruce, Alice
	iterOrderByNameDesc := client.Single().Query(ctx, spanner.Statement{SQL: "SELECT Name FROM Singers ORDER BY Name DESC"})
	fmt.Fprintln(w, "Singers ORDER BY Name DESC")
	if err := printSingerNames(w, iterOrderByNameDesc); err != nil {
		return err
	}

	// This returns the singers in order null, Alice, Bruce
	iterOrderByNameNullsFirst := client.Single().Query(ctx, spanner.Statement{SQL: "SELECT Name FROM Singers ORDER BY Name NULLS FIRST"})
	fmt.Fprintln(w, "Singers ORDER BY Name NULLS FIRST")
	if err := printSingerNames(w, iterOrderByNameNullsFirst); err != nil {
		return err
	}

	// This returns the singers in order Bruce, Alice, null
	iterOrderByNameDescNullsLast := client.Single().Query(ctx, spanner.Statement{SQL: "SELECT Name FROM Singers ORDER BY Name DESC NULLS LAST"})
	fmt.Fprintln(w, "Singers ORDER BY Name DESC NULLS LAST")
	if err := printSingerNames(w, iterOrderByNameDescNullsLast); err != nil {
		return err
	}

	return nil
}

func printSingerNames(w io.Writer, iter *spanner.RowIterator) error {
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			return nil
		}
		if err != nil {
			return err
		}
		var name spanner.NullString
		if err := row.ColumnByName("name", &name); err != nil {
			return err
		}
		fmt.Fprintf(w, "\t%s\n", name)
	}
}

// [END spanner_postgresql_order_nulls]
