// Copyright 2025 Google LLC
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

// [START spanner_postgresql_dml_last_statement]

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/spanner"
)

// Inserts a record into the Singers table and then updates
// the row while also setting the update DML as the last
// statement.
func pgInsertAndUpdateDmlWithLastStatement(w io.Writer, db string) error {
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	_, err = client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		insertStmt := spanner.Statement{
			SQL: `INSERT INTO Singers (SingerId, FirstName, LastName)
					VALUES (54214, 'John', 'Do')`,
		}
		insertRowCount, err := txn.Update(ctx, insertStmt)
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "%d record(s) inserted.\n", insertRowCount)

		updateStmt := spanner.Statement{
			SQL: `UPDATE Singers SET LastName = 'Doe' WHERE SingerId = 54214`,
		}
		opts := spanner.QueryOptions{LastStatement: true}
		updateRowCount, err := txn.UpdateWithOptions(ctx, updateStmt, opts)
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "%d record(s) updated.\n", updateRowCount)
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// [END spanner_postgresql_dml_last_statement]
