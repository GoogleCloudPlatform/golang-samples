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

// [START spanner_batch_client_request_priority]

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/spanner"
	sppb "cloud.google.com/go/spanner/apiv1/spannerpb"
	"google.golang.org/api/iterator"
)

// readBatchDataRequestPriority reads a database in batch read-only transaction with request priority set
func readBatchDataRequestPriority(w io.Writer, db string) error {
	// db = `projects/<project>/instances/<instance-id>/database/<database-id>`
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	txn, err := client.BatchReadOnlyTransaction(ctx, spanner.StrongRead())
	if err != nil {
		return err
	}
	defer txn.Close()

	// Singer represents a row in the Singers table.
	type Singer struct {
		SingerID   int64
		FirstName  string
		LastName   string
		SingerInfo []byte
	}
	stmt := spanner.Statement{SQL: "SELECT SingerId, FirstName, LastName FROM Singers;"}
	partitions, err := txn.PartitionQueryWithOptions(ctx, stmt, spanner.PartitionOptions{}, spanner.QueryOptions{Priority: sppb.RequestOptions_PRIORITY_LOW})
	if err != nil {
		return err
	}
	recordCount := 0
	for i, p := range partitions {
		iter := txn.Execute(ctx, p)
		defer iter.Stop()
		for {
			row, err := iter.Next()
			if err == iterator.Done {
				break
			} else if err != nil {
				return err
			}
			var s Singer
			if err := row.ToStruct(&s); err != nil {
				return err
			}
			fmt.Fprintf(w, "Partition (%d) %v\n", i, s)
			recordCount++
		}
	}
	fmt.Fprintf(w, "Total partition count: %v\n", len(partitions))
	fmt.Fprintf(w, "Total record count: %v\n", recordCount)
	return nil
}

// [END spanner_batch_client_request_priority]
