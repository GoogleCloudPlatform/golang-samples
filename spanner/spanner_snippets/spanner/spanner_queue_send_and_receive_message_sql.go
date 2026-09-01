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

// [START spanner_queue_send_and_receive_message_sql]

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
)

func sendAndReceiveQueueMessageSQL(w io.Writer, db string) error {
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	_, err = client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		stmt := spanner.Statement{
			SQL: `INSERT INTO MyQueue (Id, Payload) VALUES (5, b'message5')`,
		}
		_, err := txn.Update(ctx, stmt)
		return err
	})
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Message sent to queue via SQL\n")

	fmt.Fprintf(w, "Receiving message from queue (max_duration 1min)...\n")
	stmt := spanner.Statement{
		SQL: `SELECT * FROM RECEIVE_MyQueue(max_duration => '1m')`,
	}
	iter := client.Single().Query(ctx, stmt)
	defer iter.Stop()

	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		var id int64
		var payload []byte
		if err := row.ColumnByName("Id", &id); err != nil {
			return err
		}
		if err := row.ColumnByName("Payload", &payload); err != nil {
			return err
		}
		fmt.Fprintf(w, "Received message ID: %d, Payload: %s\n", id, string(payload))
	}

	return nil
}

// [END spanner_queue_send_and_receive_message_sql]
