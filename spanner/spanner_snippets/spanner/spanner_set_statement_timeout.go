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

// [START spanner_set_statement_timeout]

import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/spanner"
	"google.golang.org/grpc/codes"
)

func setStatementTimeout(w io.Writer, db string) error {
	client, err := spanner.NewClient(context.Background(), db)
	if err != nil {
		return err
	}
	defer client.Close()

	_, err = client.ReadWriteTransaction(context.Background(),
		func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
			// Create a context with a 60-second timeout and apply this timeout to the insert statement.
			ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()
			stmt := spanner.Statement{
				SQL: `INSERT Singers (SingerId, FirstName, LastName)
					VALUES (39, 'George', 'Washington')`,
			}
			rowCount, err := txn.Update(ctxWithTimeout, stmt)
			// Get the error code from the error. This function returns codes.OK if err == nil.
			code := spanner.ErrCode(err)
			if code == codes.DeadlineExceeded {
				fmt.Fprintf(w, "Insert statement timed out.\n")
			} else if code == codes.OK {
				fmt.Fprintf(w, "%d record(s) inserted.\n", rowCount)
			} else {
				fmt.Fprintf(w, "Insert statement failed with error %v\n", err)
			}
			return err
		})
	if err != nil {
		return err
	}
	return nil
}

// [END spanner_set_statement_timeout]
