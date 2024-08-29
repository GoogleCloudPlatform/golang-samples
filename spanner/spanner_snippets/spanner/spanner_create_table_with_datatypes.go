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

// [START spanner_create_table_with_datatypes]

import (
	"context"
	"fmt"
	"io"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
	adminpb "cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
)

// Creates a Cloud Spanner table comprised of columns for each supported data type
// See https://cloud.google.com/spanner/docs/data-types
func createTableWithDatatypes(ctx context.Context, w io.Writer, db string) error {
	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return err
	}
	defer adminClient.Close()

	op, err := adminClient.UpdateDatabaseDdl(ctx, &adminpb.UpdateDatabaseDdlRequest{
		Database: db,
		Statements: []string{
			`CREATE TABLE Venues (
				VenueId	INT64 NOT NULL,
				VenueName STRING(100),
				VenueInfo BYTES(MAX),
				Capacity INT64,
				AvailableDates ARRAY<DATE>,
				LastContactDate DATE,
				OutdoorVenue BOOL,
				PopularityScore FLOAT64,
				Revenue NUMERIC,
				LastUpdateTime TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true)
			) PRIMARY KEY (VenueId)`,
		},
	})
	if err != nil {
		return fmt.Errorf("UpdateDatabaseDdl: %w", err)
	}
	if err := op.Wait(ctx); err != nil {
		return err
	}
	fmt.Fprintf(w, "Created Venues table in database [%s]\n", db)
	return nil
}

// [END spanner_create_table_with_datatypes]
