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

// [START spanner_query_with_proto_field_parameter]

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/spanner"
	pb "github.com/GoogleCloudPlatform/golang-samples/spanner/spanner_snippets/spanner/testdata"
	"google.golang.org/api/iterator"
)

func queryWithProtoFieldParameter(w io.Writer, db string) error {
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	// Filtering on Proto message fields
	var exampleAuthor = "Ron"
	stmt := spanner.Statement{
		SQL: `SELECT Id, BookInfo, BookGenre FROM Library
	            	WHERE BookInfo.Author = @author`,
		Params: map[string]interface{}{
			"author": exampleAuthor,
		},
	}

	/*var exampleInt = 0
		stmt := spanner.Statement{
			SQL: `SELECT Id, BookInfo, BookGenre FROM Library
	            	WHERE BookInfo.Isbn >= @num`,
			Params: map[string]interface{}{
				"num": exampleInt,
			},
		}*/

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
		var bookID int64
		bookProto := &pb.Book{}
		var genre pb.Genre
		if err := row.Columns(&bookID, bookProto, &genre); err != nil {
			return err
		}
		fmt.Fprintf(w, "%d %s %s\n", bookID, bookProto, genre)
	}
}

// [END spanner_query_with_proto_field_parameter]
