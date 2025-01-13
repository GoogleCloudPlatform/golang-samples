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

// [START spanner_set_exclude_txn_from_change_streams]

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/spanner"
)

// rwTxnExcludedFromChangeStreams executes the insert and update DMLs on Singers table excluded from allowed tracking change streams
func rwTxnExcludedFromChangeStreams(w io.Writer, db string) error {
	// db = `projects/<project>/instances/<instance-id>/database/<database-id>`
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return fmt.Errorf("rwTxnExcludedFromChangeStreams.NewClient: %w", err)
	}
	defer client.Close()

	_, err = client.ReadWriteTransactionWithOptions(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		stmt := spanner.Statement{
			SQL: `INSERT Singers (SingerId, FirstName, LastName)
					VALUES (111, 'Virginia', 'Watson')`,
		}
		_, err := txn.Update(ctx, stmt)
		if err != nil {
			return fmt.Errorf("rwTxnExcludedFromChangeStreams.Update: %w", err)
		}
		fmt.Fprintf(w, "New singer inserted.")
		stmt = spanner.Statement{
			SQL: `UPDATE Singers SET FirstName = 'Hi' WHERE SingerId = 111`,
		}
		_, err = txn.Update(ctx, stmt)
		if err != nil {
			return fmt.Errorf("rwTxnExcludedFromChangeStreams.Update: %w", err)
		}
		fmt.Fprint(w, "Singer first name updated.")
		return nil
	}, spanner.TransactionOptions{ExcludeTxnFromChangeStreams: true})
	if err != nil {
		return err
	}
	return nil
}

// [END spanner_set_exclude_txn_from_change_streams]
