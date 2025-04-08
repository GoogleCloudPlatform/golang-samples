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

// [START spanner_update_data_with_proto_types_with_dml]
import (
	"context"
	"fmt"
	"io"
	"regexp"

	"cloud.google.com/go/spanner"
	pb "github.com/GoogleCloudPlatform/golang-samples/spanner/spanner_snippets/spanner/testdata/protos"
	"google.golang.org/protobuf/proto"
)

// updateDataWithProtoColumnWithDml updates database with Proto type values using DML
func updateDataWithProtoColumnWithDml(w io.Writer, db string) error {
	// db = `projects/<project>/instances/<instance-id>/database/<database-id>`
	matches := regexp.MustCompile("^(.*)/databases/(.*)$").FindStringSubmatch(db)
	if matches == nil || len(matches) != 3 {
		return fmt.Errorf("updateDataWithProtoColumnWithDml: invalid database id %s", db)
	}

	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	singerGenre := pb.Genre_ROCK
	singerInfo := &pb.SingerInfo{
		SingerId:    proto.Int64(1),
		BirthDate:   proto.String("January"),
		Nationality: proto.String("Country1"),
		Genre:       &singerGenre,
	}
	singerInfoArray := []*pb.SingerInfo{singerInfo}
	singerGenreArray := []*pb.Genre{&singerGenre}

	_, err = client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		stmt := spanner.Statement{
			SQL: `UPDATE Singers SET SingerInfo = @singerInfo,  SingerInfoArray=@singerInfoArray,
                   SingerGenre=@singerGenre, SingerGenreArray=@singerGenreArray WHERE SingerId = 1`,
			Params: map[string]interface{}{
				"singerInfo":       singerInfo,
				"singerInfoArray":  singerInfoArray,
				"singerGenre":      &singerGenre,
				"singerGenreArray": singerGenreArray,
			},
		}

		rowCount, err := txn.Update(ctx, stmt)
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "%d record(s) updated.\n", rowCount)
		return err
	})

	if err != nil {
		return err
	}

	_, err = client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		stmt := spanner.Statement{
			SQL: `UPDATE Singers SET SingerInfo.nationality = @singerNationality WHERE SingerId = 1`,
			Params: map[string]interface{}{
				"singerNationality": "Country2",
			},
		}

		rowCount, err := txn.Update(ctx, stmt)
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "%d record(s) updated.\n", rowCount)
		return err
	})

	return err
}

// [END spanner_update_data_with_proto_types_with_dml]
