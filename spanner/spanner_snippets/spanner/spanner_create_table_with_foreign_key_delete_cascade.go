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

// [START spanner_create_table_with_foreign_key_delete_cascade]

import (
	"context"
	"fmt"
	"io"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
	adminpb "cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
)

func createTableWithForeignKeyDeleteCascade(w io.Writer, db string) error {
	// db := "projects/my-project/instances/my-instance/databases/my-database"
	ctx := context.Background()
	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return err
	}
	defer adminClient.Close()

	// List of DDL statements to be applied to the database.
	// Create a parent table, and then use the primary key of parent table as auto foreign key in ShoppingCarts table.
	ddl := []string{
		"CREATE TABLE Customers ( CustomerId INT64 NOT NULL, CustomerName STRING(62) NOT NULL, ) PRIMARY KEY (CustomerId)",
		"CREATE TABLE ShoppingCarts ( CartId INT64 NOT NULL, CustomerId INT64 NOT NULL, CustomerName STRING(62) NOT NULL, CONSTRAINT FKShoppingCartsCustomerId FOREIGN KEY (CustomerId) REFERENCES Customers (CustomerId) ON DELETE CASCADE ) PRIMARY KEY (CartId)",
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
	fmt.Fprintf(w, "Created Customers and ShoppingCarts table with FKShoppingCartsCustomerId foreign key constraint on database %v.\n", db)
	return err
}

// [END spanner_create_table_with_foreign_key_delete_cascade]
