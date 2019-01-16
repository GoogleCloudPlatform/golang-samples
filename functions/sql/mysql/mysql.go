// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START functions_sql_mysql]

// Package sql contains examples of using to Cloud SQL.
package sql

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	// Import the MySQL SQL driver.
	_ "github.com/go-sql-driver/mysql"
)

var (
	db *sql.DB

	connectionName = os.Getenv("MYSQL_INSTANCE_CONNECTION_NAME")
	dbUser         = os.Getenv("MYSQL_USER")
	dbPassword     = os.Getenv("MYSQL_PASSWORD")
	dsn            = fmt.Sprintf("%s:%s@unix(/cloudsql/%s)/", dbUser, dbPassword, connectionName)
)

func init() {
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Could not open db: %v", err)
	}

	// Only allow 1 connection to the database to avoid overloading it.
	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(1)
}

// MySQLDemo is an example of making a MySQL database query.
func MySQLDemo(w http.ResponseWriter, r *http.Request) {
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

// [END functions_sql_mysql]
