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

// [START spanner_postgresql_insert_dml_returning]

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
)

func pgInsertUsingDMLReturning(w io.Writer, db string) error {
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	// Insert records into the SINGERS table and returns the
	// generated column FullName of the inserted records using
	// 'RETURNING FullName'.
	// It is also possible to return all columns of all the
	// inserted records by using 'RETURNING *'.
	_, err = client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		stmt := spanner.Statement{
			SQL: `INSERT INTO Singers (SingerId, FirstName, LastName)
			        VALUES (21, 'Melissa', 'Garcia'),
			               (22, 'Russell', 'Morales'),
			               (23, 'Jacqueline', 'Long'),
			               (24, 'Dylan', 'Shaw')
			        RETURNING FullName`,
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
			var fullName string
			if err := row.Columns(&fullName); err != nil {
				return err
			}
			fmt.Fprintf(w, "%s\n", fullName)
		}
		fmt.Fprintf(w, "%d record(s) inserted.\n", iter.RowCount)
		return nil
	})
	return err
}

// [END spanner_postgresql_insert_dml_returning]
