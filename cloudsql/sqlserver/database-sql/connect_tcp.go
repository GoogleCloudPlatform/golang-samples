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

// [START cloud_sql_sqlserver_databasesql_connect_tcp]
package cloudsql

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
)

// connectTCPSocket initializes a TCP connection pool for a Cloud SQL
// instance of SQL Server.
func connectTCPSocket() (*sql.DB, error) {
	mustGetenv := func(k string) string {
		v := os.Getenv(k)
		if v == "" {
			log.Fatalf("Warning: %s environment variable not set.", k)
		}
		return v
	}

	var (
		dbUser    = mustGetenv("DB_USER")       // e.g. 'my-db-user'
		dbPwd     = mustGetenv("DB_PASS")       // e.g. 'my-db-password'
		dbTCPHost = mustGetenv("INSTANCE_HOST") // e.g. '127.0.0.1' ('172.17.0.1' if deployed to GAE Flex)
		dbPort    = mustGetenv("DB_PORT")       // e.g. '1433'
		dbName    = mustGetenv("DB_NAME")       // e.g. 'my-database'
	)

	dbURI := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;",
		dbTCPHost, dbUser, dbPwd, dbPort, dbName)

	// [START_EXCLUDE]
	// [START cloud_sql_sqlserver_databasesql_sslcerts]
	// (OPTIONAL) Configure SSL certificates
	// For deployments that connect directly to a Cloud SQL instance without
	// using the Cloud SQL Proxy, configuring SSL certificates will ensure the
	// connection is encrypted. This step is entirely OPTIONAL.
	dbRootCert := os.Getenv("DB_ROOT_CERT") // e.g., '/path/to/my/server-ca.pem'
	if dbRootCert != "" {
		// Get connection host name to be matched with host name in SSL certificate.
		var instanceConnectionName = mustGetenv("INSTANCE_CONNECTION_NAME")
		// Expected format of INSTANCE_CONNECTION_NAME is project_id:region:instance_name
		hostNameParts := strings.Split(instanceConnectionName, ":")
		if len(hostNameParts) != 3 {
			log.Fatalf("Invalid format for INSTANCE_CONNECTION_NAME environment variable, got = %q", instanceConnectionName)
		}
		// Specify encrypted connection, host name and certificate.
		dbURI += fmt.Sprintf("encrypt=true;hostnameincertificate=%s:%s;certificate=%s;",
			hostNameParts[0], hostNameParts[2], dbRootCert)
	}
	// [END cloud_sql_sqlserver_databasesql_sslcerts]
	// [END_EXCLUDE]

	// dbPool is the pool of database connections.
	dbPool, err := sql.Open("sqlserver", dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %v", err)
	}

	// [START_EXCLUDE]
	configureConnectionPool(dbPool)
	// [END_EXCLUDE]

	return dbPool, nil
}

// [END cloud_sql_sqlserver_databasesql_connect_tcp]
