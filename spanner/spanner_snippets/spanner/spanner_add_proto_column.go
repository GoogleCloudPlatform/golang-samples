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

// [START spanner_add_proto_type_columns]

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
	adminpb "cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
)

func addProtoColumn(ctx context.Context, w io.Writer, db string) error {
	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return err
	}
	defer adminClient.Close()

	file, err := os.Open(filepath.Join("testdata", "protos", "descriptors.pb"))
	if err != nil {
		return err
	}
	defer file.Close()

	protoFileDescriptor, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	op, err := adminClient.UpdateDatabaseDdl(ctx, &adminpb.UpdateDatabaseDdlRequest{
		Database: db,
		Statements: []string{
			"CREATE PROTO BUNDLE (" +
				"examples.spanner.music.SingerInfo," +
				"examples.spanner.music.Genre," +
				")",
			"ALTER TABLE Singers ADD COLUMN SingerInfo examples.spanner.music.SingerInfo",
			"ALTER TABLE Singers ADD COLUMN SingerInfoArray ARRAY<examples.spanner.music.SingerInfo>",
			"ALTER TABLE Singers ADD COLUMN SingerGenre examples.spanner.music.Genre",
			"ALTER TABLE Singers ADD COLUMN SingerGenreArray ARRAY<examples.spanner.music.Genre>",
		},
		ProtoDescriptors: protoFileDescriptor,
	})
	if err != nil {
		return err
	}
	if err := op.Wait(ctx); err != nil {
		return err
	}
	fmt.Fprintf(w, "Added Proto columns\n")
	return nil
}

// [END spanner_add_proto_type_columns]
