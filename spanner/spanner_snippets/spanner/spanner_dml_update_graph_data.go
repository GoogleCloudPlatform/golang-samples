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

// [START spanner_update_graph_data_with_dml]

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/spanner"
)

func updateGraphDataWithDml(w io.Writer, db string) error {
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	// Execute a ReadWriteTransaction to update the 'Account' table underpinning
	// 'Account' nodes in 'FinGraph'. The function run by ReadWriteTransaction
	// executes an 'UPDATE' SQL DML statement. Graph queries run after this
	// transaction is committed will observe the effects of the update to 'Account'
	// with 'id' = 2.
	_, err1 := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		stmt := spanner.Statement{
			SQL: `UPDATE Account SET is_blocked = false WHERE id = 2`,
		}
		rowCount, err := txn.Update(ctx, stmt)
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "%d Account record(s) updated.\n", rowCount)
		return err
	})

	if err1 != nil {
		return err1
	}

	// Execute a ReadWriteTransaction to update the 'AccountTransferAccount' table
	// underpinning 'AccountTransferAccount' edges in 'FinGraph'. The function run
	// by ReadWriteTransaction executes an 'UPDATE' SQL DML statement.
	// Graph queries run after this transaction is committed will observe the effects
	// of the update to 'AccountTransferAccount' where the source of the transfer has
	// 'id' 1 and the destination has 'id' 2.
	_, err2 := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		stmt := spanner.Statement{
			SQL: `UPDATE AccountTransferAccount SET amount = 300 WHERE id = 1 AND to_id = 2`,
		}
		rowCount, err := txn.Update(ctx, stmt)
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "%d AccountTransferAccount record(s) updated.\n", rowCount)
		return err
	})

	return err2
}

// [END spanner_update_graph_data_with_dml]
