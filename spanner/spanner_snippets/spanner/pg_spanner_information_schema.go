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

// [START spanner_postgresql_information_schema]

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/spanner"
	database "cloud.google.com/go/spanner/admin/database/apiv1"
	"google.golang.org/api/iterator"
	adminpb "google.golang.org/genproto/googleapis/spanner/admin/database/v1"
)

// pgInformationSchema shows how to query the information schema metadata in a
// Spanner PostgreSQL database.
func pgInformationSchema(w io.Writer, db string) error {
	// db := "projects/my-project/instances/my-instance/databases/my-database"
	ctx := context.Background()

	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return err
	}
	defer adminClient.Close()

	// Create a table, and then get the metadata of the table from the INFORMATION_SCHEMA.
	op, err := adminClient.UpdateDatabaseDdl(ctx, &adminpb.UpdateDatabaseDdlRequest{
		Database: db,
		Statements: []string{
			`CREATE TABLE Venues (
				VenueId  bigint NOT NULL PRIMARY KEY,
				Name     varchar(1024) NOT NULL,
				Revenues numeric,
				Picture  bytea
			)`},
	})
	if err != nil {
		return err
	}
	if err := op.Wait(ctx); err != nil {
		return err
	}

	// The Spanner INFORMATION_SCHEMA tables can be used to query the metadata of tables and
	// columns of PostgreSQL databases. The returned results will include additional PostgreSQL
	// metadata columns.

	// Get all the user tables in the database. PostgreSQL uses the `public` schema for user
	// tables. The table_catalog is equal to the database name.

	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	// The `user_defined_...` columns are only available for PostgreSQL databases.
	iter := client.Single().Query(ctx, spanner.Statement{SQL: `SELECT table_catalog, table_schema, table_name, 
				user_defined_type_catalog, 
				user_defined_type_schema, 
				user_defined_type_name 
		FROM INFORMATION_SCHEMA.tables 
		WHERE table_schema='public'`,
	})
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		var tableCatalog, tableSchema, tableName string
		var typeCatalog, typeSchema, typeName spanner.NullString
		if err := row.Columns(&tableCatalog, &tableSchema, &tableName, &typeCatalog, &typeSchema, &typeName); err != nil {
			return err
		}
		userDefinedType := "null"
		if typeCatalog.Valid {
			userDefinedType = fmt.Sprintf("%s.%s.%s", typeCatalog, typeSchema, typeName)
		}
		fmt.Fprintf(w, "Table: %s.%s.%s (User defined type: %s)\n", tableCatalog, tableSchema, tableName, userDefinedType)
	}

	return nil
}

// [END spanner_postgresql_information_schema]
