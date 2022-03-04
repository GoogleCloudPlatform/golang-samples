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

// [START cloud_sql_sqlserver_databasesql_connect_connector]
package cloudsql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	"cloud.google.com/go/cloudsqlconn"
	mssql "github.com/denisenkom/go-mssqldb"
)

type csqlDialer struct {
	dialer     *cloudsqlconn.Dialer
	connName   string
	usePrivate bool
}

// DialContext adheres to the mssql.Dialer interface.
func (c *csqlDialer) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	var opts []cloudsqlconn.DialOption
	if c.usePrivate {
		opts = append(opts, cloudsqlconn.WithPrivateIP())
	}
	return c.dialer.Dial(ctx, c.connName, opts...)
}

func connectWithConnector() (*sql.DB, error) {
	mustGetenv := func(k string) string {
		v := os.Getenv(k)
		if v == "" {
			log.Fatalf("Warning: %s environment variable not set.\n", k)
		}
		return v
	}

	var (
		dbUser                 = mustGetenv("DB_USER")                  // e.g. 'my-db-user'
		dbPwd                  = mustGetenv("DB_PASS")                  // e.g. 'my-db-password'
		dbName                 = mustGetenv("DB_NAME")                  // e.g. 'my-database'
		instanceConnectionName = mustGetenv("INSTANCE_CONNECTION_NAME") // e.g. 'project:region:instance'
		usePrivate             = os.Getenv("PRIVATE_IP")
	)

	dbURI := fmt.Sprintf("user id=%s;password=%s;database=%s;", dbUser, dbPwd, dbName)
	c, err := mssql.NewConnector(dbURI)
	if err != nil {
		return nil, fmt.Errorf("mssql.NewConnector: %v", err)
	}
	dialer, err := cloudsqlconn.NewDialer(context.Background())
	if err != nil {
		return nil, fmt.Errorf("cloudsqlconn.NewDailer: %v", err)
	}
	c.Dialer = &csqlDialer{
		dialer:     dialer,
		connName:   instanceConnectionName,
		usePrivate: usePrivate != "",
	}

	dbPool := sql.OpenDB(c)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %v", err)
	}
	return dbPool, nil
}

// [END cloud_sql_sqlserver_databasesql_connect_connector]
