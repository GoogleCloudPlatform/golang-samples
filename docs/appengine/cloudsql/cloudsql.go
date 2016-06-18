// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Sample cloudsql demonstrates connection to a Cloud SQL instance from App Engine standard.
package cloudsql

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	projectID := os.Getenv("CLOUDSQL_PROJECT_ID")
	instanceName := os.Getenv("CLOUDSQL_INSTANCE")

	if projectID == "" {
		http.Error(w, "Missing project ID environment variable.", 500)
		return
	}
	if instanceName == "" {
		http.Error(w, "Missing instance name environment variable.", 500)
		return
	}

	w.Header().Set("Content-Type", "text/plain")

	db, err := sql.Open("mysql", fmt.Sprintf("root@cloudsql(%s:%s)/", projectID, instanceName))
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not open db: %v", err), 500)
		return
	}
	defer db.Close()

	rows, err := db.Query("SHOW DATABASES")
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not query db: %v", err), 500)
		return
	}
	defer rows.Close()

	buf := bytes.NewBufferString("Databases:\n")
	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			http.Error(w, fmt.Sprintf("Could not scan result: %v", err), 500)
			return
		}
		fmt.Fprintf(buf, "- %s\n", dbName)
	}
	w.Write(buf.Bytes())
}
