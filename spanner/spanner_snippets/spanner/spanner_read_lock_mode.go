// Copyright 2026 Google LLC
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

// [START spanner_read_lock_mode]

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/spanner"
	pb "cloud.google.com/go/spanner/apiv1/spannerpb"
)

// writeWithTransactionUsingReadLockMode sets the ReadLockMode globally
// by using ClientConfig and shows how to override it for a specific
// transaction. ReadLockMode determines the locking strategy used during
// transaction execution.
func writeWithTransactionUsingReadLockMode(w io.Writer, db string) error {
	ctx := context.Background()

    // Client-level configuration: Applies to all read-write transactions
	// for this client. OPTIMISTIC mode avoids locks during reads and
	// verifies changes during the commit phase.
	cfg := spanner.ClientConfig{
		TransactionOptions: spanner.TransactionOptions{
			ReadLockMode: pb.TransactionOptions_ReadWrite_OPTIMISTIC,
		},
	}
	client, err := spanner.NewClientWithConfig(ctx, db, cfg)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer client.Close()

    // Transaction-level options take precedence over client-level
	// configuration. PESSIMISTIC mode is used here to override the
	// client-level setting and ensure immediate locking during reads.
	txnOpts := spanner.TransactionOptions{
		ReadLockMode: pb.TransactionOptions_ReadWrite_PESSIMISTIC,
	}

	_, err = client.ReadWriteTransactionWithOptions(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
        // In PESSIMISTIC mode with SERIALIZABLE isolation, the transaction
		// acquires a shared lock during this read.
		key := spanner.Key{1, 2}
		row, err := txn.ReadRow(ctx, "Albums", key, []string{"AlbumTitle"})
		if err != nil {
			return fmt.Errorf("failed to read album: %w", err)
		}
		var title string
		if err := row.Column(0, &title); err != nil {
			return fmt.Errorf("failed to get album title: %w", err)
		}
		fmt.Fprintf(w, "Current album title: %s\n", title)

		// Update the album title
		stmt := spanner.Statement{
			SQL: `UPDATE Albums
				SET AlbumTitle = @AlbumTitle
				WHERE SingerId = @SingerId AND AlbumId = @AlbumId`,
			Params: map[string]interface{}{
				"SingerId":   1,
				"AlbumId":    2,
				"AlbumTitle": "New Album Title",
			},
		}
		count, err := txn.Update(ctx, stmt)
		if err != nil {
			return fmt.Errorf("failed to update album: %w", err)
		}
		fmt.Fprintf(w, "Updated %d record(s).\n", count)
		return nil
	}, txnOpts)

	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}
	return nil
}

// [END spanner_read_lock_mode]
