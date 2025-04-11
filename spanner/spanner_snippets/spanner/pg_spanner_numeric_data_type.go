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

// [START spanner_postgresql_numeric_data_type]

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/spanner"
	database "cloud.google.com/go/spanner/admin/database/apiv1"
	adminpb "cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
	"google.golang.org/api/iterator"
)

// pgNumericDataType shows how to work with the PostgreSQL NUMERIC/DECIMAL data
// type on a Spanner PostgreSQL database.
func pgNumericDataType(w io.Writer, db string) error {
	// db := "projects/my-project/instances/my-instance/databases/my-database"
	ctx := context.Background()

	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return err
	}
	defer adminClient.Close()

	// Create a table that includes a column with data type NUMERIC. As the database has been
	// created with the PostgreSQL dialect, the data type that is used will be the PostgreSQL
	// NUMERIC / DECIMAL data type.
	op, err := adminClient.UpdateDatabaseDdl(ctx, &adminpb.UpdateDatabaseDdlRequest{
		Database: db,
		Statements: []string{
			`CREATE TABLE Venues (
				VenueId  bigint NOT NULL PRIMARY KEY,
				Name     varchar(1024) NOT NULL,
				Revenues numeric
			)`},
	})
	if err != nil {
		return err
	}
	if err := op.Wait(ctx); err != nil {
		return err
	}
	fmt.Fprintf(w, "Created Venues table\n")

	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	var updateCount int64
	insertSQL := `INSERT INTO Venues (VenueId, Name, Revenues) VALUES ($1, $2, $3)`

	// Insert a Venue using DML.
	insertStmt := spanner.Statement{
		SQL: insertSQL,
		Params: map[string]interface{}{
			"p1": 1,
			"p2": "Venue 1",
			"p3": spanner.PGNumeric{Numeric: "3150.25", Valid: true},
		},
	}
	_, err = client.ReadWriteTransaction(ctx, func(ctx context.Context, transaction *spanner.ReadWriteTransaction) error {
		updateCount, err = transaction.Update(ctx, insertStmt)
		return err
	})
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Inserted %d venue(s)\n", updateCount)

	// Insert a Venue with a NULL value for the Revenues column.
	nullRevenueStmt := spanner.Statement{
		SQL: insertSQL,
		Params: map[string]interface{}{
			"p1": 2,
			"p2": "Venue 2",
			"p3": spanner.PGNumeric{Valid: false},
		},
	}
	_, err = client.ReadWriteTransaction(ctx, func(ctx context.Context, transaction *spanner.ReadWriteTransaction) error {
		updateCount, err = transaction.Update(ctx, nullRevenueStmt)
		return err
	})
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Inserted %d venue(s) with NULL revenues\n", updateCount)

	// Insert a Venue with a NaN (Not a Number) value for the Revenues column.
	nanRevenueStmt := spanner.Statement{
		SQL: insertSQL,
		Params: map[string]interface{}{
			"p1": 3,
			"p2": "Venue 3",
			"p3": spanner.PGNumeric{Numeric: "NaN", Valid: true},
		},
	}
	_, err = client.ReadWriteTransaction(ctx, func(ctx context.Context, transaction *spanner.ReadWriteTransaction) error {
		updateCount, err = transaction.Update(ctx, nanRevenueStmt)
		return err
	})
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Inserted %d venue(s) with NaN revenues\n", updateCount)

	// Get all Venues and inspect the Revenues values.
	iter := client.Single().Query(ctx, spanner.Statement{
		SQL: "SELECT Name, Revenues FROM Venues",
	})
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		var name string
		var revenues spanner.PGNumeric
		if err := row.Columns(&name, &revenues); err != nil {
			return err
		}
		fmt.Fprintf(w, "Revenues of %s: %s\n", name, revenues)
	}

	// Mutations can also be used to insert/update NUMERIC values, including NaN values.
	ts, err := client.Apply(ctx, []*spanner.Mutation{
		spanner.InsertMap("Venues", map[string]interface{}{
			"VenueId":  4,
			"Name":     "Venue 4",
			"Revenues": spanner.PGNumeric{Numeric: "125.10", Valid: true},
		}),
		spanner.InsertMap("Venues", map[string]interface{}{
			"VenueId":  5,
			"Name":     "Venue 5",
			"Revenues": spanner.PGNumeric{Numeric: "NaN", Valid: true},
		}),
	})
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Inserted 2 Venues using mutations at %s\n", ts)

	return nil
}

// [END spanner_postgresql_numeric_data_type]
