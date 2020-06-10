// Copyright 2020 Google LLC
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

// [START spanner_dml_batch_update]

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/spanner"
)

func updateUsingBatchDML(w io.Writer, db string) error {
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	_, err = client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		stmts := []spanner.Statement{
			{SQL: `INSERT INTO Albums
				(SingerId, AlbumId, AlbumTitle, MarketingBudget)
				VALUES (1, 3, 'Test Album Title', 10000)`},
			{SQL: `UPDATE Albums
				SET MarketingBudget = MarketingBudget * 2
				WHERE SingerId = 1 and AlbumId = 3`},
		}
		rowCounts, err := txn.BatchUpdate(ctx, stmts)
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "Executed %d SQL statements using Batch DML.\n", len(rowCounts))
		return nil
	})
	return err
}

// [END spanner_dml_batch_update]
