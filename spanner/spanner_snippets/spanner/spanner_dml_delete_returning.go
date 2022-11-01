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

// [START spanner_dml_delete_returning]

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
)

func deleteUsingDMLReturning(w io.Writer, db string) error {
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	// Delete records from SINGERS table satisfying a
	// particular condition and returns the SingerId
	// and FullName column of the deleted records using
	// 'THEN RETURN SingerId, FullName'.
	// It is also possible to return all columns of all the
	// deleted records by using 'THEN RETURN *'.
	_, err = client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		stmt := spanner.Statement{
			SQL: `DELETE FROM Singers WHERE FirstName = 'Alice'
			        THEN RETURN SingerId, FullName`,
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
			var (
				singerID int64
				fullName string
			)
			if err := row.Columns(&singerID, &fullName); err != nil {
				return err
			}
			fmt.Fprintf(w, "%d %s\n", singerID, fullName)
		}
		fmt.Fprintf(w, "%d record(s) deleted.\n", iter.RowCount)
		return nil
	})
	return err
}

// [END spanner_dml_delete_returning]
