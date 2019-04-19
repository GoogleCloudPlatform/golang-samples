// Copyright 2019 Google LLC
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

// [START functions_sql_postgres]

// Package sql contains examples of using to Cloud SQL.
package sql

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	// Import the Postgres SQL driver.
	_ "github.com/lib/pq"
)

var (
	db *sql.DB

	connectionName = os.Getenv("POSTGRES_INSTANCE_CONNECTION_NAME")
	dbUser         = os.Getenv("POSTGRES_USER")
	dbPassword     = os.Getenv("POSTGRES_PASSWORD")
	dsn            = fmt.Sprintf("user=%s password=%s host=/cloudsql/%s", dbUser, dbPassword, connectionName)
)

func init() {
	var err error
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Could not open db: %v", err)
	}

	// Only allow 1 connection to the database to avoid overloading it.
	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(1)
}

// PostgresDemo is an example of making a Postgres database query.
func PostgresDemo(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT NOW() as now")
	if err != nil {
		log.Printf("db.Query: %v", err)
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	now := ""
	rows.Next()
	if err := rows.Scan(&now); err != nil {
		log.Printf("rows.Scan: %v", err)
		http.Error(w, "Error scanning database", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Now: %v", now)
}

// [END functions_sql_postgres]
