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

// Command spanner_leaderboard contains runnable snippet code for Cloud Spanner.
package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"time"

	"cloud.google.com/go/spanner"
	database "cloud.google.com/go/spanner/admin/database/apiv1"

	adminpb "cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
	"google.golang.org/api/iterator"
)

type command func(ctx context.Context, w io.Writer, client *spanner.Client) error
type uniqueRand struct {
	used map[int64]bool
	rand *rand.Rand
}

func newUniqueRand() uniqueRand {
	return uniqueRand{
		used: map[int64]bool{},
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (r *uniqueRand) rnd(min, max int64) int64 {
	for {
		rnd := r.rand.Int63n(max) + min
		if !r.used[rnd] {
			r.used[rnd] = true
			return rnd
		}
	}
}

var (
	commands = map[string]command{
		"insertplayers": insertPlayers,
		"insertscores":  insertScores,
		"query":         query,
	}
)

func createDatabase(ctx context.Context, w io.Writer, adminClient *database.DatabaseAdminClient, db string) error {
	matches := regexp.MustCompile("^(.*)/databases/(.*)$").FindStringSubmatch(db)
	if matches == nil || len(matches) != 3 {
		return fmt.Errorf("Invalid database id %s", db)
	}
	op, err := adminClient.CreateDatabase(ctx, &adminpb.CreateDatabaseRequest{
		Parent:          matches[1],
		CreateStatement: "CREATE DATABASE `" + matches[2] + "`",
		ExtraStatements: []string{
			`CREATE TABLE Players(
			    PlayerId INT64 NOT NULL,
			    PlayerName STRING(2048) NOT NULL
			) PRIMARY KEY(PlayerId)`,
			`CREATE TABLE Scores(
			    PlayerId INT64 NOT NULL,
			    Score INT64 NOT NULL,
			    Timestamp TIMESTAMP NOT NULL
			    OPTIONS(allow_commit_timestamp=true)
			) PRIMARY KEY(PlayerId, Timestamp),
			INTERLEAVE IN PARENT Players ON DELETE NO ACTION`,
		},
	})
	if err != nil {
		return err
	}
	if _, err := op.Wait(ctx); err != nil {
		return err
	}
	fmt.Fprintf(w, "Created database [%s]\n", db)
	return nil
}

func insertPlayers(ctx context.Context, w io.Writer, client *spanner.Client) error {
	// Get number of players to use as an incrementing value for each PlayerName to be inserted
	stmt := spanner.Statement{
		SQL: `SELECT Count(PlayerId) as PlayerCount FROM Players`,
	}
	iter := client.Single().Query(ctx, stmt)
	defer iter.Stop()
	row, err := iter.Next()
	if err != nil {
		return err
	}
	var numberOfPlayers int64 = 0
	if err := row.Columns(&numberOfPlayers); err != nil {
		return err
	}
	// Intialize values for random PlayerId
	min := int64(1000000000)
	max := int64(9000000000)
	rnd := newUniqueRand()
	// Insert 100 player records into the Players table
	_, err = client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		stmts := []spanner.Statement{}
		for i := 1; i <= 100; i++ {
			numberOfPlayers++
			playerID := rnd.rnd(min, max)

			playerName := fmt.Sprintf("Player %d", numberOfPlayers)
			stmts = append(stmts, spanner.Statement{
				SQL: `INSERT INTO Players
						(PlayerId, PlayerName)
						VALUES (@playerID, @playerName)`,
				Params: map[string]interface{}{
					"playerID":   playerID,
					"playerName": playerName,
				},
			})
		}
		_, err := txn.BatchUpdate(ctx, stmts)
		if err != nil {
			return err
		}
		return nil
	})
	fmt.Fprintf(w, "Inserted players \n")
	return nil
}

func insertScores(ctx context.Context, w io.Writer, client *spanner.Client) error {
	playerRecordsFound := false
	// Create slice for insert statements
	stmts := []spanner.Statement{}
	// Select all player records
	stmt := spanner.Statement{SQL: `SELECT PlayerId FROM Players`}
	iter := client.Single().Query(ctx, stmt)
	defer iter.Stop()
	// Insert 4 score records into the Scores table for each player in the Players table
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		playerRecordsFound = true
		var playerID int64
		if err := row.ColumnByName("PlayerId", &playerID); err != nil {
			return err
		}
		// Intialize values for random score and date
		rand.Seed(time.Now().UnixNano())
		min := 1000
		max := 1000000
		for i := 0; i < 4; i++ {
			// Generate random score between 1,000 and 1,000,000
			score := rand.Intn(max-min) + min
			// Generate random day within the past two years
			now := time.Now()
			endDate := now.Unix()
			past := now.AddDate(0, -24, 0)
			startDate := past.Unix()
			randomDateInSeconds := rand.Int63n(endDate-startDate) + startDate
			randomDate := time.Unix(randomDateInSeconds, 0)
			// Add insert statement to stmts slice
			stmts = append(stmts, spanner.Statement{
				SQL: `INSERT INTO Scores
						(PlayerId, Score, Timestamp)
						VALUES (@playerID, @score, @timestamp)`,
				Params: map[string]interface{}{
					"playerID":  playerID,
					"score":     score,
					"timestamp": randomDate,
				},
			})
		}

	}
	if !playerRecordsFound {
		fmt.Fprintln(w, "No player records currently exist. First insert players then insert scores.")
	} else {
		_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
			// Commit insert statements for all scores to be inserted as a single transaction
			_, err := txn.BatchUpdate(ctx, stmts)
			return err
		})
		if err != nil {
			return err
		}
		fmt.Fprintln(w, "Inserted scores")
	}
	return nil
}

func query(ctx context.Context, w io.Writer, client *spanner.Client) error {
	stmt := spanner.Statement{
		SQL: `SELECT p.PlayerId, p.PlayerName, s.Score, s.Timestamp
		        FROM Players p
		        JOIN Scores s ON p.PlayerId = s.PlayerId
		        ORDER BY s.Score DESC LIMIT 10`}
	iter := client.Single().Query(ctx, stmt)
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			return nil
		}
		if err != nil {
			return err
		}
		var playerID, score int64
		var playerName string
		var timestamp time.Time
		if err := row.Columns(&playerID, &playerName, &score, &timestamp); err != nil {
			return err
		}
		fmt.Fprintf(w, "PlayerId: %d  PlayerName: %s  Score: %s  Timestamp: %s\n",
			playerID, playerName, formatWithCommas(score), timestamp.String()[0:10])
	}
}

func queryWithTimespan(ctx context.Context, w io.Writer, client *spanner.Client, timespan int) error {
	stmt := spanner.Statement{
		SQL: `SELECT p.PlayerId, p.PlayerName, s.Score, s.Timestamp
				FROM Players p
				JOIN Scores s ON p.PlayerId = s.PlayerId 
				WHERE s.Timestamp > TIMESTAMP_SUB(CURRENT_TIMESTAMP(), INTERVAL @Timespan HOUR)
				ORDER BY s.Score DESC LIMIT 10`,
		Params: map[string]interface{}{"Timespan": timespan},
	}
	iter := client.Single().Query(ctx, stmt)
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			return nil
		}
		if err != nil {
			return err
		}
		var playerID, score int64
		var playerName string
		var timestamp time.Time
		if err := row.Columns(&playerID, &playerName, &score, &timestamp); err != nil {
			return err
		}
		fmt.Fprintf(w, "PlayerId: %d  PlayerName: %s  Score: %s  Timestamp: %s\n",
			playerID, playerName, formatWithCommas(score), timestamp.String()[0:10])
	}
}

func formatWithCommas(n int64) string {
	numberAsString := strconv.FormatInt(n, 10)
	numberLength := len(numberAsString)
	if numberLength < 4 {
		return numberAsString
	}
	var buffer bytes.Buffer
	comma := []rune(",")
	bufferPosition := numberLength % 3
	if (bufferPosition) > 0 {
		bufferPosition = 3 - bufferPosition
	}
	for i := 0; i < numberLength; i++ {
		if bufferPosition == 3 {
			buffer.WriteRune(comma[0])
			bufferPosition = 0
		}
		bufferPosition++
		buffer.WriteByte(numberAsString[i])
	}
	return buffer.String()
}

func createClients(ctx context.Context, db string) (*database.DatabaseAdminClient, *spanner.Client) {
	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	dataClient, err := spanner.NewClient(ctx, db)
	if err != nil {
		log.Fatal(err)
	}

	return adminClient, dataClient
}

func run(ctx context.Context, adminClient *database.DatabaseAdminClient, dataClient *spanner.Client, w io.Writer,
	cmd string, db string, timespan int) error {
	// createdatabase command
	if cmd == "createdatabase" {
		err := createDatabase(ctx, w, adminClient, db)
		if err != nil {
			fmt.Fprintf(w, "%s failed with %v", cmd, err)
		}
		return err
	}

	// querywithtimespan command
	if cmd == "querywithtimespan" {
		if timespan == 0 {
			flag.Usage()
			os.Exit(2)
		}
		err := queryWithTimespan(ctx, w, dataClient, timespan)
		if err != nil {
			fmt.Fprintf(w, "%s failed with %v", cmd, err)
		}
		return err
	}

	// insert and query commands
	cmdFn := commands[cmd]
	if cmdFn == nil {
		flag.Usage()
		os.Exit(2)
	}
	err := cmdFn(ctx, w, dataClient)
	if err != nil {
		fmt.Fprintf(w, "%s failed with %v", cmd, err)
	}
	return err
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: leaderboard <command> <database_name> [command_option]

	Command can be one of: createdatabase, insertplayers, insertscores, query, querywithtimespan

Examples:
	leaderboard createdatabase projects/my-project/instances/my-instance/databases/example-db
		- Create a sample Cloud Spanner database along with sample tables in your project.
	leaderboard insertplayers projects/my-project/instances/my-instance/databases/example-db
		- Insert 100 sample Player records into the database.
	leaderboard insertscores projects/my-project/instances/my-instance/databases/example-db
		- Insert sample score data into Scores sample Cloud Spanner database table.
	leaderboard query projects/my-project/instances/my-instance/databases/example-db
		- Query players with top ten scores of all time.
	leaderboard querywithtimespan projects/my-project/instances/my-instance/databases/example-db 168
		- Query players with top ten scores within a timespan specified in hours.
`)
	}

	flag.Parse()
	flagCount := len(flag.Args())
	if flagCount < 2 || flagCount > 3 {
		flag.Usage()
		os.Exit(2)
	}

	cmd, db := flag.Arg(0), flag.Arg(1)
	// If query timespan flag is specified, parse to int
	var timespan int = 0
	if flagCount == 3 {
		parsedTimespan, err := strconv.Atoi(flag.Arg(2))
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		timespan = parsedTimespan
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	adminClient, dataClient := createClients(ctx, db)
	if err := run(ctx, adminClient, dataClient, os.Stdout, cmd, db, timespan); err != nil {
		os.Exit(1)
	}
}
