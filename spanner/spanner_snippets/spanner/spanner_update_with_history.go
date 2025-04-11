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

import (
	"context"
	"io"

	"cloud.google.com/go/spanner"
)

func updateWithHistory(w io.Writer, db string) error {
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	_, err = client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		// Create anonymous function "getContents" to read the current value of the Contents column for a given row.
		getContents := func(key spanner.Key) (string, error) {
			row, err := txn.ReadRow(ctx, "Documents", key, []string{"Contents"})
			if err != nil {
				return "", err
			}
			var content string
			if err := row.Column(0, &content); err != nil {
				return "", err
			}
			return content, nil
		}
		// Create two string arrays corresponding to the columns in each table.
		documentsColumns := []string{"UserId", "DocumentId", "Contents"}
		documentHistoryColumns := []string{"UserId", "DocumentId", "Timestamp", "PreviousContents"}
		// Get row's Contents before updating.
		previousContents, err := getContents(spanner.Key{1, 1})
		if err != nil {
			return err
		}
		// Update row's Contents while saving previous Contents in DocumentHistory table.
		txn.BufferWrite([]*spanner.Mutation{
			spanner.InsertOrUpdate("Documents", documentsColumns,
				[]interface{}{1, 1, "Hello World 1 Updated"}),
			spanner.InsertOrUpdate("DocumentHistory", documentHistoryColumns,
				[]interface{}{1, 1, spanner.CommitTimestamp, previousContents}),
		})
		previousContents, err = getContents(spanner.Key{1, 3})
		if err != nil {
			return err
		}
		txn.BufferWrite([]*spanner.Mutation{
			spanner.InsertOrUpdate("Documents", documentsColumns,
				[]interface{}{1, 3, "Hello World 3 Updated"}),
			spanner.InsertOrUpdate("DocumentHistory", documentHistoryColumns,
				[]interface{}{1, 3, spanner.CommitTimestamp, previousContents}),
		})
		previousContents, err = getContents(spanner.Key{2, 5})
		if err != nil {
			return err
		}
		txn.BufferWrite([]*spanner.Mutation{
			spanner.InsertOrUpdate("Documents", documentsColumns,
				[]interface{}{2, 5, "Hello World 5 Updated"}),
			spanner.InsertOrUpdate("DocumentHistory", documentHistoryColumns,
				[]interface{}{2, 5, spanner.CommitTimestamp, previousContents}),
		})
		return nil
	})
	return err
}
