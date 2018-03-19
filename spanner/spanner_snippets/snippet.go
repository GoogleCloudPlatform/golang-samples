// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Command spanner_snippets contains runnable snippet code for Cloud Spanner.
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	"cloud.google.com/go/spanner"
	database "cloud.google.com/go/spanner/admin/database/apiv1"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"

	adminpb "google.golang.org/genproto/googleapis/spanner/admin/database/v1"
)

type command func(ctx context.Context, w io.Writer, client *spanner.Client) error
type adminCommand func(ctx context.Context, w io.Writer, adminClient *database.DatabaseAdminClient, database string) error

var (
	commands = map[string]command{
		"write":               write,
		"query":               query,
		"read":                read,
		"update":              update,
		"writetransaction":    writeWithTransaction,
		"querynewcolumn":      queryNewColumn,
		"queryindex":          queryUsingIndex,
		"readindex":           readUsingIndex,
		"readstoringindex":    readStoringIndex,
		"readonlytransaction": readOnlyTransaction,
		"readstaledata":       readStaleData,
		"readbatchdata":       readBatchData,
	}

	adminCommands = map[string]adminCommand{
		"createdatabase":  createDatabase,
		"addnewcolumn":    addNewColumn,
		"addindex":        addIndex,
		"addstoringindex": addStoringIndex,
	}
)

// [START spanner_create_database]

func createDatabase(ctx context.Context, w io.Writer, adminClient *database.DatabaseAdminClient, db string) error {
	matches := regexp.MustCompile("^(.*)/databases/(.*)$").FindStringSubmatch(db)
	if matches == nil || len(matches) != 3 {
		return fmt.Errorf("Invalid database id %s", db)
	}
	op, err := adminClient.CreateDatabase(ctx, &adminpb.CreateDatabaseRequest{
		Parent:          matches[1],
		CreateStatement: "CREATE DATABASE `" + matches[2] + "`",
		ExtraStatements: []string{
			`CREATE TABLE Singers (
				SingerId   INT64 NOT NULL,
				FirstName  STRING(1024),
				LastName   STRING(1024),
				SingerInfo BYTES(MAX)
			) PRIMARY KEY (SingerId)`,
			`CREATE TABLE Albums (
				SingerId	INT64 NOT NULL,
				AlbumId		INT64 NOT NULL,
				AlbumTitle	STRING(MAX),
			) PRIMARY KEY (SingerId, AlbumId),
			INTERLEAVE IN PARENT Singers ON DELETE CASCADE`,
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

// [END spanner_create_database]

// [START spanner_insert_data]

func write(ctx context.Context, w io.Writer, client *spanner.Client) error {
	singerColumns := []string{"SingerId", "FirstName", "LastName"}
	albumColumns := []string{"SingerId", "AlbumId", "AlbumTitle"}
	m := []*spanner.Mutation{
		spanner.InsertOrUpdate("Singers", singerColumns, []interface{}{1, "Marc", "Richards"}),
		spanner.InsertOrUpdate("Singers", singerColumns, []interface{}{2, "Catalina", "Smith"}),
		spanner.InsertOrUpdate("Singers", singerColumns, []interface{}{3, "Alice", "Trentor"}),
		spanner.InsertOrUpdate("Singers", singerColumns, []interface{}{4, "Lea", "Martin"}),
		spanner.InsertOrUpdate("Singers", singerColumns, []interface{}{5, "David", "Lomond"}),
		spanner.InsertOrUpdate("Albums", albumColumns, []interface{}{1, 1, "Total Junk"}),
		spanner.InsertOrUpdate("Albums", albumColumns, []interface{}{1, 2, "Go, Go, Go"}),
		spanner.InsertOrUpdate("Albums", albumColumns, []interface{}{2, 1, "Green"}),
		spanner.InsertOrUpdate("Albums", albumColumns, []interface{}{2, 2, "Forever Hold Your Peace"}),
		spanner.InsertOrUpdate("Albums", albumColumns, []interface{}{2, 3, "Terrified"}),
	}
	_, err := client.Apply(ctx, m)
	return err
}

// [END spanner_insert_data]

// [START spanner_query_data]

func query(ctx context.Context, w io.Writer, client *spanner.Client) error {
	stmt := spanner.Statement{SQL: `SELECT SingerId, AlbumId, AlbumTitle FROM Albums`}
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
		var singerID, albumID int64
		var albumTitle string
		if err := row.Columns(&singerID, &albumID, &albumTitle); err != nil {
			return err
		}
		fmt.Fprintf(w, "%d %d %s\n", singerID, albumID, albumTitle)
	}
}

// [END spanner_query_data]

// [START spanner_read_data]

func read(ctx context.Context, w io.Writer, client *spanner.Client) error {
	iter := client.Single().Read(ctx, "Albums", spanner.AllKeys(),
		[]string{"SingerId", "AlbumId", "AlbumTitle"})
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			return nil
		}
		if err != nil {
			return err
		}
		var singerID, albumID int64
		var albumTitle string
		if err := row.Columns(&singerID, &albumID, &albumTitle); err != nil {
			return err
		}
		fmt.Fprintf(w, "%d %d %s\n", singerID, albumID, albumTitle)
	}
}

// [END spanner_read_data]

// [START spanner_add_column]

func addNewColumn(ctx context.Context, w io.Writer, adminClient *database.DatabaseAdminClient, database string) error {
	op, err := adminClient.UpdateDatabaseDdl(ctx, &adminpb.UpdateDatabaseDdlRequest{
		Database: database,
		Statements: []string{
			"ALTER TABLE Albums ADD COLUMN MarketingBudget INT64",
		},
	})
	if err != nil {
		return err
	}
	if err := op.Wait(ctx); err != nil {
		return err
	}
	fmt.Fprintf(w, "Added MarketingBudget column\n")
	return nil
}

// [END spanner_add_column]

// [START spanner_update_data]

func update(ctx context.Context, w io.Writer, client *spanner.Client) error {
	cols := []string{"SingerId", "AlbumId", "MarketingBudget"}
	_, err := client.Apply(ctx, []*spanner.Mutation{
		spanner.Update("Albums", cols, []interface{}{1, 1, 100000}),
		spanner.Update("Albums", cols, []interface{}{2, 2, 500000}),
	})
	return err
}

// [END spanner_update_data]

// [START spanner_read_write_transaction]

func writeWithTransaction(ctx context.Context, w io.Writer, client *spanner.Client) error {
	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		getBudget := func(key spanner.Key) (int64, error) {
			row, err := txn.ReadRow(ctx, "Albums", key, []string{"MarketingBudget"})
			if err != nil {
				return 0, err
			}
			var budget int64
			if err := row.Column(0, &budget); err != nil {
				return 0, err
			}
			return budget, nil
		}
		album2Budget, err := getBudget(spanner.Key{2, 2})
		if err != nil {
			return err
		}
		if album2Budget >= 300000 {
			album1Budget, err := getBudget(spanner.Key{1, 1})
			if err != nil {
				return err
			}
			const transfer = 200000
			album1Budget += transfer
			album2Budget -= transfer
			cols := []string{"SingerId", "AlbumId", "MarketingBudget"}
			txn.BufferWrite([]*spanner.Mutation{
				spanner.Update("Albums", cols, []interface{}{1, 1, album1Budget}),
				spanner.Update("Albums", cols, []interface{}{2, 2, album2Budget}),
			})
		}
		return nil
	})
	return err
}

// [END spanner_read_write_transaction]

// [START spanner_query_data_with_new_column]

func queryNewColumn(ctx context.Context, w io.Writer, client *spanner.Client) error {
	stmt := spanner.Statement{SQL: `SELECT SingerId, AlbumId, MarketingBudget FROM Albums`}
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
		var singerID, albumID int64
		var marketingBudget spanner.NullInt64
		if err := row.ColumnByName("SingerId", &singerID); err != nil {
			return err
		}
		if err := row.ColumnByName("AlbumId", &albumID); err != nil {
			return err
		}
		if err := row.ColumnByName("MarketingBudget", &marketingBudget); err != nil {
			return err
		}
		budget := "NULL"
		if marketingBudget.Valid {
			budget = strconv.FormatInt(marketingBudget.Int64, 10)
		}
		fmt.Fprintf(w, "%d %d %s\n", singerID, albumID, budget)
	}
}

// [END spanner_query_data_with_new_column]

// [START spanner_create_index]

func addIndex(ctx context.Context, w io.Writer, adminClient *database.DatabaseAdminClient, database string) error {
	op, err := adminClient.UpdateDatabaseDdl(ctx, &adminpb.UpdateDatabaseDdlRequest{
		Database: database,
		Statements: []string{
			"CREATE INDEX AlbumsByAlbumTitle ON Albums(AlbumTitle)",
		},
	})
	if err != nil {
		return err
	}
	if err := op.Wait(ctx); err != nil {
		return err
	}
	fmt.Fprintf(w, "Added index\n")
	return nil
}

// [END spanner_create_index]

// [START spanner_query_data_with_index]

func queryUsingIndex(ctx context.Context, w io.Writer, client *spanner.Client) error {
	stmt := spanner.Statement{
		SQL: `SELECT AlbumId, AlbumTitle, MarketingBudget
			FROM Albums@{FORCE_INDEX=AlbumsByAlbumTitle}
			WHERE AlbumTitle >= @start_title AND AlbumTitle < @end_title`,
		Params: map[string]interface{}{
			"start_title": "Aardvark",
			"end_title":   "Goo",
		},
	}
	iter := client.Single().Query(ctx, stmt)
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		var albumID int64
		var marketingBudget spanner.NullInt64
		var albumTitle string
		if err := row.ColumnByName("AlbumId", &albumID); err != nil {
			return err
		}
		if err := row.ColumnByName("AlbumTitle", &albumTitle); err != nil {
			return err
		}
		if err := row.ColumnByName("MarketingBudget", &marketingBudget); err != nil {
			return err
		}
		budget := "NULL"
		if marketingBudget.Valid {
			budget = strconv.FormatInt(marketingBudget.Int64, 10)
		}
		fmt.Fprintf(w, "%d %s %s\n", albumID, albumTitle, budget)
	}
	return nil
}

// [END spanner_query_data_with_index]

// [START spanner_read_data_with_index]

func readUsingIndex(ctx context.Context, w io.Writer, client *spanner.Client) error {
	iter := client.Single().ReadUsingIndex(ctx, "Albums", "AlbumsByAlbumTitle", spanner.AllKeys(),
		[]string{"AlbumId", "AlbumTitle"})
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			return nil
		}
		if err != nil {
			return err
		}
		var albumID int64
		var albumTitle string
		if err := row.Columns(&albumID, &albumTitle); err != nil {
			return err
		}
		fmt.Fprintf(w, "%d %s\n", albumID, albumTitle)
	}
}

// [END spanner_read_data_with_index]

// [START spanner_create_storing_index]

func addStoringIndex(ctx context.Context, w io.Writer, adminClient *database.DatabaseAdminClient, database string) error {
	op, err := adminClient.UpdateDatabaseDdl(ctx, &adminpb.UpdateDatabaseDdlRequest{
		Database: database,
		Statements: []string{
			"CREATE INDEX AlbumsByAlbumTitle2 ON Albums(AlbumTitle) STORING (MarketingBudget)",
		},
	})
	if err != nil {
		return err
	}
	if err := op.Wait(ctx); err != nil {
		return err
	}
	fmt.Fprintf(w, "Added storing index\n")
	return nil
}

// [END spanner_create_storing_index]

// [START spanner_read_data_with_storing_index]

func readStoringIndex(ctx context.Context, w io.Writer, client *spanner.Client) error {
	iter := client.Single().ReadUsingIndex(ctx, "Albums", "AlbumsByAlbumTitle2", spanner.AllKeys(),
		[]string{"AlbumId", "AlbumTitle", "MarketingBudget"})
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			return nil
		}
		if err != nil {
			return err
		}
		var albumID int64
		var marketingBudget spanner.NullInt64
		var albumTitle string
		if err := row.Columns(&albumID, &albumTitle, &marketingBudget); err != nil {
			return err
		}
		budget := "NULL"
		if marketingBudget.Valid {
			budget = strconv.FormatInt(marketingBudget.Int64, 10)
		}
		fmt.Fprintf(w, "%d %s %s\n", albumID, albumTitle, budget)
	}
}

// [END spanner_read_data_with_storing_index]

// [START spanner_read_only_transaction]

func readOnlyTransaction(ctx context.Context, w io.Writer, client *spanner.Client) error {
	ro := client.ReadOnlyTransaction()
	defer ro.Close()
	stmt := spanner.Statement{SQL: `SELECT SingerId, AlbumId, AlbumTitle FROM Albums`}
	iter := ro.Query(ctx, stmt)
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		var singerID int64
		var albumID int64
		var albumTitle string
		if err := row.Columns(&singerID, &albumID, &albumTitle); err != nil {
			return err
		}
		fmt.Fprintf(w, "%d %d %s\n", singerID, albumID, albumTitle)
	}

	iter = ro.Read(ctx, "Albums", spanner.AllKeys(), []string{"SingerId", "AlbumId", "AlbumTitle"})
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			return nil
		}
		if err != nil {
			return err
		}
		var singerID int64
		var albumID int64
		var albumTitle string
		if err := row.Columns(&singerID, &albumID, &albumTitle); err != nil {
			return err
		}
		fmt.Fprintf(w, "%d %d %s\n", singerID, albumID, albumTitle)
	}
}

// [END spanner_read_only_transaction]

// [START spanner_read_stale_data]

func readStaleData(ctx context.Context, w io.Writer, client *spanner.Client) error {
	ro := client.ReadOnlyTransaction().WithTimestampBound(spanner.ExactStaleness(15 * time.Second))
	defer ro.Close()

	iter := ro.Read(ctx, "Albums", spanner.AllKeys(), []string{"SingerId", "AlbumId", "AlbumTitle"})
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			return nil
		}
		if err != nil {
			return err
		}
		var singerID int64
		var albumID int64
		var albumTitle string
		if err := row.Columns(&singerID, &albumID, &albumTitle); err != nil {
			return err
		}
		fmt.Fprintf(w, "%d %d %s\n", singerID, albumID, albumTitle)
	}
}

// [END spanner_read_stale_data]

// [START spanner_batch_client]

func readBatchData(ctx context.Context, w io.Writer, client *spanner.Client) error {
	txn, err := client.BatchReadOnlyTransaction(ctx, spanner.StrongRead())
	if err != nil {
		return err
	}
	defer txn.Close()

	// Singer represents a row in the Singers table.
	type Singer struct {
		SingerID   int64
		FirstName  string
		LastName   string
		SingerInfo []byte
	}
	stmt := spanner.Statement{SQL: "SELECT SingerId, FirstName, LastName FROM Singers;"}
	partitions, err := txn.PartitionQuery(ctx, stmt, spanner.PartitionOptions{})
	if err != nil {
		return err
	}
	recordCount := 0
	for i, p := range partitions {
		iter := txn.Execute(ctx, p)
		defer iter.Stop()
		for {
			row, err := iter.Next()
			if err == iterator.Done {
				break
			} else if err != nil {
				return err
			}
			var s Singer
			if err := row.ToStruct(&s); err != nil {
				return err
			}
			fmt.Fprintf(w, "Partition (%d) %v\n", i, s)
			recordCount++
		}
	}
	fmt.Fprintf(w, "Total partition count: %v\n", len(partitions))
	fmt.Fprintf(w, "Total record count: %v\n", recordCount)
	return nil
}

// [END spanner_batch_client]

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

func run(ctx context.Context, adminClient *database.DatabaseAdminClient, dataClient *spanner.Client, w io.Writer, cmd string, db string) error {
	if adminCmdFn := adminCommands[cmd]; adminCmdFn != nil {
		err := adminCmdFn(ctx, w, adminClient, db)
		if err != nil {
			fmt.Fprintf(w, "%s failed with %v", cmd, err)
		}
		return err
	}

	// Normal mode
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
		fmt.Fprintf(os.Stderr, `Usage: spanner_snippets <command> <database_name>

	Command can be one of: createdatabase, write, query, read, update,
		writetransaction, addnewcolumn, querynewcolumn, addindex, queryindex, readindex,
		addstoringindex, readstoringindex, readonlytransaction, readstaledata, readbatchdata

Examples:
	spanner_snippets createdatabase projects/my-project/instances/my-instance/databases/example-db
	spanner_snippets write projects/my-project/instances/my-instance/databases/example-db
`)
	}

	flag.Parse()
	if len(flag.Args()) != 2 {
		flag.Usage()
		os.Exit(2)
	}

	cmd, db := flag.Arg(0), flag.Arg(1)
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Minute)
	adminClient, dataClient := createClients(ctx, db)
	if err := run(ctx, adminClient, dataClient, os.Stdout, cmd, db); err != nil {
		os.Exit(1)
	}
}
