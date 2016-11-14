// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Sample cloudsql demonstrates usage of Cloud SQL from App Engine flexible environment.
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"google.golang.org/appengine"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	// Set this in app.yaml when running in production.
	datastoreName := os.Getenv("MYSQL_CONNECTION")

	var err error
	db, err = sql.Open("mysql", datastoreName)
	if err != nil {
		log.Fatal(err)
	}

	// Ensure the table exists.
	// Running an SQL query also checks the connection to the MySQL server
	// is authenticated and valid.
	if err := createTable(); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", handle)
	appengine.Main()
}

func createTable() error {
	stmt := `CREATE TABLE IF NOT EXISTS visits (
			timestamp  BIGINT,
			userip     VARCHAR(255)
		)`
	_, err := db.Exec(stmt)
	return err
}

func handle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Get a list of the most recent visits.
	visits, err := queryVisits(10)
	if err != nil {
		msg := fmt.Sprintf("Could not get recent visits: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	// Record this visit.
	if err := recordVisit(time.Now().UnixNano(), r.RemoteAddr); err != nil {
		msg := fmt.Sprintf("Could not save visit: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Previous visits:")
	for _, v := range visits {
		fmt.Fprintf(w, "[%s] %s\n", time.Unix(0, v.timestamp), v.userIP)
	}
	fmt.Fprintln(w, "\nSuccessfully stored an entry of the current request.")
}

type visit struct {
	timestamp int64
	userIP    string
}

func recordVisit(timestamp int64, userIP string) error {
	stmt := "INSERT INTO visits (timestamp, userip) VALUES (?, ?)"
	_, err := db.Exec(stmt, timestamp, userIP)
	return err
}

func queryVisits(limit int64) ([]visit, error) {
	rows, err := db.Query("SELECT timestamp, userip FROM visits ORDER BY timestamp DESC LIMIT ?", limit)
	if err != nil {
		return nil, fmt.Errorf("Could not get recent visits: %v", err)
	}
	defer rows.Close()

	var visits []visit

	for rows.Next() {
		var v visit
		if err := rows.Scan(&v.timestamp, &v.userIP); err != nil {
			return nil, fmt.Errorf("Could not get timestamp/user IP out of row: %v", err)
		}
		visits = append(visits, v)
	}

	return visits, rows.Err()
}
