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

// writeWithTransactionUsingReadLockMode demonstrates how to set the
// ReadLockMode for Cloud Spanner transactions. It shows how to configure a
// default mode at the client level and how to override it for specific
// transactions. The ReadLockMode dictates the locking strategy employed during
// the transaction's execution, influencing concurrency and performance.

// Key ReadLockMode options:
// - PESSIMISTIC: Spanner acquires locks on data as it is read, following a
//   traditional locking model. The exact locking behavior depends on the
//   isolation level (SERIALIZABLE or REPEATABLE_READ).
// - OPTIMISTIC: Spanner avoids acquiring locks during the read phase.
//   Instead, it verifies at commit time whether the data read has remained
//   unchanged. This can improve concurrency but may result in transaction
//   aborts if conflicting modifications are detected.
// - READ_LOCK_MODE_UNSPECIFIED: The default behavior. This typically defaults
//   to PESSIMISTIC for SERIALIZABLE isolation and OPTIMISTIC for
//   REPEATABLE_READ isolation.

// Please refer to the official Spanner concurrency control documentation for
// in-depth details on how ReadLockMode interacts with different isolation
// levels.

// [START spanner_read_lock_mode]

import (
"context"
"fmt"
"io"

"cloud.google.com/go/spanner"
pb "cloud.google.com/go/spanner/apiv1/spannerpb"
)

func writeWithTransactionUsingReadLockMode(w io.Writer, db string) error {
	ctx := context.Background()

	// The read lock mode specified at the client-level will be applied
	// to all RW transactions.
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

	// The read lock mode specified at the transaction-level takes precedence over
	// the read lock mode configured at the client-level. PESSIMISTIC is used here
	// to demonstrate overriding the client-level setting.
	txnOpts := spanner.TransactionOptions{
		ReadLockMode: pb.TransactionOptions_ReadWrite_PESSIMISTIC,
	}

	_, err = client.ReadWriteTransactionWithOptions(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		// Since PESSIMISTIC lock mode with the default isolation level
		// (SERIALIZABLE) is used here, the lock will be acquired immediately after
		// this read.
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

