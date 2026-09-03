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

// [START spanner_queue_send_sql_future]

import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/spanner"
)

func sendToQueueSQLFuture(w io.Writer, db string) error {
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	futureTime := time.Now().Add(10 * time.Minute)
	_, err = client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		stmt := spanner.Statement{
			SQL: `INSERT INTO MyQueue (Id, Payload, DeliverTime) VALUES (4, b'message4', @deliveryTime)`,
			Params: map[string]interface{}{
				"deliveryTime": futureTime,
			},
		}
		_, err := txn.Update(ctx, stmt)
		return err
	})
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Message scheduled for future delivery via SQL\n")
	return nil
}

// [END spanner_queue_send_sql_future]
