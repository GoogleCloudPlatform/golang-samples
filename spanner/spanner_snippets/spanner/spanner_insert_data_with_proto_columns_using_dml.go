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

// [START spanner_insert_data_with_proto_columns_using_dml]

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/spanner"
	pb "github.com/GoogleCloudPlatform/golang-samples/spanner/spanner_snippets/spanner/testdata/protos"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/proto"
)

func insertDataWithProtoMsgAndEnumUsingDML(w io.Writer, db string) error {
	ctx := context.Background()
	endpoint := "staging-wrenchworks.sandbox.googleapis.com:443"
	client, err := spanner.NewClient(ctx, db, option.WithEndpoint(endpoint))
	if err != nil {
		return err
	}
	defer client.Close()

	// Using Protocol Buffers: https://developers.google.com/protocol-buffers/docs/gotutorial
	// Creating instance of SingerInfo and Genre from user-defined Proto Message and Enum
	singer5ProtoEnum := pb.Genre_POP
	singer5ProtoMsg := &pb.SingerInfo{
		SingerId:    proto.Int64(5),
		BirthDate:   proto.String("May"),
		Nationality: proto.String("Country5"),
		Genre:       &singer5ProtoEnum,
	}
	_, err = client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		stmt := spanner.Statement{
			SQL: `INSERT INTO Singers (SingerId, FirstName, LastName, SingerInfo, SingerGenre) 
                   VALUES (@singerId, @firstName, @lastName, @singerInfo, @singerGenre)`,
			Params: map[string]interface{}{
				"singerId":    5,
				"firstName":   "Singer5",
				"lastName":    "Singer5",
				"singerInfo":  singer5ProtoMsg,
				"singerGenre": singer5ProtoEnum,
			},
		}

		rowCount, err := txn.Update(ctx, stmt)
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "%d record(s) inserted.\n", rowCount)
		return err
	})
	return err
}

// [END spanner_insert_data_with_proto_columns_using_dml]
