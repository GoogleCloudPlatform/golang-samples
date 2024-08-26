// Copyright 2020 Google LLC
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

// [START spanner_list_database_operations]

import (
	"context"
	"fmt"
	"io"
	"regexp"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
	adminpb "cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/api/iterator"
)

func listDatabaseOperations(ctx context.Context, w io.Writer, db string) error {
	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return err
	}
	defer adminClient.Close()

	matches := regexp.MustCompile("^(.*)/databases/(.*)$").FindStringSubmatch(db)
	if matches == nil || len(matches) != 3 {
		return fmt.Errorf("Invalid database id %s", db)
	}
	instanceName := matches[1]
	// List the databases that are being optimized after a restore operation.
	filter := "(metadata.@type:type.googleapis.com/google.spanner.admin.database.v1.OptimizeRestoredDatabaseMetadata)"
	iter := adminClient.ListDatabaseOperations(ctx, &adminpb.ListDatabaseOperationsRequest{
		Parent: instanceName,
		Filter: filter,
	})
	for {
		resp, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		metadata := &adminpb.OptimizeRestoredDatabaseMetadata{}
		if err := ptypes.UnmarshalAny(resp.Metadata, metadata); err != nil {
			return err
		}
		fmt.Fprintf(w, "Database %s restored from backup is %d%% optimized.\n",
			metadata.Name,
			metadata.Progress.ProgressPercent,
		)
	}
	return nil
}

// [END spanner_list_database_operations]
