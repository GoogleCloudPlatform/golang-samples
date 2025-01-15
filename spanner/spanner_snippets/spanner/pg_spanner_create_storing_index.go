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

// [START spanner_postgresql_create_storing_index]

import (
	"context"
	"fmt"
	"io"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
	adminpb "cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
)

// pgAddStoringIndex shows how to create 'STORING' indexes on a Spanner
// PostgreSQL database. The PostgreSQL dialect uses INCLUDE keyword, as
// opposed to the STORING keyword of Cloud Spanner.
func pgAddStoringIndex(ctx context.Context, w io.Writer, db string) error {
	// db := "projects/my-project/instances/my-instance/databases/my-database"
	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize spanner database admin client: %w", err)
	}
	defer adminClient.Close()

	op, err := adminClient.UpdateDatabaseDdl(ctx, &adminpb.UpdateDatabaseDdlRequest{
		Database: db,
		Statements: []string{
			"CREATE INDEX AlbumsByAlbumTitle2 ON Albums(AlbumTitle) INCLUDE (MarketingBudget)",
		},
	})
	if err != nil {
		return fmt.Errorf("failed to execute spanner database DDL request: %w", err)
	}
	if err := op.Wait(ctx); err != nil {
		return fmt.Errorf("failed to complete spanner database DDL request: %w", err)
	}
	fmt.Fprintf(w, "Added storing index\n")
	return nil
}

// [END spanner_postgresql_create_storing_index]
