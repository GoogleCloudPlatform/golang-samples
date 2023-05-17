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

// Sample database-sql demonstrates connecting to a Cloud SQL instance.
// The application is a Go version of the "Tabs vs Spaces"
// web app presented at Google Cloud Next 2019 as seen in this video:
// https://www.youtube.com/watch?v=qVgzP3PsXFw&t=1833s
package cloudsql

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	indexTmpl = template.Must(template.New("index").Parse(indexHTML))
	db        *sql.DB
	once      sync.Once
)

// getDB lazily instantiates a database connection pool. Users of Cloud Run or
// Cloud Functions may wish to skip this lazy instantiation and connect as soon
// as the function is loaded. This is primarily to help testing.
func getDB() *sql.DB {
	once.Do(func() {
		db = mustConnect()
	})
	return db
}

// migrateDB creates the votes table if it does not already exist.
func migrateDB(db *sql.DB) error {
	createVotes := `CREATE TABLE IF NOT EXISTS votes (
		id SERIAL NOT NULL,
		created_at datetime NOT NULL,
		candidate VARCHAR(6) NOT NULL,
		PRIMARY KEY (id)
	);`
	_, err := db.Exec(createVotes)
	return err
}

// recentVotes returns the last five votes cast.
func recentVotes(db *sql.DB) ([]vote, error) {
	rows, err := db.Query("SELECT candidate, created_at FROM votes ORDER BY created_at DESC LIMIT 5")
	if err != nil {
		return nil, fmt.Errorf("DB.Query: %w", err)
	}
	defer rows.Close()

	var votes []vote
	for rows.Next() {
		var (
			candidate string
			voteTime  time.Time
		)
		err := rows.Scan(&candidate, &voteTime)
		if err != nil {
			return nil, fmt.Errorf("Rows.Scan: %w", err)
		}
		votes = append(votes, vote{Candidate: candidate, VoteTime: voteTime})
	}
	return votes, nil
}

// formatMargin calculates the difference between votes and returns a human
// friendly margin (e.g., 2 votes)
func formatMargin(a, b int) string {
	diff := int(math.Abs(float64(a - b)))
	margin := fmt.Sprintf("%d votes", diff)
	// remove pluralization when diff is just one
	if diff == 1 {
		margin = "1 vote"
	}
	return margin
}

// votingData is used to pass data to the HTML template.
type votingData struct {
	TabsCount   int
	SpacesCount int
	VoteMargin  string
	RecentVotes []vote
}

// currentTotals retrieves all voting data from the database.
func currentTotals(db *sql.DB) (votingData, error) {
	var (
		tabs   int
		spaces int
	)
	err := db.QueryRow("SELECT count(id) FROM votes WHERE candidate='TABS'").Scan(&tabs)
	if err != nil {
		return votingData{}, fmt.Errorf("DB.QueryRow: %w", err)
	}
	err = db.QueryRow("SELECT count(id) FROM votes WHERE candidate='SPACES'").Scan(&spaces)
	if err != nil {
		return votingData{}, fmt.Errorf("DB.QueryRow: %w", err)
	}

	recent, err := recentVotes(db)
	if err != nil {
		return votingData{}, fmt.Errorf("recentVotes: %w", err)
	}

	return votingData{
		TabsCount:   tabs,
		SpacesCount: spaces,
		VoteMargin:  formatMargin(tabs, spaces),
		RecentVotes: recent,
	}, nil
}

// mustConnect creates a connection to the database based on environment
// variables. Setting one of INSTANCE_HOST, INSTANCE_UNIX_SOCKET, or
// INSTANCE_CONNECTION_NAME will establish a connection using a TCP socket, a
// Unix socket, or a connector respectively.
func mustConnect() *sql.DB {
	var (
		db  *sql.DB
		err error
	)

	// Use a TCP socket when INSTANCE_HOST (e.g., 127.0.0.1) is defined
	if os.Getenv("INSTANCE_HOST") != "" {
		db, err = connectTCPSocket()
		if err != nil {
			log.Fatalf("connectTCPSocket: unable to connect: %s", err)
		}
	}
	// Use a Unix socket when INSTANCE_UNIX_SOCKET (e.g., /cloudsql/proj:region:instance) is defined.
	if os.Getenv("INSTANCE_UNIX_SOCKET") != "" {
		db, err = connectUnixSocket()
		if err != nil {
			log.Fatalf("connectUnixSocket: unable to connect: %s", err)
		}
	}

	// Use the connector when INSTANCE_CONNECTION_NAME (proj:region:instance) is defined.
	if os.Getenv("INSTANCE_CONNECTION_NAME") != "" {
		if os.Getenv("DB_USER") == "" && os.Getenv("DB_IAM_USER") == "" {
			log.Fatal("Warning: One of DB_USER or DB_IAM_USER must be defined")
		}
		// Use IAM Authentication (recommended) if DB_IAM_USER is set
		if os.Getenv("DB_IAM_USER") != "" {
			db, err = connectWithConnectorIAMAuthN()
		} else {
			db, err = connectWithConnector()
		}
		if err != nil {
			log.Fatalf("connectConnector: unable to connect: %s", err)
		}
	}

	if db == nil {
		log.Fatal("Missing database connection type. Please define one of INSTANCE_HOST, INSTANCE_UNIX_SOCKET, or INSTANCE_CONNECTION_NAME")
	}

	if err := migrateDB(db); err != nil {
		log.Fatalf("unable to create table: %s", err)
	}

	return db
}

// configureConnectionPool sets database connection pool properties.
// For more information, see https://golang.org/pkg/database/sql
func configureConnectionPool(db *sql.DB) {
	// [START cloud_sql_mysql_databasesql_limit]
	// Set maximum number of connections in idle connection pool.
	db.SetMaxIdleConns(5)

	// Set maximum number of open connections to the database.
	db.SetMaxOpenConns(7)
	// [END cloud_sql_mysql_databasesql_limit]

	// [START cloud_sql_mysql_databasesql_lifetime]
	// Set Maximum time (in seconds) that a connection can remain open.
	db.SetConnMaxLifetime(1800 * time.Second)
	// [END cloud_sql_mysql_databasesql_lifetime]

	// [START cloud_sql_mysql_databasesql_backoff]
	// database/sql does not support specifying backoff
	// [END cloud_sql_mysql_databasesql_backoff]
	// [START cloud_sql_mysql_databasesql_timeout]
	// The database/sql package currently doesn't offer any functionality to
	// configure connection timeout.
	// [END cloud_sql_mysql_databasesql_timeout]
}

// Votes handles HTTP requests to alternatively show the voting app or to save a
// vote.
func Votes(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		renderIndex(w, r, getDB())
	case http.MethodPost:
		saveVote(w, r, getDB())
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// vote contains a single row from the votes table in the database. Each vote
// includes a candidate ("TABS" or "SPACES") and a timestamp.
type vote struct {
	Candidate string
	VoteTime  time.Time
}

// renderIndex renders the HTML application with the voting form, current
// totals, and recent votes.
func renderIndex(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	t, err := currentTotals(db)
	if err != nil {
		log.Printf("renderIndex: failed to read current totals: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = indexTmpl.Execute(w, t)
	if err != nil {
		log.Printf("renderIndex: failed to render template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// saveVote saves a vote passed as http.Request form data.
func saveVote(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if err := r.ParseForm(); err != nil {
		log.Printf("saveVote: failed to parse form: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	team := r.FormValue("team")
	if team == "" {
		log.Printf("saveVote: \"team\" property missing from form submission")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if team != "TABS" && team != "SPACES" {
		log.Printf("saveVote: \"team\" property should be \"TABS\" or \"SPACES\", was %q", team)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// [START cloud_sql_mysql_databasesql_connection]
	insertVote := "INSERT INTO votes(candidate, created_at) VALUES(?, NOW())"
	_, err := db.Exec(insertVote, team)
	// [END cloud_sql_mysql_databasesql_connection]

	if err != nil {
		log.Printf("saveVote: unable to save vote: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	fmt.Fprintf(w, "Vote successfully cast for %s!", team)
}

var indexHTML = `
<html lang="en">
<head>
    <title>Tabs VS Spaces</title>
    <link rel="icon" type="image/png" href="data:image/png;base64,iVBORw0KGgo=">
    <link rel="stylesheet"
          href="https://cdnjs.cloudflare.com/ajax/libs/materialize/1.0.0/css/materialize.min.css">
    <link href="https://fonts.googleapis.com/icon?family=Material+Icons" rel="stylesheet">
    <script src="https://cdnjs.cloudflare.com/ajax/libs/materialize/1.0.0/js/materialize.min.js"></script>
</head>
<body>
<nav class="red lighten-1">
    <div class="nav-wrapper">
        <a href="#" class="brand-logo center">Tabs VS Spaces</a>
    </div>
</nav>
<div class="section">
    <div class="center">
        <h4>
            {{ if eq .TabsCount .SpacesCount }}
                TABS and SPACES are evenly matched!
            {{ else if gt .TabsCount .SpacesCount }}
                TABS are winning by {{ .VoteMargin }}
            {{ else if gt .SpacesCount .TabsCount }}
                SPACES are winning by {{ .VoteMargin }}
            {{ end }}
        </h4>
    </div>
    <div class="row center">
        <div class="col s6 m5 offset-m1">
            {{ if gt .TabsCount .SpacesCount }}
			<div class="card-panel green lighten-3">
			{{ else }}
			<div class="card-panel">
			{{ end }}
                <i class="material-icons large">keyboard_tab</i>
                <h3>{{ .TabsCount }} votes</h3>
                <button id="voteTabs" class="btn green">Vote for TABS</button>
            </div>
        </div>
        <div class="col s6 m5">
            {{ if lt .TabsCount .SpacesCount }}
			<div class="card-panel blue lighten-3">
			{{ else }}
			<div class="card-panel">
			{{ end }}
                <i class="material-icons large">space_bar</i>
                <h3>{{ .SpacesCount }} votes</h3>
                <button id="voteSpaces" class="btn blue">Vote for SPACES</button>
            </div>
        </div>
    </div>
    <h4 class="header center">Recent Votes</h4>
    <ul class="container collection center">
        {{ range .RecentVotes }}
            <li class="collection-item avatar">
                {{ if eq .Candidate "TABS" }}
                    <i class="material-icons circle green">keyboard_tab</i>
                {{ else if eq .Candidate "SPACES" }}
                    <i class="material-icons circle blue">space_bar</i>
                {{ end }}
                <span class="title">
                    A vote for <b>{{.Candidate}}</b> was cast at {{.VoteTime.Format "2006-01-02T15:04:05Z07:00" }}
                </span>
            </li>
        {{ end }}
    </ul>
</div>
<script>
    function vote(team) {
        var xhr = new XMLHttpRequest();
        xhr.onreadystatechange = function () {
            if (this.readyState == 4) {
                window.location.reload();
            }
        };
        xhr.open("POST", "/Votes", true);
        xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
        xhr.send("team=" + team);
    }

    document.getElementById("voteTabs").addEventListener("click", function () {
        vote("TABS");
    });
    document.getElementById("voteSpaces").addEventListener("click", function () {
        vote("SPACES");
    });
</script>
</body>
</html>
`
