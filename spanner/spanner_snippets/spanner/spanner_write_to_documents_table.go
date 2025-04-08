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

func writeToDocumentsTable(w io.Writer, db string) error {
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	documentsColumns := []string{"UserId", "DocumentId", "Timestamp", "Contents"}
	m := []*spanner.Mutation{
		spanner.InsertOrUpdate("DocumentsWithTimestamp", documentsColumns,
			[]interface{}{1, 1, spanner.CommitTimestamp, "Hello World 1"}),
		spanner.InsertOrUpdate("DocumentsWithTimestamp", documentsColumns,
			[]interface{}{1, 2, spanner.CommitTimestamp, "Hello World 2"}),
		spanner.InsertOrUpdate("DocumentsWithTimestamp", documentsColumns,
			[]interface{}{1, 3, spanner.CommitTimestamp, "Hello World 3"}),
		spanner.InsertOrUpdate("DocumentsWithTimestamp", documentsColumns,
			[]interface{}{2, 4, spanner.CommitTimestamp, "Hello World 4"}),
		spanner.InsertOrUpdate("DocumentsWithTimestamp", documentsColumns,
			[]interface{}{2, 5, spanner.CommitTimestamp, "Hello World 5"}),
		spanner.InsertOrUpdate("DocumentsWithTimestamp", documentsColumns,
			[]interface{}{3, 6, spanner.CommitTimestamp, "Hello World 6"}),
		spanner.InsertOrUpdate("DocumentsWithTimestamp", documentsColumns,
			[]interface{}{3, 7, spanner.CommitTimestamp, "Hello World 7"}),
		spanner.InsertOrUpdate("DocumentsWithTimestamp", documentsColumns,
			[]interface{}{3, 8, spanner.CommitTimestamp, "Hello World 8"}),
		spanner.InsertOrUpdate("DocumentsWithTimestamp", documentsColumns,
			[]interface{}{3, 9, spanner.CommitTimestamp, "Hello World 9"}),
		spanner.InsertOrUpdate("DocumentsWithTimestamp", documentsColumns,
			[]interface{}{3, 10, spanner.CommitTimestamp, "Hello World 10"}),
	}
	_, err = client.Apply(ctx, m)
	return err
}
