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

// [START spanner_create_sequence]

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/spanner"
	database "cloud.google.com/go/spanner/admin/database/apiv1"
	adminpb "cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
	"google.golang.org/api/iterator"
)

func createSequence(ctx context.Context, w io.Writer, db string) error {
	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return err
	}
	defer adminClient.Close()

	ddl := []string{
		"CREATE SEQUENCE Seq OPTIONS (sequence_kind = 'bit_reversed_positive')",
		"CREATE TABLE Customers (CustomerId INT64 DEFAULT (GET_NEXT_SEQUENCE_VALUE('Seq')), CustomerName STRING(1024)) PRIMARY KEY (CustomerId)",
	}
	op, err := adminClient.UpdateDatabaseDdl(ctx, &adminpb.UpdateDatabaseDdlRequest{
		Database:   db,
		Statements: ddl,
	})
	if err != nil {
		return err
	}
	if err := op.Wait(ctx); err != nil {
		return err
	}
	fmt.Fprintf(w, "Created Customers table with bit reverse sequence keys\n")

	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	_, err = client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		stmt := spanner.Statement{
			SQL: `INSERT INTO Customers (CustomerName) VALUES ('Alice'), ('David'), ('Marc') THEN RETURN CustomerId`,
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
		fmt.Fprintf(w, "%d record(s) inserted.\n", iter.RowCount)
		return nil
	})
	return err
}

// [END spanner_create_sequence]
