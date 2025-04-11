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
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
)

func queryWithHistory(w io.Writer, db string) error {
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	stmt := spanner.Statement{
		SQL: `SELECT d.UserId, d.DocumentId, d.Contents, dh.Timestamp, dh.PreviousContents
				FROM Documents d JOIN DocumentHistory dh
				ON dh.UserId = d.UserId AND dh.DocumentId = d.DocumentId
				ORDER BY dh.Timestamp DESC LIMIT 3`}
	iter := client.Single().Query(ctx, stmt)
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			return nil
		}
		if err != nil {
			return err
		}
		var userID, documentID int64
		var timestamp time.Time
		var contents, previousContents string
		if err := row.Columns(&userID, &documentID, &contents, &timestamp, &previousContents); err != nil {
			return err
		}
		fmt.Fprintf(w, "%d %d %s %s %s\n", userID, documentID, contents, timestamp, previousContents)
	}
}
