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

// [START spanner_postgresql_update_dml_returning]

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
)

func pgUpdateUsingDMLReturning(w io.Writer, db string) error {
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	// Update MarketingBudget column for records satisfying
	// a particular condition and returns the modified
	// MarketingBudget column of the updated records using
	// 'RETURNING MarketingBudget'.
	// It is also possible to return all columns of all the
	// updated records by using 'RETURNING *'.
	_, err = client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		stmt := spanner.Statement{
			SQL: `UPDATE Albums
				SET MarketingBudget = MarketingBudget * 2
				WHERE SingerId = 1 and AlbumId = 1
				RETURNING MarketingBudget`,
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
			var marketingBudget int64
			if err := row.Columns(&marketingBudget); err != nil {
				return err
			}
			fmt.Fprintf(w, "%d\n", marketingBudget)
		}
		fmt.Fprintf(w, "%d record(s) updated.\n", iter.RowCount)
		return nil
	})
	return err
}

// [END spanner_postgresql_update_dml_returning]
