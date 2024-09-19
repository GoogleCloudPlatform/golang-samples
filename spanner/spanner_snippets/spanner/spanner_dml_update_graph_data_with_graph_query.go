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

// [START spanner_update_graph_data_with_graph_query_in_dml]

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/spanner"
)

func updateGraphDataWithGraphQueryInDml(w io.Writer, db string) error {
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	// Execute a ReadWriteTransaction to update the 'Account' table underpinning
	// 'Account' nodes in 'FinGraph'. The function run by ReadWriteTransaction
	// executes an 'UPDATE' SQL DML statement. Graph queries run after this
	// transaction is committed will observe the effects of the updates to 'Account's
	//
	// The update is performed for all 'Account's whose 'id' is returned by
	// the graph query in the 'IN' subquery, i.e., all 'Account's that have
	// received transfers directly or via one intermediary from an 'Account'
	// whose 'id' is 1.
	_, err = client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		stmt := spanner.Statement{
			SQL: `UPDATE Account SET is_blocked = true 
            	  WHERE id IN {
            	    GRAPH FinGraph 
            	    MATCH (a:Account WHERE a.id = 1)-[:TRANSFERS]->{1,2}(b:Account)
            	    RETURN b.id}`,
		}
		rowCount, err := txn.Update(ctx, stmt)
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "%d Account record(s) updated.\n", rowCount)
		return err
	})

	return err
}

// [END spanner_update_graph_data_with_graph_query_in_dml]
