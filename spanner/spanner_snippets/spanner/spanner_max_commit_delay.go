// Copyright 2021 Google LLC
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

// [START spanner_set_max_commit_delay]

import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/spanner"
)

func maxCommitDelay(w io.Writer, db string) error {
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return fmt.Errorf("maxCommitDelay.NewClient: %w", err)
	}
	defer client.Close()

	resp, err := client.ReadWriteTransactionWithOptions(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		stmt := spanner.Statement{
			SQL: `INSERT Singers (SingerId, FirstName, LastName)
					VALUES (111, 'Virginia', 'Watson')`,
		}
		rowCount, err := txn.Update(ctx, stmt)
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "%d record(s) inserted.\n", rowCount)
		return nil
	}, spanner.TransactionOptions{CommitOptions: spanner.CommitOptions{MaxCommitDelay: time.ParseDuration("100ms"}})
	if err != nil {
		return fmt.Errorf("maxCommitDelay.ReadWriteTransactionWithOptions: %w", err)
	}
	return nil
}

// [END spanner_set_max_commit_delay]
