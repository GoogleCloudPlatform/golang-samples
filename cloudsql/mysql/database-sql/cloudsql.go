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

// [START gae_cloudsql]

// This sample demonstrates connection to a Cloud SQL instance from App Engine
// standard. The application is a Golang version of the "Tabs vs Spaces" web
// app presented at Cloud Next '19 as seen in this video:
// https://www.youtube.com/watch?v=qVgzP3PsXFw&t=1833s
package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"os"
	"runtime"
	"strconv"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

// Global db variable holds the database connection.
var db *sql.DB

// The vote struct stores a row from the votes table in the Cloud SQL instance.
// Each vote includes a candidate ("TABS" or "SPACES") and a timestamp.
type vote struct {
	Candidate string
	VoteTime  mysql.NullTime
}

// The templateData struct is used to pass data to the HTML template.
type templateData struct {
	TabsCount   uint
	SpacesCount uint
	VoteMargin  string
	RecentVotes []vote
}

func main() {
	var err error

	db, err = initConnectionPool()
	if err != nil {
		log.Fatalf("unable to initialize database connection pool: %s", err)
	}

	if err = initDBSchema(); err != nil {
		log.Fatalf("unable to create table: %s", err)
	}

	http.HandleFunc("/", indexHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

// indexHandler handles requests to the / route.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		showTotals(w, r)
	case "POST":
		saveVote(w, r)
	default:
		fmt.Fprintf(w, "Unsupported HTTP verb: %s\n", r.Method)
	}
}

// recentVotes returns a slice of the last 5 votes cast.
func recentVotes() []vote {
	var votes []vote
	rows, err := db.Query(`SELECT candidate, time_cast FROM votes ORDER BY time_cast DESC LIMIT 5`)
	if err != nil {
		fmt.Println(err)
		return votes
	}
	defer rows.Close()
	for rows.Next() {
		nextVote := vote{}
		if err := rows.Scan(&nextVote.Candidate, &nextVote.VoteTime); err != nil {
			log.Fatalf("unable to scan row returned by SELECT statement: %s", err)
		}
		votes = append(votes, nextVote)
	}
	if err = rows.Err(); err != nil {
		log.Fatalf("error reading selected rows: %s", err)
	}
	return votes
}

// currentTotals returns a templateData structure for populating the web page.
func currentTotals() templateData {

	// get total votes for each candidate
	var tabVotes, spaceVotes uint
	_ = db.QueryRow(`select count(vote_id) from votes where candidate='TABS'`).Scan(&tabVotes)
	_ = db.QueryRow(`select count(vote_id) from votes where candidate='SPACES'`).Scan(&spaceVotes)

	// voteMargin is string representation of current voting margin.
	voteDiff := int(math.Abs(float64(tabVotes) - float64(spaceVotes)))
	var voteMargin string
	if voteDiff == 1 {
		voteMargin = "1 vote"
	}
	voteMargin = strconv.Itoa(voteDiff) + " votes"

	return templateData{tabVotes, spaceVotes, voteMargin, recentVotes()}

}

// showTotals renders an HTML template showing the current vote totals.
func showTotals(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatalf("unable to parse template file: %s", err)
	}

	err = t.Execute(w, currentTotals())
	if err != nil {
		log.Fatalf("unable to execute template: %s", err)
	}
}

// saveVote handles POST requests and saves a vote passed as http.Request form data.
func saveVote(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	team := r.Form["team"][0]
	// [START cloud_sql_mysql_database-sql_connection]
	sqlInsert := "INSERT INTO votes (candidate) VALUES (?)"
	if team == "TABS" || team == "SPACES" {
		if _, err := db.Exec(sqlInsert, team); err != nil {
			log.Fatalf("unable to save vote: %s", err)
		} else {
			fmt.Fprintf(w, "Vote successfully cast for %s!\n", team)
		}
	}
	// [END cloud_sql_mysql_database-sql_connection]
}

// mustGetEnv is a helper function for getting environment variables.
// Displays a warning if the environment variable is not set.
func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		fmt.Printf("Warning: %s environment variable not set.\n", k)
	}
	return v
}

// [START cloud_sql_mysql_database-sql_create]

// initConnectionPool initializes a database connection for a Cloud SQL instance.
func initConnectionPool() (*sql.DB, error) {
	var (
		dbUser                 = mustGetenv("DB_USER")
		dbPass                 = mustGetenv("DB_PASS")
		instanceConnectionName = mustGetenv("INSTANCE_CONNECTION_NAME")
		dbName                 = mustGetenv("DB_NAME")
	)

	// Default connection string format for App Engine deployment (UNIX socket syntax).
	dbURI := fmt.Sprintf("%s:%s@unix(/cloudsql/%s)/%s", dbUser, dbPass, instanceConnectionName, dbName)
	if runtime.GOOS == "windows" {
		// If running on Windows (local dev machine), connect via Cloud SQL Proxy.
		instanceConnectionName = "127.0.0.1:3306"
		dbURI = fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPass, instanceConnectionName, dbName)
	}

	dbConn, err := sql.Open("mysql", dbURI)
	if err != nil {
		return nil, err
	}

	// [START cloud_sql_mysql_database-sql_limit]
	dbConn.SetMaxIdleConns(5)
	dbConn.SetMaxOpenConns(7)
	// [END cloud_sql_mysql_database-sql_limit]

	// [START cloud_sql_mysql_database-sql_lifetime]
	dbConn.SetConnMaxLifetime(1800)
	// [END cloud_sql_mysql_database-sql_lifetime]

	return dbConn, nil
}

// [END cloud_sql_mysql_database-sql_create]

// initDBSchema creates the votes table if it does not exist.
func initDBSchema() error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS votes
	( vote_id SERIAL NOT NULL, time_cast timestamp NOT NULL,
	candidate CHAR(6) NOT NULL, PRIMARY KEY (vote_id) );`)
	return err
}

// [END gae_cloudsql]
