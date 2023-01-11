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

// [START cloud_sql_postgres_databasesql_connect_connector]
package cloudsql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	"cloud.google.com/go/cloudsqlconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
)

func connectWithConnector() (*sql.DB, error) {
	mustGetenv := func(k string) string {
		v := os.Getenv(k)
		if v == "" {
			log.Fatalf("Warning: %s environment variable not set.\n", k)
		}
		return v
	}
	// Note: Saving credentials in environment variables is convenient, but not
	// secure - consider a more secure solution such as
	// Cloud Secret Manager (https://cloud.google.com/secret-manager) to help
	// keep secrets safe.
	var (
		// If a password is not defined, Automatic IAM Authentication will be used
		dbUser                 = mustGetenv("DB_USER")                   // e.g. 'my-db-user' or 'sa-name@project-id.iam' for iam
		dbPwd                  = os.Getenv("DB_PASS")                   // e.g. 'my-db-password'
		dbName                 = mustGetenv("DB_NAME")                  // e.g. 'my-database'
		instanceConnectionName = mustGetenv("INSTANCE_CONNECTION_NAME") // e.g. 'project:region:instance'
		usePrivate             = os.Getenv("PRIVATE_IP")
	)

	dsn := fmt.Sprintf("user=%s database=%s", dbUser, dbName)
	if (dbPwd != "") {
		dsn += fmt.Sprintf(" password=%s", dbPwd)
	}
	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	config.DialFunc = func(ctx context.Context, network, instance string) (net.Conn, error) {
		if (dbPwd != "") {
			// [START cloud_sql_postgres_databasesql_auto_iam_authn]
			d, err := cloudsqlconn.NewDialer(ctx, cloudsqlconn.WithIAMAuthN())
			if err != nil {
				return nil, err
			}
			return d.Dial(ctx, instanceConnectionName)
			// [END cloud_sql_postgres_databasesql_auto_iam_authn]
		}
		if usePrivate != "" {
			d, err := cloudsqlconn.NewDialer(
				ctx,
				cloudsqlconn.WithDefaultDialOptions(cloudsqlconn.WithPrivateIP()),
			)
			if err != nil {
				return nil, err
			}
			return d.Dial(ctx, instanceConnectionName)
		}
		// Use the Cloud SQL connector to handle connecting to the instance.
		// This approach does *NOT* require the Cloud SQL proxy.
		d, err := cloudsqlconn.NewDialer(ctx)
		if err != nil {
			return nil, err
		}
		return d.Dial(ctx, instanceConnectionName)
	}
	dbURI := stdlib.RegisterConnConfig(config)
	dbPool, err := sql.Open("pgx", dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %v", err)
	}
	return dbPool, nil
}

// [END cloud_sql_postgres_databasesql_connect_connector]
