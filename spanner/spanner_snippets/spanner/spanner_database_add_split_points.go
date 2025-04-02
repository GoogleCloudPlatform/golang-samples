// Copyright 2025 Google LLC
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

// [START spanner_database_add_split_points]

import (
	"context"
	"fmt"
	"io"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
	"cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
	adminpb "cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
	"google.golang.org/protobuf/types/known/structpb"
)

func addSplitpoints(w io.Writer, dbName string) error {
	ctx := context.Background()

	client, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	// Singers table is already created int the createDatabase test
	ddl := []string{
		"CREATE INDEX IF NOT EXISTS SingersByFirstLastName ON Singers(FirstName, LastName)",
	}
	op, err := client.UpdateDatabaseDdl(ctx, &adminpb.UpdateDatabaseDdlRequest{
		Database:   dbName,
		Statements: ddl,
	})

	// Wait for the UpdateDatabaseDdl operation to finish.
	if err := op.Wait(ctx); err != nil {
		return fmt.Errorf("waiting for UpdateDatabaseDdlRequest to finish failed: %w", err)
	}
	fmt.Fprintf(w, "Added indexes for Split testing\n")

	splitTableKey := databasepb.SplitPoints_Key{
		KeyParts: &structpb.ListValue{
			Values: []*structpb.Value{
				structpb.NewStringValue("42"),
			},
		},
	}

	splitsKeysArrayForTable := [...]databasepb.SplitPoints_Key{splitTableKey}

	splitForTable := databasepb.SplitPoints{
		Table: "Singers",
		Keys:  &splitsKeysArrayForTable,
	}

	splitIndexKey := databasepb.SplitPoints_Key{
		KeyParts: &structpb.ListValue{
			Values: []*structpb.Value{
				structpb.NewStringValue("John"),
				structpb.NewStringValue("Doe"),
			},
		},
	}

	splitsKeysArrayForIndex := [...]databasepb.SplitPoints_Key{splitIndexKey}

	splitForindex := databasepb.SplitPoints{
		Index: "SingersByFirstLastName",
		Keys:  &splitsKeysArrayForIndex,
	}

	splitsArray := [...]databasepb.SplitPoints{splitForTable, splitForindex}
	// Add split points to table and index
	req := databasepb.AddSplitPointsRequest{
		Database:    dbName,
		SplitPoints: &splitsArray,
	}

	res, err := client.AddSplitPoints(ctx, &req)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "Added split points %s", res)
	return nil
}

// [END spanner_database_add_split_points]
