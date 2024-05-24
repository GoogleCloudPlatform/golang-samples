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

// [START spanner_update_data_with_proto_types]
import (
	"context"
	"fmt"
	"io"
	"regexp"

	"cloud.google.com/go/spanner"
	pb "github.com/GoogleCloudPlatform/golang-samples/spanner/spanner_snippets/spanner/testdata/protos"
	"google.golang.org/protobuf/proto"
)

// updateDataWithProtoColumn updates database with Proto type values
func updateDataWithProtoColumn(w io.Writer, db string) error {
	// db = `projects/<project>/instances/<instance-id>/database/<database-id>`
	matches := regexp.MustCompile("^(.*)/databases/(.*)$").FindStringSubmatch(db)
	if matches == nil || len(matches) != 3 {
		return fmt.Errorf("updateDataWithProtoColumn: invalid database id %s", db)
	}

	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	singerGenre := pb.Genre_FOLK
	singerInfo := &pb.SingerInfo{
		SingerId:    proto.Int64(2),
		BirthDate:   proto.String("February"),
		Nationality: proto.String("Country2"),
		Genre:       &singerGenre,
	}
	singerInfoArray := []*pb.SingerInfo{singerInfo}
	singerGenreArray := []*pb.Genre{&singerGenre}

	cols := []string{"SingerId", "SingerInfo", "SingerInfoArray", "SingerGenre", "SingerGenreArray"}
	_, err = client.Apply(ctx, []*spanner.Mutation{
		spanner.Update("Singers", cols, []interface{}{2, singerInfo, singerInfoArray, &singerGenre, singerGenreArray}),
		spanner.Update("Singers", cols, []interface{}{3, nil, nil, nil, nil}),
	})

	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Data updated\n")

	return nil
}

// [END spanner_update_data_with_proto_types]
