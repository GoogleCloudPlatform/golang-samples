// Copyright 2024 Google LLC
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

// [START spanner_transaction_timeout]

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"io"
	"time"

	"cloud.google.com/go/spanner"
)

func transactionTimeout(w io.Writer, db string) error {
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	// Create a context with a 60-second timeout and use this context to run a read/write transaction.
	// This context timeout will be applied to the entire transaction, and the transaction will fail
	// if it cannot finish within the specified timeout value.
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	_, err = client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		selectStmt := spanner.Statement{
			SQL: `SELECT SingerId, FirstName, LastName FROM Singers ORDER BY LastName, FirstName`,
		}
		iter := txn.Query(ctx, selectStmt)
		defer iter.Stop()
		for {
			row, err := iter.Next()
			if errors.Is(err, iterator.Done) {
				break
			}
			if err != nil {
				return err
			}
			var singerID int64
			var firstName, lastName string
			if err := row.Columns(&singerID, &firstName, &lastName); err != nil {
				return err
			}
			fmt.Fprintf(w, "%d %s %s\n", singerID, firstName, lastName)
		}
		stmt := spanner.Statement{
			SQL: `INSERT INTO Singers (SingerId, FirstName, LastName)
					VALUES (20, 'George', 'Washington')`,
		}
		rowCount, err := txn.Update(ctx, stmt)
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "%d record(s) inserted.\n", rowCount)
		return nil
	})
	// Check if an error was returned by the transaction.
	// The spanner.ErrCode(err) function will return codes.OK if err == nil.
	code := spanner.ErrCode(err)
	if code == codes.OK {
		fmt.Fprintf(w, "Transaction with timeout was executed successfully\n")
	} else if code == codes.DeadlineExceeded {
		fmt.Fprintf(w, "Transaction timed out\n")
	} else {
		fmt.Fprintf(w, "Transaction failed with error code %v\n", code)
	}
	return err
}

// [END spanner_transaction_timeout]
