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

// [START spanner_drop_foreign_key_constraint_delete_cascade]

import (
	"context"
	"fmt"
	"io"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
	adminpb "cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
)

func dropForeignKeyDeleteCascade(w io.Writer, db string) error {
	// db := "projects/my-project/instances/my-instance/databases/my-database"
	ctx := context.Background()
	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return err
	}
	defer adminClient.Close()

	// List of DDL statements to be applied to the database.
	// Alter the ShoppingCarts to drop foreign key constraint on CustomerName column of parent table.
	ddl := []string{
		"ALTER TABLE ShoppingCarts DROP CONSTRAINT FKShoppingCartsCustomerName",
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
		return fmt.Errorf("waiting for operation to complete failed: %w", err)
	}
	fmt.Fprintf(w, "Altered ShoppingCarts table to drop FKShoppingCartsCustomerName foreign key constraint on database %v.\n", db)
	return err
}

// [END spanner_drop_foreign_key_constraint_delete_cascade]
