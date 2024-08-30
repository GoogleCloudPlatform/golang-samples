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

// [START spanner_create_database_with_property_graph]

import (
	"context"
	"fmt"
	"io"
	"regexp"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
	adminpb "google.golang.org/genproto/googleapis/spanner/admin/database/v1"
)

func makeCreateDatabaseWithPropertyGraphRequest(instance string, dbName string) *adminpb.CreateDatabaseRequest {
	// The schema defintion for a database with a property graph comprises table
	// definitions one or more `CREATE PROPERTY GRAPH` statements to define the
	// property graph(s).
	// Consult https://cloud.google.com/spanner/docs/reference/standard-sql/graph-schema-statements
	// for a description of property graph schemas.
	var schema_statements = []string {
		`CREATE TABLE Person (
			id               INT64 NOT NULL,
			name             STRING(MAX),
			birthday         TIMESTAMP,
			country          STRING(MAX),
			city             STRING(MAX),
		) PRIMARY KEY (id)`,
		`CREATE TABLE Account (
			id               INT64 NOT NULL,
			create_time      TIMESTAMP,
			is_blocked       BOOL,
			nick_name        STRING(MAX),
		) PRIMARY KEY (id)`,
		`CREATE TABLE PersonOwnAccount (
			id               INT64 NOT NULL,
			account_id       INT64 NOT NULL,
			create_time      TIMESTAMP,
			FOREIGN KEY (account_id)
				REFERENCES Account (id)
		) PRIMARY KEY (id, account_id),
		INTERLEAVE IN PARENT Person ON DELETE CASCADE`,
		`CREATE TABLE AccountTransferAccount (
			id               INT64 NOT NULL,
			to_id            INT64 NOT NULL,
			amount           FLOAT64,
			create_time      TIMESTAMP NOT NULL,
			order_number     STRING(MAX),
			FOREIGN KEY (to_id) REFERENCES Account (id)
		) PRIMARY KEY (id, to_id, create_time),
		INTERLEAVE IN PARENT Account ON DELETE CASCADE`,
		`CREATE OR REPLACE PROPERTY GRAPH FinGraph
			NODE TABLES (Account, Person)
			EDGE TABLES (
				PersonOwnAccount
					SOURCE KEY(id) REFERENCES Person(id)
					DESTINATION KEY(account_id) REFERENCES Account(id)
					LABEL Owns,
				AccountTransferAccount
					SOURCE KEY(id) REFERENCES Account(id)
					DESTINATION KEY(to_id) REFERENCES Account(id)
					LABEL Transfers)`,
		};

	return &adminpb.CreateDatabaseRequest{
		Parent:          instance,
		CreateStatement: "CREATE DATABASE `" + dbName + "`",
		ExtraStatements: schema_statements,
	}
}

func createDatabaseWithPropertyGraph(ctx context.Context, w io.Writer, dbId string) error {
	// dbId is of the form:
	// 	projects/<project>/instances/<instance>/databases/<database>
	matches := regexp.MustCompile("^(.*)/databases/(.*)$").FindStringSubmatch(dbId)
	if matches == nil || len(matches) != 3 {
		return fmt.Errorf("Invalid database id %s", dbId)
	}

	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return err
	}
	defer adminClient.Close()

	var instance = matches[1]
	var dbName = matches[2]

	op, err := adminClient.CreateDatabase(ctx, makeCreateDatabaseWithPropertyGraphRequest(instance, dbName))
	if err != nil {
		return err
	}
	if _, err := op.Wait(ctx); err != nil {
		return err
	}
	fmt.Fprintf(w, "Created database [%s]\n", dbId)
	return nil
}

// [END spanner_create_database_with_property_graph]
