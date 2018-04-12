// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package main contains runnable code demonstrating Cloud Spanner's commit timestamp feature.
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
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
		"write":             write,
		"update":            update,
		"query":             query,
		"writewithhistory":  writeWithHistory,
		"updatewithhistory": updateWithHistory,
		"querywithhistory":  queryWithHistory,
	}

	adminCommands = map[string]adminCommand{
		"createdatabase":              createDatabase,
		"createtablewithtimestamp":    createTableWithTimestamp,
		"createtablewithhistorytable": createTableWithHistoryTable,
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

func createTableWithTimestamp(ctx context.Context, w io.Writer, adminClient *database.DatabaseAdminClient, database string) error {
	op, err := adminClient.UpdateDatabaseDdl(ctx, &adminpb.UpdateDatabaseDdlRequest{
		Database: database,
		Statements: []string{
			`CREATE TABLE DocumentsWithTimestamp(
				UserId INT64 NOT NULL,
				DocumentId INT64 NOT NULL,
			    Ts TIMESTAMP NOT NULL OPTIONS(allow_commit_timestamp=true),
				Contents STRING(MAX) NOT NULL
			) PRIMARY KEY(UserId, DocumentId)`,
		},
	})
	if err != nil {
		return err
	}
	if err := op.Wait(ctx); err != nil {
		return err
	}
	fmt.Fprintf(w, "Created DocumentsWithTimestamp table in database [%s]\n", database)
	return nil
}

func write(ctx context.Context, w io.Writer, client *spanner.Client) error {
	DocumentsColumns := []string{"UserId", "DocumentId", "Ts", "Contents"}
	m := []*spanner.Mutation{
		spanner.InsertOrUpdate("DocumentsWithTimestamp", DocumentsColumns,
			[]interface{}{1, 1, spanner.CommitTimestamp, "Hello World 1"}),
		spanner.InsertOrUpdate("DocumentsWithTimestamp", DocumentsColumns,
			[]interface{}{1, 2, spanner.CommitTimestamp, "Hello World 2"}),
		spanner.InsertOrUpdate("DocumentsWithTimestamp", DocumentsColumns,
			[]interface{}{1, 3, spanner.CommitTimestamp, "Hello World 3"}),
		spanner.InsertOrUpdate("DocumentsWithTimestamp", DocumentsColumns,
			[]interface{}{2, 4, spanner.CommitTimestamp, "Hello World 4"}),
		spanner.InsertOrUpdate("DocumentsWithTimestamp", DocumentsColumns,
			[]interface{}{2, 5, spanner.CommitTimestamp, "Hello World 5"}),
		spanner.InsertOrUpdate("DocumentsWithTimestamp", DocumentsColumns,
			[]interface{}{3, 6, spanner.CommitTimestamp, "Hello World 6"}),
		spanner.InsertOrUpdate("DocumentsWithTimestamp", DocumentsColumns,
			[]interface{}{3, 7, spanner.CommitTimestamp, "Hello World 7"}),
		spanner.InsertOrUpdate("DocumentsWithTimestamp", DocumentsColumns,
			[]interface{}{3, 8, spanner.CommitTimestamp, "Hello World 8"}),
		spanner.InsertOrUpdate("DocumentsWithTimestamp", DocumentsColumns,
			[]interface{}{3, 9, spanner.CommitTimestamp, "Hello World 9"}),
		spanner.InsertOrUpdate("DocumentsWithTimestamp", DocumentsColumns,
			[]interface{}{3, 10, spanner.CommitTimestamp, "Hello World 10"}),
	}
	_, err := client.Apply(ctx, m)
	return err
}

func update(ctx context.Context, w io.Writer, client *spanner.Client) error {
	cols := []string{"UserId", "DocumentId", "Ts", "Contents"}
	_, err := client.Apply(ctx, []*spanner.Mutation{
		spanner.Update("DocumentsWithTimestamp", cols,
			[]interface{}{1, 1, spanner.CommitTimestamp, "Hello World 1 Updated"}),
		spanner.Update("DocumentsWithTimestamp", cols,
			[]interface{}{1, 3, spanner.CommitTimestamp, "Hello World 3 Updated"}),
		spanner.Update("DocumentsWithTimestamp", cols,
			[]interface{}{2, 5, spanner.CommitTimestamp, "Hello World 5 Updated"}),
		spanner.Update("DocumentsWithTimestamp", cols,
			[]interface{}{3, 7, spanner.CommitTimestamp, "Hello World 7 Updated"}),
		spanner.Update("DocumentsWithTimestamp", cols,
			[]interface{}{3, 9, spanner.CommitTimestamp, "Hello World 9 Updated"}),
	})
	return err
}

func query(ctx context.Context, w io.Writer, client *spanner.Client) error {
	stmt := spanner.Statement{SQL: `SELECT UserId, DocumentId, Ts, Contents FROM DocumentsWithTimestamp
		ORDER BY Ts DESC Limit 5`}
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
		var userID, documentID int64
		var ts time.Time
		var contents string
		if err := row.Columns(&userID, &documentID, &ts, &contents); err != nil {
			return err
		}
		fmt.Fprintf(w, "%d %d %s %s\n", userID, documentID, ts, contents)
	}
}

func createTableWithHistoryTable(ctx context.Context, w io.Writer, adminClient *database.DatabaseAdminClient, database string) error {
	op, err := adminClient.UpdateDatabaseDdl(ctx, &adminpb.UpdateDatabaseDdlRequest{
		Database: database,
		Statements: []string{
			`CREATE TABLE Documents(
				UserId INT64 NOT NULL,
				DocumentId INT64 NOT NULL,
				Contents STRING(MAX) NOT NULL
			) PRIMARY KEY(UserId, DocumentId)`,
			`CREATE TABLE DocumentHistory(
				UserId INT64 NOT NULL,
				DocumentId INT64 NOT NULL,
				Ts TIMESTAMP NOT NULL OPTIONS(allow_commit_timestamp=true),
				PreviousContents STRING(MAX)
			) PRIMARY KEY(UserId, DocumentId, Ts), INTERLEAVE IN PARENT Documents ON DELETE NO ACTION`,
		},
	})
	if err != nil {
		return err
	}
	if err := op.Wait(ctx); err != nil {
		return err
	}
	fmt.Fprintf(w, "Created Documents and DocumentHistory tables in database [%s]\n", database)
	return nil
}

func writeWithHistory(ctx context.Context, w io.Writer, client *spanner.Client) error {
	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		DocumentsColumns := []string{"UserId", "DocumentId", "Contents"}
		DocumentHistoryColumns := []string{"UserId", "DocumentId", "Ts", "PreviousContents"}
		txn.BufferWrite([]*spanner.Mutation{
			spanner.InsertOrUpdate("Documents", DocumentsColumns,
				[]interface{}{1, 1, "Hello World 1"}),
			spanner.InsertOrUpdate("Documents", DocumentsColumns,
				[]interface{}{1, 2, "Hello World 2"}),
			spanner.InsertOrUpdate("Documents", DocumentsColumns,
				[]interface{}{1, 3, "Hello World 3"}),
			spanner.InsertOrUpdate("Documents", DocumentsColumns,
				[]interface{}{2, 4, "Hello World 4"}),
			spanner.InsertOrUpdate("Documents", DocumentsColumns,
				[]interface{}{2, 5, "Hello World 5"}),
			spanner.InsertOrUpdate("DocumentHistory", DocumentHistoryColumns,
				[]interface{}{1, 1, spanner.CommitTimestamp, "Hello World 1"}),
			spanner.InsertOrUpdate("DocumentHistory", DocumentHistoryColumns,
				[]interface{}{1, 2, spanner.CommitTimestamp, "Hello World 2"}),
			spanner.InsertOrUpdate("DocumentHistory", DocumentHistoryColumns,
				[]interface{}{1, 3, spanner.CommitTimestamp, "Hello World 3"}),
			spanner.InsertOrUpdate("DocumentHistory", DocumentHistoryColumns,
				[]interface{}{2, 4, spanner.CommitTimestamp, "Hello World 4"}),
			spanner.InsertOrUpdate("DocumentHistory", DocumentHistoryColumns,
				[]interface{}{2, 5, spanner.CommitTimestamp, "Hello World 5"}),
		})
		return nil
	})
	return err
}

func updateWithHistory(ctx context.Context, w io.Writer, client *spanner.Client) error {
	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		// Create anonymous function "getContents" to read the current value of the Contents column for a given row.
		getContents := func(key spanner.Key) (string, error) {
			row, err := txn.ReadRow(ctx, "Documents", key, []string{"Contents"})
			if err != nil {
				return "", err
			}
			var content string
			if err := row.Column(0, &content); err != nil {
				return "", err
			}
			return content, nil
		}
		// Create two string arrays corresponding to the columns in each table.
		DocumentsColumns := []string{"UserId", "DocumentId", "Contents"}
		DocumentHistoryColumns := []string{"UserId", "DocumentId", "Ts", "PreviousContents"}
		// Get row's Contents before updating.
		previousContents, err := getContents(spanner.Key{1, 1})
		if err != nil {
			return err
		}
		// Update row's Contents while saving previous Contents in DocumentHistory table.
		txn.BufferWrite([]*spanner.Mutation{
			spanner.InsertOrUpdate("Documents", DocumentsColumns,
				[]interface{}{1, 1, "Hello World 1 Updated"}),
			spanner.InsertOrUpdate("DocumentHistory", DocumentHistoryColumns,
				[]interface{}{1, 1, spanner.CommitTimestamp, previousContents}),
		})
		previousContents, err = getContents(spanner.Key{1, 3})
		if err != nil {
			return err
		}
		txn.BufferWrite([]*spanner.Mutation{
			spanner.InsertOrUpdate("Documents", DocumentsColumns,
				[]interface{}{1, 3, "Hello World 3 Updated"}),
			spanner.InsertOrUpdate("DocumentHistory", DocumentHistoryColumns,
				[]interface{}{1, 3, spanner.CommitTimestamp, previousContents}),
		})
		previousContents, err = getContents(spanner.Key{2, 5})
		if err != nil {
			return err
		}
		txn.BufferWrite([]*spanner.Mutation{
			spanner.InsertOrUpdate("Documents", DocumentsColumns,
				[]interface{}{2, 5, "Hello World 5 Updated"}),
			spanner.InsertOrUpdate("DocumentHistory", DocumentHistoryColumns,
				[]interface{}{2, 5, spanner.CommitTimestamp, previousContents}),
		})
		return nil
	})
	return err
}

func queryWithHistory(ctx context.Context, w io.Writer, client *spanner.Client) error {
	stmt := spanner.Statement{
		SQL: `SELECT d.UserId, d.DocumentId, d.Contents, dh.Ts, dh.PreviousContents
				FROM Documents d JOIN DocumentHistory dh
				ON dh.UserId = d.UserId AND dh.DocumentId = d.DocumentId
				ORDER BY dh.Ts DESC LIMIT 3`}
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
		var userID, documentID int64
		var ts time.Time
		var contents, previousContents string
		if err := row.Columns(&userID, &documentID, &contents, &ts, &previousContents); err != nil {
			return err
		}
		fmt.Fprintf(w, "%d %d %s %s %s\n", userID, documentID, contents, ts, previousContents)
	}
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
		fmt.Fprintf(os.Stderr, `Usage: go run main.go <command> <database_name>

	Command can be one of: createdatabase, createtablewithtimestamp, write, update, query,
			createtablewithhistorytable, writewithhistory, updatewithhistory, querywithhistory

Examples:
		go run main.go createdatabase projects/my-project/instances/my-instance/databases/example-db
		go run main.go createtablewithtimestamp projects/my-project/instances/my-instance/databases/example-db
		go run main.go write projects/my-project/instances/my-instance/databases/example-db
		go run main.go update projects/my-project/instances/my-instance/databases/example-db
		go run main.go query projects/my-project/instances/my-instance/databases/example-db
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
