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

import (
	"context"
	"fmt"
	"io"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
	adminpb "cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
)

func createTableDocumentsWithHistoryTable(ctx context.Context, w io.Writer, db string) error {
	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return err
	}
	defer adminClient.Close()

	op, err := adminClient.UpdateDatabaseDdl(ctx, &adminpb.UpdateDatabaseDdlRequest{
		Database: db,
		Statements: []string{
			`CREATE TABLE Documents(
				UserId INT64 NOT NULL,
				DocumentId INT64 NOT NULL,
				Contents STRING(MAX) NOT NULL
			) PRIMARY KEY(UserId, DocumentId)`,
			`CREATE TABLE DocumentHistory(
				UserId INT64 NOT NULL,
				DocumentId INT64 NOT NULL,
				Timestamp TIMESTAMP NOT NULL OPTIONS(allow_commit_timestamp=true),
				PreviousContents STRING(MAX)
			) PRIMARY KEY(UserId, DocumentId, Timestamp), INTERLEAVE IN PARENT Documents ON DELETE NO ACTION`,
		},
	})
	if err != nil {
		return err
	}
	if err := op.Wait(ctx); err != nil {
		return err
	}
	fmt.Fprintf(w, "Created Documents and DocumentHistory tables in database [%s]\n", db)
	return nil
}
