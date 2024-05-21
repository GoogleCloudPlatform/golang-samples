// Copyright 2024 Google LLC
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

// [START spanner_query_with_proto_types_parameter]

import (
	"context"
	"fmt"
	"io"
	"regexp"

	"cloud.google.com/go/spanner"
	pb "github.com/GoogleCloudPlatform/golang-samples/spanner/spanner_snippets/spanner/testdata/protos"
	"google.golang.org/api/iterator"
)

// queryWithProtoParameter queries data on the JSON type column of the database
func queryWithProtoParameter(w io.Writer, db string) error {
	// db = `projects/<project>/instances/<instance-id>/database/<database-id>`
	matches := regexp.MustCompile("^(.*)/databases/(.*)$").FindStringSubmatch(db)
	if matches == nil || len(matches) != 3 {
		return fmt.Errorf("queryWithProtoParameter: invalid database id %s", db)
	}
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	stmt := spanner.Statement{
		SQL: `SELECT SingerId, SingerInfo, SingerInfo.nationality, SingerInfoArray, SingerGenre, SingerGenreArray FROM Singers WHERE SingerInfo.Nationality=@country and SingerGenre=@singerGenre`,
		Params: map[string]interface{}{
			"country": "Country2",
			"singerGenre": spanner.NullProtoEnum{
				ProtoEnumVal: pb.Genre_FOLK,
				Valid:        true,
			},
		},
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

		var singerId int64
		singerInfo := &pb.SingerInfo{}
		var nationality string
		var singerGenre pb.Genre
		var singerInfoArray []*pb.SingerInfo
		var singerGenreArray []*pb.Genre
		if err := row.Columns(&singerId, singerInfo, &nationality, &singerInfoArray, &singerGenre, &singerGenreArray); err != nil {
			return err
		}
		fmt.Fprintf(w, "%d %v %v %v %v\n", singerId, singerInfo, singerGenre, singerInfoArray, singerGenreArray)
	}
}

// [END spanner_query_with_proto_types_parameter]
