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

// Sample database-sql demonstrates connection to a Cloud SQL instance from App Engine
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
)

// db is global database connection, parsedTemplate is parsed HTML template.
var db *sql.DB
var parsedTemplate *template.Template

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

	parsedTemplate, err = template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatalf("unable to parse template file: %s", err)
	}

	db, err = initConnectionPool()
	if err != nil {
		log.Fatalf("initConnectionPool: unable to initialize database connection pool: %s", err)
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
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}

}

// indexHandler handles requests to the / route.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := showTotals(w, r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	case "POST":
		saveVote(w, r)
	default:
		http.Error(w, fmt.Sprintf("HTTP Method %s Not Allowed", r.Method), http.StatusMethodNotAllowed)
	}
}

// recentVotes returns a slice of the last 5 votes cast.
func recentVotes() ([]vote, error) {
	var votes []vote
	rows, err := db.Query(`SELECT candidate, time_cast FROM votes ORDER BY time_cast DESC LIMIT 5`)
	if err != nil {
		return votes, err
	}
	defer rows.Close()
	for rows.Next() {
		nextVote := vote{}
		rowError := rows.Scan(&nextVote.Candidate, &nextVote.VoteTime)
		if err != nil {
			return votes, rowError
		}
		votes = append(votes, nextVote)
	}
	if err = rows.Err(); err != nil {
		return votes, err
	}
	return votes, nil
}

// currentTotals returns a templateData structure for populating the web page.
func currentTotals() (templateData, error) {

	// get total votes for each candidate
	var tabVotes, spaceVotes uint
	err := db.QueryRow(`SELECT count(vote_id) FROM votes WHERE candidate='TABS'`).Scan(&tabVotes)
	if err != nil {
		return templateData{}, err
	}
	err = db.QueryRow(`SELECT count(vote_id) FROM votes WHERE candidate='SPACES'`).Scan(&spaceVotes)
	if err != nil {
		return templateData{}, err
	}

	// voteMargin is string representation of the current voting margin,
	// such as "1 vote" (singular) or "2 votes" (plural).
	voteDiff := int(math.Abs(float64(tabVotes) - float64(spaceVotes)))
	var voteMargin string
	if voteDiff == 1 {
		voteMargin = "1 vote"
	}
	voteMargin = strconv.Itoa(voteDiff) + " votes"

	latestVotesCast, err := recentVotes()
	if err != nil {
		return templateData{}, err
	}
	return templateData{tabVotes, spaceVotes, voteMargin, latestVotesCast}, nil

}

// showTotals renders an HTML template showing the current vote totals.
func showTotals(w http.ResponseWriter, r *http.Request) error {

	totals, err := currentTotals()
	if err != nil {
		return err
	}
	err = parsedTemplate.Execute(w, totals)
	if err != nil {
		return err
	}
	return nil
}

// saveVote handles POST requests and saves a vote passed as http.Request form data.
func saveVote(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	team := r.Form["team"][0]
	// [START cloud_sql_mysql_databasesql_connection]
	sqlInsert := "INSERT INTO votes (candidate) VALUES (?)"
	if team == "TABS" || team == "SPACES" {
		if _, err := db.Exec(sqlInsert, team); err != nil {
			log.Fatalf("unable to save vote: %s", err)
		} else {
			fmt.Fprintf(w, "Vote successfully cast for %s!\n", team)
		}
	}
	// [END cloud_sql_mysql_databasesql_connection]
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

// [START cloud_sql_mysql_databasesql_create]

// initConnectionPool initializes a database connection for a Cloud SQL instance.
func initConnectionPool() (*sql.DB, error) {
	var (
		dbUser                 = mustGetenv("DB_USER")
		dbPass                 = mustGetenv("DB_PASS")
		instanceConnectionName = mustGetenv("INSTANCE_CONNECTION_NAME")
		dbName                 = mustGetenv("DB_NAME")
	)

	// connectionType can be "Unix socket" or "TCP". Unix socket is used for
	// deployment on App Engine.
	connectionType := "Unix socket"
	if runtime.GOOS == "windows" {
		// The Cloud SQL Proxy currently only supports TCP connections on
		// Windows, so must use TCP if running locally on Windows.
		connectionType = "TCP"
	}

	var dbURI string
	switch connectionType {
	case "Unix socket":
		dbURI = fmt.Sprintf("%s:%s@unix(/cloudsql/%s)/%s", dbUser, dbPass, instanceConnectionName, dbName)
	case "TCP":
		instanceConnectionName = "127.0.0.1:3306"
		dbURI = fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPass, instanceConnectionName, dbName)
	}

	// Open database connection.
	dbConn, err := sql.Open("mysql", dbURI)
	if err != nil {
		return nil, err
	}

	// [START cloud_sql_mysql_databasesql_limit]

	// Set maximum number of connections in idle connection pool.
	// For more information see https://golang.org/pkg/database/sql/#DB.SetMaxIdleConns
	dbConn.SetMaxIdleConns(5)

	// Set maximum number of open connections to the database.
	// For more information see https://golang.org/pkg/database/sql/#DB.SetMaxOpenConns
	dbConn.SetMaxOpenConns(7)

	// [END cloud_sql_mysql_databasesql_limit]

	// [START cloud_sql_mysql_databasesql_lifetime]

	// Set Maximum time (in seconds) that a connection can remain open.
	// For more information see https://golang.org/pkg/database/sql/#DB.SetConnMaxLifetime
	dbConn.SetConnMaxLifetime(1800)

	// [END cloud_sql_mysql_databasesql_lifetime]

	return dbConn, nil
}

// [END cloud_sql_mysql_databasesql_create]

// initDBSchema creates the votes table if it does not exist.
func initDBSchema() error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS votes
	( vote_id SERIAL NOT NULL, time_cast timestamp NOT NULL,
	candidate CHAR(6) NOT NULL, PRIMARY KEY (vote_id) );`)
	return err
}

// [END gae_cloudsql]
