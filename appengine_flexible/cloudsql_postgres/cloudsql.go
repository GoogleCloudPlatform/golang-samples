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

// +build go1.8

// [START gae_flex_postgres_app]

// Sample cloudsql demonstrates usage of Cloud SQL from App Engine flexible environment.
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	// Set this in app.yaml when running in production.
	datastoreName := os.Getenv("POSTGRES_CONNECTION")

	var err error
	db, err = sql.Open("postgres", datastoreName)
	if err != nil {
		log.Fatal(err)
	}

	// Ensure the table exists.
	// Running an SQL query also checks the connection to the PostgreSQL server
	// is authenticated and valid.
	if err := createTable(); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", handle)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
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

func maskIP(addr string) (string, error) {
	// Discard the port portion of the address.
	userIP, _, err := net.SplitHostPort(addr)
	if err != nil {
		return "", err
	}

	// Parse the IP and mask a portion of it.
	ip := net.ParseIP(userIP)
	if ip.To4() != nil {
		ipV4Mask := net.CIDRMask(16, 32)
		ip = ip.Mask(ipV4Mask)
	} else {
		ipv6Mask := net.CIDRMask(32, 128)
		ip = ip.Mask(ipv6Mask)
	}
	return ip.String(), nil
}

// recordVisit records the masked IP of the visit.
func recordVisit(timestamp int64, addr string) error {
	userIP, err := maskIP(addr)
	if err != nil {
		return err
	}
	stmt := "INSERT INTO visits (timestamp, userip) VALUES ($1, $2)"
	_, err = db.Exec(stmt, timestamp, userIP)
	return err
}

func queryVisits(limit int64) ([]visit, error) {
	rows, err := db.Query("SELECT timestamp, userip FROM visits ORDER BY timestamp DESC LIMIT $1", limit)
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

// [END gae_flex_postgres_app]
