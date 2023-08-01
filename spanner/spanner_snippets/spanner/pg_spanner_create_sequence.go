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

// [START spanner_postgresql_create_sequence]

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/spanner"
	database "cloud.google.com/go/spanner/admin/database/apiv1"
	adminpb "cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
	"google.golang.org/api/iterator"
)

func pgCreateSequence(w io.Writer, db string) error {
	// db := "projects/my-project/instances/my-instance/databases/my-database"
	ctx := context.Background()
	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return err
	}
	defer adminClient.Close()

	// List of DDL statements to be applied to the database.
	// Create a sequence, and then use the sequence as auto generated primary key in Customers table.
	ddl := []string{
		"CREATE SEQUENCE Seq BIT_REVERSED_POSITIVE",
		"CREATE TABLE Customers (CustomerId BIGINT DEFAULT nextval('Seq'), CustomerName character varying(1024), PRIMARY KEY (CustomerId))",
	}
	op, err := adminClient.UpdateDatabaseDdl(ctx, &adminpb.UpdateDatabaseDdlRequest{
		Database:   db,
		Statements: ddl,
	})
	if err != nil {
		return err
	}
	// Wait for the UpdateDatabaseDdl operation to finish.
	if err := op.Wait(ctx); err != nil {
		return fmt.Errorf("waiting for bit reverse sequence creation to finish failed: %w", err)
	}
	fmt.Fprintf(w, "Created Seq sequence and Customers table, where its key column CustomerId uses the sequence as a default value\n")

	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	// Inserts records into the Customers table.
	// The ReadWriteTransaction function returns the commit timestamp and an error.
	// The commit timestamp is ignored in this case.
	_, err = client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		stmt := spanner.Statement{
			SQL: `INSERT INTO Customers (CustomerName) VALUES ('Alice'), ('David'), ('Marc') RETURNING CustomerId`,
		}
		iter := txn.Query(ctx, stmt)
		defer iter.Stop()
		for {
			row, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return err
			}
			var customerId int64
			if err := row.Columns(&customerId); err != nil {
				return err
			}
			fmt.Fprintf(w, "Inserted customer record with CustomerId: %d\n", customerId)
		}
		fmt.Fprintf(w, "Number of customer records inserted is: %d\n", iter.RowCount)
		return nil
	})
	return err
}

// [END spanner_postgresql_create_sequence]
