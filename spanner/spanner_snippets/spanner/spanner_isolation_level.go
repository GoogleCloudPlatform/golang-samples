// Copyright 2025 Google LLC
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

// [START spanner_isolation_level]

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/spanner"
	pb "cloud.google.com/go/spanner/apiv1/spannerpb"
)

func writeWithTransactionUsingIsolationLevel(w io.Writer, db string) error {
	ctx := context.Background()

	// The isolation level specified at the client-level will be applied
	// to all RW transactions.
	cfg := spanner.ClientConfig{
		TransactionOptions: spanner.TransactionOptions{
			IsolationLevel: pb.TransactionOptions_SERIALIZABLE,
		},
	}
	client, err := spanner.NewClientWithConfig(ctx, db, cfg)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer client.Close()

	// The isolation level specified at the transaction-level takes
	// precedence over the isolation level configured at the client-level.
	// REPEATABLE_READ is used here to demonstrate overriding the client-level setting.
	txnOpts := spanner.TransactionOptions{
		IsolationLevel: pb.TransactionOptions_REPEATABLE_READ,
	}

	_, err = client.ReadWriteTransactionWithOptions(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		// Read the current album title
		key := spanner.Key{1, 1}
		row, err := txn.ReadRow(ctx, "Albums", key, []string{"AlbumTitle"})
		if err != nil {
			return fmt.Errorf("failed to read album: %v", err)
		}
		var title string
		if err := row.Column(0, &title); err != nil {
			return fmt.Errorf("failed to get album title: %v", err)
		}
		fmt.Fprintf(w, "Current album title: %s\n", title)

		// Update the album title
		stmt := spanner.Statement{
			SQL: `UPDATE Albums
				SET AlbumTitle = @AlbumTitle
				WHERE SingerId = @SingerId AND AlbumId = @AlbumId`,
			Params: map[string]interface{}{
				"SingerId":   1,
				"AlbumId":    1,
				"AlbumTitle": "New Album Title",
			},
		}
		count, err := txn.Update(ctx, stmt)
		if err != nil {
			return fmt.Errorf("failed to update album: %v", err)
		}
		fmt.Fprintf(w, "Updated %d record(s).\n", count)
		return nil
	}, txnOpts)

	if err != nil {
		return fmt.Errorf("transaction failed: %v", err)
	}
	return nil
}

// [END spanner_isolation_level]
