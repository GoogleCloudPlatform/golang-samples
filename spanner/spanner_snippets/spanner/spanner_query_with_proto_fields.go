// Copyright 2022 Google LLC
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

// [START spanner_query_with_proto_fields]

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/spanner"
	pb "github.com/GoogleCloudPlatform/golang-samples/spanner/spanner_snippets/spanner/testdata"
	"google.golang.org/api/iterator"
)

// This covers Query Proto Columns, Query ENUM Columns, Query proto fields using dot operator
func queryWithProtoFields(w io.Writer, db string) error {
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	stmt := spanner.Statement{
		SQL: `SELECT BookInfo.Title, BookInfo.Genre, BookInfo.Author, BookGenre FROM Library`,
	}

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
		var bookInfoProtoTitle string
		var bookInfoProtoGenre pb.Genre
		var bookInfoProtoAuthor string
		var bookGenre pb.Genre
		if err := row.Columns(&bookInfoProtoTitle, &bookInfoProtoGenre, &bookInfoProtoAuthor, &bookGenre); err != nil {
			return err
		}
		fmt.Fprintf(w, "%s %s %s %s\n", bookInfoProtoTitle, bookInfoProtoGenre, bookInfoProtoAuthor, bookGenre)
	}
}

// [END spanner_query_with_proto_fields]
