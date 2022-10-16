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
	"fmt"
	"io"

	"cloud.google.com/go/spanner"
	pb "github.com/GoogleCloudPlatform/golang-samples/spanner/spanner_snippets/spanner/testdata"
	"google.golang.org/api/iterator"
	"google.golang.org/protobuf/reflect/protoreflect"

	//"google.golang.org/protobuf/proto"
	"github.com/golang/protobuf/proto"
)

// [START spanner_read_data_with_array_of_proto_column]

func readArrayOfProtoMsgEnum(w io.Writer, db string) error {
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	iter := client.Single().Read(ctx, "Library", spanner.AllKeys(),
		[]string{"Id", "BookInfo", "BookGenre"})
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
		bookArray := []proto.Message{&pb.Book{}}
		genreArray := []protoreflect.Enum{pb.Genre_BLUES}
		if err := row.Columns(&bookID, &bookArray, &genreArray); err != nil {
			return err
		}
		fmt.Fprintf(w, "%d %s %s\n", bookID, bookArray, genreArray)
	}
}

// [END spanner_read_data_with_array_of_proto_column]
