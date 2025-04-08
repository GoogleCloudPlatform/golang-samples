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

// [START spanner_dml_getting_started_update]

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/spanner"
)

func writeWithTransactionUsingDML(w io.Writer, db string) error {
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	_, err = client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		// getBudget returns the budget for a record with a given albumId and singerId.
		getBudget := func(albumID, singerID int64) (int64, error) {
			key := spanner.Key{albumID, singerID}
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
		// updateBudget updates the budget for a record with a given albumId and singerId.
		updateBudget := func(singerID, albumID, albumBudget int64) error {
			stmt := spanner.Statement{
				SQL: `UPDATE Albums
					SET MarketingBudget = @AlbumBudget
					WHERE SingerId = @SingerId and AlbumId = @AlbumId`,
				Params: map[string]interface{}{
					"SingerId":    singerID,
					"AlbumId":     albumID,
					"AlbumBudget": albumBudget,
				},
			}
			_, err := txn.Update(ctx, stmt)
			return err
		}

		// Transfer the marketing budget from one album to another. By keeping the actions
		// in a single transaction, it ensures the movement is atomic.
		const transferAmt = 200000
		album2Budget, err := getBudget(2, 2)
		if err != nil {
			return err
		}
		// The transaction will only be committed if this condition still holds at the time
		// of commit. Otherwise it will be aborted and the callable will be rerun by the
		// client library.
		if album2Budget >= transferAmt {
			album1Budget, err := getBudget(1, 1)
			if err != nil {
				return err
			}
			if err = updateBudget(1, 1, album1Budget+transferAmt); err != nil {
				return err
			}
			if err = updateBudget(2, 2, album2Budget-transferAmt); err != nil {
				return err
			}
			fmt.Fprintf(w, "Moved %d from Album2's MarketingBudget to Album1's.", transferAmt)
		}
		return nil
	})
	return err
}

// [END spanner_dml_getting_started_update]
