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
	pb "github.com/GoogleCloudPlatform/golang-samples/spanner/spanner_snippets/spanner/testdata/protos"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/proto"
)

// [START spanner_insert_data_with_proto_columns]

func insertDataWithProtoMsgAndEnum(w io.Writer, db string) error {
	ctx := context.Background()
	endpoint := "staging-wrenchworks.sandbox.googleapis.com:443"
	client, err := spanner.NewClient(ctx, db, option.WithEndpoint(endpoint))
	if err != nil {
		return err
	}
	defer client.Close()

	// Using Protocol Buffers: https://developers.google.com/protocol-buffers/docs/gotutorial
	// Creating instance of SingerInfo and Genre from user-defined Proto Message and Enum
	singer1ProtoEnum := pb.Genre_ROCK
	singer1ProtoMsg := &pb.SingerInfo{
		SingerId:    proto.Int64(1),
		BirthDate:   proto.String("January"),
		Nationality: proto.String("Country1"),
		Genre:       &singer1ProtoEnum,
	}

	singer2ProtoEnum := pb.Genre_FOLK
	singer2ProtoMsg := &pb.SingerInfo{
		SingerId:    proto.Int64(2),
		BirthDate:   proto.String("February"),
		Nationality: proto.String("Country2"),
		Genre:       &singer2ProtoEnum,
	}

	singer3ProtoEnum := pb.Genre_JAZZ
	singer3ProtoMsg := &pb.SingerInfo{
		SingerId:    proto.Int64(3),
		BirthDate:   proto.String("March"),
		Nationality: proto.String("Country3"),
		Genre:       &singer3ProtoEnum,
	}

	singer4ProtoEnum := pb.Genre_POP
	singer4ProtoMsg := &pb.SingerInfo{
		SingerId:    proto.Int64(4),
		BirthDate:   proto.String("April"),
		Nationality: proto.String("Country4"),
		Genre:       &singer4ProtoEnum,
	}

	cols := []string{"SingerId", "FirstName", "LastName", "SingerInfo", "SingerGenre"}
	m := []*spanner.Mutation{
		spanner.InsertOrUpdate("Singers", cols, []interface{}{1, "Singer1", "Singer1", singer1ProtoMsg, &singer1ProtoEnum}),
		spanner.InsertOrUpdate("Singers", cols, []interface{}{2, "Singer2", "Singer2", singer2ProtoMsg, &singer2ProtoEnum}),
		spanner.InsertOrUpdate("Singers", cols, []interface{}{3, "Singer3", "Singer3", singer3ProtoMsg, &singer3ProtoEnum}),
		spanner.InsertOrUpdate("Singers", cols, []interface{}{4, "Singer4", "Singer4", singer4ProtoMsg, &singer4ProtoEnum}),
	}
	_, err = client.Apply(ctx, m)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Inserted data to SingerInfo and SingerGenre columns")
	return nil
}

// [END spanner_insert_data_with_proto_columns]
