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

// [START spanner_read_write_transaction]

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/spanner"
)

func writeWithTransaction(w io.Writer, db string) error {
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	_, err = client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		getBudget := func(key spanner.Key) (int64, error) {
			row, err := txn.ReadRow(ctx, "Albums", key, []string{"MarketingBudget"})
			if err != nil {
				return 0, err
			}
			var budget int64
			if err := row.Column(0, &budget); err != nil {
				return 0, err
			}
			return budget, nil
		}
		album2Budget, err := getBudget(spanner.Key{2, 2})
		if err != nil {
			return err
		}
		const transferAmt = 200000
		if album2Budget >= transferAmt {
			album1Budget, err := getBudget(spanner.Key{1, 1})
			if err != nil {
				return err
			}
			album1Budget += transferAmt
			album2Budget -= transferAmt
			cols := []string{"SingerId", "AlbumId", "MarketingBudget"}
			txn.BufferWrite([]*spanner.Mutation{
				spanner.Update("Albums", cols, []interface{}{1, 1, album1Budget}),
				spanner.Update("Albums", cols, []interface{}{2, 2, album2Budget}),
			})
			fmt.Fprintf(w, "Moved %d from Album2's MarketingBudget to Album1's.", transferAmt)
		}
		return nil
	})
	return err
}

// [END spanner_read_write_transaction]
