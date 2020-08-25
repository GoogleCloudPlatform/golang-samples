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
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

// vote struct contains a single row from the votes table in the database.
// Each vote includes a candidate ("TABS" or "SPACES") and a timestamp.
type vote struct {
	Candidate string
	VoteTime  sql.NullTime
}

// voteDiff is used to provide a string representation of the current voting
// margin, such as "1 vote" (singular) or "2 votes" (plural).
type voteDiff int

func (v voteDiff) String() string {
	if v == 1 {
		return "1 vote"
	}
	return strconv.Itoa(int(v)) + " votes"
}

// templateData struct is used to pass data to the HTML template.
type templateData struct {
	TabsCount   uint
	SpacesCount uint
	VoteMargin  string
	RecentVotes []vote
}

// app struct contains global state.
type app struct {
	// db is the global database connection pool.
	db *sql.DB
	// tmpl is the parsed HTML template.
	tmpl *template.Template
}

// indexHandler handles requests to the / route.
func (app *app) indexHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if err := showTotals(w, r, app); err != nil {
			log.Printf("showTotals: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	case "POST":
		if err := saveVote(w, r, app); err != nil {
			log.Printf("saveVote: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	default:
		http.Error(w, fmt.Sprintf("HTTP Method %s Not Allowed", r.Method), http.StatusMethodNotAllowed)
	}
}

func main() {
	parsedTemplate, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatalf("unable to parse template file: %s", err)
	}

	app := &app{
		tmpl: parsedTemplate,
	}

	// If the optional DB_TCP_HOST environment variable is set, it contains
	// the IP address and port number of a TCP connection pool to be created,
	// such as "127.0.0.1:3306". If DB_TCP_HOST is not set, a Unix socket
	// connection pool will be created instead.
	if os.Getenv("DB_TCP_HOST") != "" {
		app.db, err = initTCPConnectionPool()
		if err != nil {
			log.Fatalf("initTCPConnectionPool: unable to connect: %v", err)
		}
	} else {
		app.db, err = initSocketConnectionPool()
		if err != nil {
			log.Fatalf("initSocketConnectionPool: unable to connect: %v", err)
		}
	}

	// Create the votes table if it does not already exist.
	if _, err = app.db.Exec(`CREATE TABLE IF NOT EXISTS votes
	( vote_id SERIAL NOT NULL, time_cast timestamp NOT NULL,
	candidate CHAR(6) NOT NULL, PRIMARY KEY (vote_id) );`); err != nil {
		log.Fatalf("DB.Exec: unable to create table: %s", err)
	}

	http.HandleFunc("/", app.indexHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}

}

// recentVotes returns a slice of the last 5 votes cast.
func recentVotes(app *app) ([]vote, error) {
	var votes []vote
	rows, err := app.db.Query(`SELECT candidate, time_cast FROM votes ORDER BY time_cast DESC LIMIT 5`)
	if err != nil {
		return votes, fmt.Errorf("DB.Query: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		nextVote := vote{}
		err := rows.Scan(&nextVote.Candidate, &nextVote.VoteTime)
		if err != nil {
			return votes, fmt.Errorf("Rows.Scan: %v", err)
		}
		votes = append(votes, nextVote)
	}
	return votes, nil
}

// currentTotals returns a templateData structure for populating the web page.
func currentTotals(app *app) (*templateData, error) {
	// get total votes for each candidate
	var tabVotes, spaceVotes uint
	err := app.db.QueryRow(`SELECT count(vote_id) FROM votes WHERE candidate='TABS'`).Scan(&tabVotes)
	if err != nil {
		return nil, fmt.Errorf("DB.QueryRow: %v", err)
	}
	err = app.db.QueryRow(`SELECT count(vote_id) FROM votes WHERE candidate='SPACES'`).Scan(&spaceVotes)
	if err != nil {
		return nil, fmt.Errorf("DB.QueryRow: %v", err)
	}

	var voteDiffStr string = voteDiff(int(math.Abs(float64(tabVotes) - float64(spaceVotes)))).String()

	latestVotesCast, err := recentVotes(app)
	if err != nil {
		return nil, fmt.Errorf("recentVotes: %v", err)
	}

	pageData := templateData{
		TabsCount:   tabVotes,
		SpacesCount: spaceVotes,
		VoteMargin:  voteDiffStr,
		RecentVotes: latestVotesCast,
	}

	return &pageData, nil
}

// showTotals renders an HTML template showing the current vote totals.
func showTotals(w http.ResponseWriter, r *http.Request, app *app) error {
	totals, err := currentTotals(app)
	if err != nil {
		return fmt.Errorf("currentTotals: %v", err)
	}
	err = app.tmpl.Execute(w, totals)
	if err != nil {
		return fmt.Errorf("Template.Execute: %v", err)
	}
	return nil
}

// saveVote saves a vote passed as http.Request form data.
func saveVote(w http.ResponseWriter, r *http.Request, app *app) error {
	if err := r.ParseForm(); err != nil {
		return fmt.Errorf("Request.ParseForm: %v", err)
	}

	team := r.FormValue("team")
	if team == "" {
		return fmt.Errorf("team property missing from form submission")
	}

	// [START cloud_sql_mysql_databasesql_connection]
	sqlInsert := "INSERT INTO votes(candidate, time_cast) VALUES(?, NOW())"
	if team == "TABS" || team == "SPACES" {
		if _, err := app.db.Exec(sqlInsert, team); err != nil {
			fmt.Fprintf(w, "unable to save vote: %s", err)
			return fmt.Errorf("DB.Exec: %v", err)
		} else {
			fmt.Fprintf(w, "Vote successfully cast for %s!\n", team)
		}
	}
	return nil
	// [END cloud_sql_mysql_databasesql_connection]
}

// mustGetEnv is a helper function for getting environment variables.
// Displays a warning if the environment variable is not set.
func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("Warning: %s environment variable not set.\n", k)
	}
	return v
}

// initSocketConnectionPool initializes a Unix socket connection pool for
// a Cloud SQL instance of SQL Server.
func initSocketConnectionPool() (*sql.DB, error) {
	// [START cloud_sql_mysql_databasesql_create_socket]
	var (
		dbUser                 = mustGetenv("DB_USER")
		dbPwd                  = mustGetenv("DB_PASS")
		instanceConnectionName = mustGetenv("INSTANCE_CONNECTION_NAME")
		dbName                 = mustGetenv("DB_NAME")
	)

	socketDir, isSet := os.LookupEnv("DB_SOCKET_DIR")
	if !isSet {
		socketDir = "/cloudsql"
	}

	var dbURI string
	dbURI = fmt.Sprintf("%s:%s@unix(/%s/%s)/%s?parseTime=true", dbUser, dbPwd, socketDir, instanceConnectionName, dbName)

	// dbPool is the pool of database connections.
	dbPool, err := sql.Open("mysql", dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %v", err)
	}

	// [START_EXCLUDE]
	configureConnectionPool(dbPool)
	// [END_EXCLUDE]

	return dbPool, nil
	// [END cloud_sql_mysql_databasesql_create_socket]
}

// initTCPConnectionPool initializes a TCP connection pool for a Cloud SQL
// instance of SQL Server.
func initTCPConnectionPool() (*sql.DB, error) {
	// [START cloud_sql_mysql_databasesql_create_tcp]
	var (
		dbUser    = mustGetenv("DB_USER")
		dbPwd     = mustGetenv("DB_PASS")
		dbTcpHost = mustGetenv("DB_TCP_HOST")
		dbPort    = mustGetenv("DB_PORT")
		dbName    = mustGetenv("DB_NAME")
	)

	var dbURI string
	dbURI = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPwd, dbTcpHost, dbPort, dbName)

	// dbPool is the pool of database connections.
	dbPool, err := sql.Open("mysql", dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %v", err)
	}

	// [START_EXCLUDE]
	configureConnectionPool(dbPool)
	// [END_EXCLUDE]

	return dbPool, nil
	// [END cloud_sql_mysql_databasesql_create_tcp]
}

// configureConnectionPool sets database connection pool properties.
// For more information, see https://golang.org/pkg/database/sql
func configureConnectionPool(dbPool *sql.DB) {
	// [START cloud_sql_mysql_databasesql_limit]

	// Set maximum number of connections in idle connection pool.
	dbPool.SetMaxIdleConns(5)

	// Set maximum number of open connections to the database.
	dbPool.SetMaxOpenConns(7)

	// [END cloud_sql_mysql_databasesql_limit]

	// [START cloud_sql_mysql_databasesql_lifetime]

	// Set Maximum time (in seconds) that a connection can remain open.
	dbPool.SetConnMaxLifetime(1800)

	// [END cloud_sql_mysql_databasesql_lifetime]
}
