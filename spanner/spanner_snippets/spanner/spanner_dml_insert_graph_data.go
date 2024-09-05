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

// [START spanner_insert_graph_data_with_dml]

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/spanner"
)

func insertGraphDataWithDml(w io.Writer, db string) error {
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	// Execute a ReadWriteTransaction to insert values into the 'Account' table
	// underpinning 'Account' nodes in 'FinGraph'. The function run by ReadWriteTransaction
	// executes an 'INSERT' SQL DML statement. Graph queries run after this
	// transaction is committed will observe the effects of the new 'Account's
	// added to the graph.
	_, err1 := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		stmt := spanner.Statement{
			SQL: `INSERT INTO Account (id, create_time, is_blocked)
            		VALUES
            	    	(1, CAST('2000-08-10 08:18:48.463959-07:52' AS TIMESTAMP), false),
            			(2, CAST('2000-08-12 07:13:16.463959-03:41' AS TIMESTAMP), true)`,
		}
		rowCount, err := txn.Update(ctx, stmt)
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "%d Account record(s) inserted.\n", rowCount)
		return err
	})

	if err1 != nil {
		return err1
	}

	// Execute a ReadWriteTransaction to insert values into the 'AccountTransferAccount'
	// table underpinning 'AccountTransferAccount' edges in 'FinGraph'. The function run
	// by ReadWriteTransaction executes an 'INSERT' SQL DML statement.
	// Graph queries run after this transaction is committed will observe the effects
	// of the edges added to the graph.
	_, err2 := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		stmt := spanner.Statement{
			SQL: `INSERT INTO AccountTransferAccount (id, to_id, create_time, amount)
					VALUES
						(1, 2, CAST('2000-09-11 03:11:18.463959-06:36' AS TIMESTAMP), 100),
						(1, 1, CAST('2000-09-12 04:09:34.463959-05:12' AS TIMESTAMP), 200)`,
		}
		rowCount, err := txn.Update(ctx, stmt)
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "%d AccountTransferAccount record(s) inserted.\n", rowCount)
		return err
	})

	return err2
}

// [END spanner_insert_graph_data_with_dml]
