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

import (
	"context"
	"io"

	"cloud.google.com/go/spanner"
	pb "github.com/GoogleCloudPlatform/golang-samples/spanner/spanner_snippets/spanner/testdata"
)

// [START spanner_insert_data_with_proto_columns]

func insertDataWithProtoMsgAndEnum(w io.Writer, db string) error {
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	book1 := &pb.Book{
		Isbn:   1,
		Title:  "Harry Potter",
		Author: "JK Rowling",
		Genre:  pb.Genre_CLASSICAL,
	}

	book2 := &pb.Book{
		Isbn:   2,
		Title:  "New Arrival",
		Author: "Ron",
		Genre:  pb.Genre_ROCK,
	}

	cols := []string{"Id", "BookInfo", "BookGenre"}
	m := []*spanner.Mutation{
		spanner.InsertOrUpdate("Library", cols, []interface{}{1, book1, pb.Genre_CLASSICAL}),
		spanner.InsertOrUpdate("Library", cols, []interface{}{2, book2, pb.Genre_ROCK}),
	}
	_, err = client.Apply(ctx, m)
	return err
}

// [END spanner_insert_data_with_proto_columns]
