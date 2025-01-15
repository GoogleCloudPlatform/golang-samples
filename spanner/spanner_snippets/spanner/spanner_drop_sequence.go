// Copyright 2023 Google LLC
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

// [START spanner_drop_sequence]
// [START spanner_postgresql_drop_sequence]

import (
	"context"
	"fmt"
	"io"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
	adminpb "cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
)

func dropSequence(w io.Writer, db string) error {
	// db := "projects/my-project/instances/my-instance/databases/my-database"
	ctx := context.Background()
	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return err
	}
	defer adminClient.Close()

	// List of DDL statements to be applied to the database.
	// Drop the DEFAULT from CustomerId column and drop the sequence.
	ddl := []string{
		"ALTER TABLE Customers ALTER COLUMN CustomerId DROP DEFAULT",
		"DROP SEQUENCE Seq",
	}
	op, err := adminClient.UpdateDatabaseDdl(ctx, &adminpb.UpdateDatabaseDdlRequest{
		Database:   db,
		Statements: ddl,
	})
	if err != nil {
		return err
	}
	// Wait for the UpdateDatabaseDdl operation to finish.
	if err := op.Wait(ctx); err != nil {
		return fmt.Errorf("waiting for bit reverse sequence drop to finish failed: %w", err)
	}
	fmt.Fprintf(w, "Altered Customers table to drop DEFAULT from CustomerId column and dropped the Seq sequence\n")
	return nil
}

// [END spanner_drop_sequence]
// [END spanner_postgresql_drop_sequence]
