// Copyright 2018 Google Inc. All rights reserved.
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
		"write":                      write,
		"query":                      query,
		"read":                       read,
		"update":                     update,
		"writetransaction":           writeWithTransaction,
		"querynewcolumn":             queryNewColumn,
		"queryindex":                 queryUsingIndex,
		"readindex":                  readUsingIndex,
		"readstoringindex":           readStoringIndex,
		"readonlytransaction":        readOnlyTransaction,
		"readstaledata":              readStaleData,
		"readbatchdata":              readBatchData,
		"updatewithtimestamp":        updateWithTimestamp,
		"querywithtimestamp":         queryWithTimestamp,
		"writewithtimestamp":         writeWithTimestamp,
		"querynewtable":              queryNewTable,
		"writetodocstable":           writeToDocumentsTable,
		"updatedocstable":            updateDocumentsTable,
		"querydocstable":             queryDocumentsTable,
		"writewithhistory":           writeWithHistory,
		"updatewithhistory":          updateWithHistory,
		"querywithhistory":           queryWithHistory,
		"writestructdata":            writeStructData,
		"querywithstruct":            queryWithStruct,
		"querywitharrayofstruct":     queryWithArrayOfStruct,
		"querywithstructfield":       queryWithStructField,
		"querywithnestedstructfield": queryWithNestedStructField,
	}

	adminCommands = map[string]adminCommand{
		"createdatabase":                  createDatabase,
		"addnewcolumn":                    addNewColumn,
		"addindex":                        addIndex,
		"addstoringindex":                 addStoringIndex,
		"addcommittimestamp":              addCommitTimestamp,
		"createtablewithtimestamp":        createTableWithTimestamp,
		"createtabledocswithtimestamp":    createTableDocumentsWithTimestamp,
		"createtabledocswithhistorytable": createTableDocumentsWithHistoryTable,
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
				SingerId     INT64 NOT NULL,
				AlbumId      INT64 NOT NULL,
				AlbumTitle   STRING(MAX)
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

// [START spanner_create_table_with_timestamp_column]

func createTableWithTimestamp(ctx context.Context, w io.Writer, adminClient *database.DatabaseAdminClient, database string) error {
	op, err := adminClient.UpdateDatabaseDdl(ctx, &adminpb.UpdateDatabaseDdlRequest{
		Database: database,
		Statements: []string{
			`CREATE TABLE Performances (
				SingerId        INT64 NOT NULL,
				VenueId         INT64 NOT NULL,
				EventDate       Date,
				Revenue         INT64,
				LastUpdateTime  TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true)
			) PRIMARY KEY (SingerId, VenueId, EventDate),
			INTERLEAVE IN PARENT Singers ON DELETE CASCADE`,
		},
	})
	if err != nil {
		return err
	}
	if err := op.Wait(ctx); err != nil {
		return err
	}
	fmt.Fprintf(w, "Created Performances table in database [%s]\n", database)
	return nil
}

// [END spanner_create_table_with_timestamp_column]

func createTableDocumentsWithTimestamp(ctx context.Context, w io.Writer, adminClient *database.DatabaseAdminClient, database string) error {
	op, err := adminClient.UpdateDatabaseDdl(ctx, &adminpb.UpdateDatabaseDdlRequest{
		Database: database,
		Statements: []string{
			`CREATE TABLE DocumentsWithTimestamp(
				UserId INT64 NOT NULL,
				DocumentId INT64 NOT NULL,
			    Timestamp TIMESTAMP NOT NULL OPTIONS(allow_commit_timestamp=true),
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

func createTableDocumentsWithHistoryTable(ctx context.Context, w io.Writer, adminClient *database.DatabaseAdminClient, database string) error {
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
				Timestamp TIMESTAMP NOT NULL OPTIONS(allow_commit_timestamp=true),
				PreviousContents STRING(MAX)
			) PRIMARY KEY(UserId, DocumentId, Timestamp), INTERLEAVE IN PARENT Documents ON DELETE NO ACTION`,
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

// [START spanner_insert_data_with_timestamp_column]

func writeWithTimestamp(ctx context.Context, w io.Writer, client *spanner.Client) error {
	performanceColumns := []string{"SingerId", "VenueId", "EventDate", "Revenue", "LastUpdateTime"}
	m := []*spanner.Mutation{
		spanner.InsertOrUpdate("Performances", performanceColumns, []interface{}{1, 4, "2017-10-05", 11000, spanner.CommitTimestamp}),
		spanner.InsertOrUpdate("Performances", performanceColumns, []interface{}{1, 19, "2017-11-02", 15000, spanner.CommitTimestamp}),
		spanner.InsertOrUpdate("Performances", performanceColumns, []interface{}{2, 42, "2017-12-23", 7000, spanner.CommitTimestamp}),
	}
	_, err := client.Apply(ctx, m)
	return err
}

// [END spanner_insert_data_with_timestamp_column]

// [START spanner_write_data_for_struct_queries]

func writeStructData(ctx context.Context, w io.Writer, client *spanner.Client) error {
	singerColumns := []string{"SingerId", "FirstName", "LastName"}
	m := []*spanner.Mutation{
		spanner.InsertOrUpdate("Singers", singerColumns, []interface{}{6, "Elena", "Campbell"}),
		spanner.InsertOrUpdate("Singers", singerColumns, []interface{}{7, "Gabriel", "Wright"}),
		spanner.InsertOrUpdate("Singers", singerColumns, []interface{}{8, "Benjamin", "Martinez"}),
		spanner.InsertOrUpdate("Singers", singerColumns, []interface{}{9, "Hannah", "Harris"}),
	}
	_, err := client.Apply(ctx, m)
	return err
}

// [END spanner_write_data_for_struct_queries]

func queryWithStruct(ctx context.Context, w io.Writer, client *spanner.Client) error {

	// [START spanner_create_struct_with_data]

	type nameStruct struct {
		FirstName string
		LastName  string
	}
	var singerInfo = nameStruct{"Elena", "Campbell"}

	// [END spanner_create_struct_with_data]

	// [START spanner_query_data_with_struct]

	stmt := spanner.Statement{
		SQL: `SELECT SingerId FROM SINGERS
				WHERE (FirstName, LastName) = @singerinfo`,
		Params: map[string]interface{}{"singerinfo": singerInfo},
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
		var singerID int64
		if err := row.Columns(&singerID); err != nil {
			return err
		}
		fmt.Fprintf(w, "%d\n", singerID)
	}

	// [END spanner_query_data_with_struct]
}

func queryWithArrayOfStruct(ctx context.Context, w io.Writer, client *spanner.Client) error {

	// [START spanner_create_user_defined_struct]

	type nameType struct {
		FirstName string
		LastName  string
	}

	// [END spanner_create_user_defined_struct]

	// [START spanner_create_array_of_struct_with_data]

	var bandMembers = []nameType{
		{"Elena", "Campbell"},
		{"Gabriel", "Wright"},
		{"Benjamin", "Martinez"},
	}

	// [END spanner_create_array_of_struct_with_data]

	// [START spanner_query_data_with_array_of_struct]

	stmt := spanner.Statement{
		SQL: `SELECT SingerId FROM SINGERS
			WHERE STRUCT<FirstName STRING, LastName STRING>(FirstName, LastName)
			IN UNNEST(@names)`,
		Params: map[string]interface{}{"names": bandMembers},
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
		var singerID int64
		if err := row.Columns(&singerID); err != nil {
			return err
		}
		fmt.Fprintf(w, "%d\n", singerID)
	}

	// [END spanner_query_data_with_array_of_struct]
}

// [START spanner_field_access_on_struct_parameters]

func queryWithStructField(ctx context.Context, w io.Writer, client *spanner.Client) error {
	type structParam struct {
		FirstName string
		LastName  string
	}
	var singerInfo = structParam{"Elena", "Campbell"}
	stmt := spanner.Statement{
		SQL: `SELECT SingerId FROM SINGERS
			WHERE FirstName = @name.FirstName`,
		Params: map[string]interface{}{"name": singerInfo},
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
		var singerID int64
		if err := row.Columns(&singerID); err != nil {
			return err
		}
		fmt.Fprintf(w, "%d\n", singerID)
	}
}

// [END spanner_field_access_on_struct_parameters]

// [START spanner_field_access_on_nested_struct_parameters]

func queryWithNestedStructField(ctx context.Context, w io.Writer, client *spanner.Client) error {
	type nameType struct {
		FirstName string
		LastName  string
	}
	type songInfoStruct struct {
		SongName    string
		ArtistNames []nameType
	}
	var songInfo = songInfoStruct{
		SongName: "Imagination",
		ArtistNames: []nameType{
			{FirstName: "Elena", LastName: "Campbell"},
			{FirstName: "Hannah", LastName: "Harris"},
		},
	}
	stmt := spanner.Statement{
		SQL: `SELECT SingerId, @songinfo.SongName FROM Singers
			WHERE STRUCT<FirstName STRING, LastName STRING>(FirstName, LastName)
			IN UNNEST(@songinfo.ArtistNames)`,
		Params: map[string]interface{}{"songinfo": songInfo},
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
		var singerID int64
		var songName string
		if err := row.Columns(&singerID, &songName); err != nil {
			return err
		}
		fmt.Fprintf(w, "%d %s\n", singerID, songName)
	}
}

// [END spanner_field_access_on_nested_struct_parameters]

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

// [START spanner_add_timestamp_column]

func addCommitTimestamp(ctx context.Context, w io.Writer, adminClient *database.DatabaseAdminClient, database string) error {
	op, err := adminClient.UpdateDatabaseDdl(ctx, &adminpb.UpdateDatabaseDdlRequest{
		Database: database,
		Statements: []string{
			"ALTER TABLE Albums ADD COLUMN LastUpdateTime TIMESTAMP " +
				"OPTIONS (allow_commit_timestamp=true)",
		},
	})
	if err != nil {
		return err
	}
	if err := op.Wait(ctx); err != nil {
		return err
	}
	fmt.Fprintf(w, "Added LastUpdateTime as a commit timestamp column in Albums table\n")
	return nil
}

// [END spanner_add_timestamp_column]

// [START spanner_update_data_with_timestamp_column]

func updateWithTimestamp(ctx context.Context, w io.Writer, client *spanner.Client) error {
	cols := []string{"SingerId", "AlbumId", "MarketingBudget", "LastUpdateTime"}
	_, err := client.Apply(ctx, []*spanner.Mutation{
		spanner.Update("Albums", cols, []interface{}{1, 1, 1000000, spanner.CommitTimestamp}),
		spanner.Update("Albums", cols, []interface{}{2, 2, 750000, spanner.CommitTimestamp}),
	})
	return err
}

// [END spanner_update_data_with_timestamp_column]

// [START spanner_query_data_with_timestamp_column]

func queryWithTimestamp(ctx context.Context, w io.Writer, client *spanner.Client) error {
	stmt := spanner.Statement{
		SQL: `SELECT SingerId, AlbumId, MarketingBudget, LastUpdateTime
	    FROM Albums ORDER BY LastUpdateTime DESC`}
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
		var lastUpdateTime spanner.NullTime
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
		if err := row.ColumnByName("LastUpdateTime", &lastUpdateTime); err != nil {
			return err
		}
		timestamp := "NULL"
		if lastUpdateTime.Valid {
			timestamp = lastUpdateTime.String()
		}
		fmt.Fprintf(w, "%d %d %s %s\n", singerID, albumID, budget, timestamp)
	}
}

// [END spanner_query_data_with_timestamp_column]

func queryNewTable(ctx context.Context, w io.Writer, client *spanner.Client) error {
	stmt := spanner.Statement{
		SQL: `SELECT SingerId, VenueId, EventDate, Revenue, LastUpdateTime FROM Performances
		  ORDER BY LastUpdateTime DESC`}
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
		var singerID, venueID int64
		var revenue spanner.NullInt64
		var eventDate, lastUpdateTime time.Time
		if err := row.ColumnByName("SingerId", &singerID); err != nil {
			return err
		}
		if err := row.ColumnByName("VenueId", &venueID); err != nil {
			return err
		}
		if err := row.ColumnByName("EventDate", &eventDate); err != nil {
			return err
		}
		if err := row.ColumnByName("Revenue", &revenue); err != nil {
			return err
		}
		currentRevenue := "NULL"
		if revenue.Valid {
			currentRevenue = strconv.FormatInt(revenue.Int64, 10)
		}
		if err := row.ColumnByName("LastUpdateTime", &lastUpdateTime); err != nil {
			return err
		}

		fmt.Fprintf(w, "%d %d %s %s %s\n", singerID, venueID, eventDate, currentRevenue, lastUpdateTime)
	}
}

func writeToDocumentsTable(ctx context.Context, w io.Writer, client *spanner.Client) error {
	documentsColumns := []string{"UserId", "DocumentId", "Timestamp", "Contents"}
	m := []*spanner.Mutation{
		spanner.InsertOrUpdate("DocumentsWithTimestamp", documentsColumns,
			[]interface{}{1, 1, spanner.CommitTimestamp, "Hello World 1"}),
		spanner.InsertOrUpdate("DocumentsWithTimestamp", documentsColumns,
			[]interface{}{1, 2, spanner.CommitTimestamp, "Hello World 2"}),
		spanner.InsertOrUpdate("DocumentsWithTimestamp", documentsColumns,
			[]interface{}{1, 3, spanner.CommitTimestamp, "Hello World 3"}),
		spanner.InsertOrUpdate("DocumentsWithTimestamp", documentsColumns,
			[]interface{}{2, 4, spanner.CommitTimestamp, "Hello World 4"}),
		spanner.InsertOrUpdate("DocumentsWithTimestamp", documentsColumns,
			[]interface{}{2, 5, spanner.CommitTimestamp, "Hello World 5"}),
		spanner.InsertOrUpdate("DocumentsWithTimestamp", documentsColumns,
			[]interface{}{3, 6, spanner.CommitTimestamp, "Hello World 6"}),
		spanner.InsertOrUpdate("DocumentsWithTimestamp", documentsColumns,
			[]interface{}{3, 7, spanner.CommitTimestamp, "Hello World 7"}),
		spanner.InsertOrUpdate("DocumentsWithTimestamp", documentsColumns,
			[]interface{}{3, 8, spanner.CommitTimestamp, "Hello World 8"}),
		spanner.InsertOrUpdate("DocumentsWithTimestamp", documentsColumns,
			[]interface{}{3, 9, spanner.CommitTimestamp, "Hello World 9"}),
		spanner.InsertOrUpdate("DocumentsWithTimestamp", documentsColumns,
			[]interface{}{3, 10, spanner.CommitTimestamp, "Hello World 10"}),
	}
	_, err := client.Apply(ctx, m)
	return err
}

func updateDocumentsTable(ctx context.Context, w io.Writer, client *spanner.Client) error {
	cols := []string{"UserId", "DocumentId", "Timestamp", "Contents"}
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

func queryDocumentsTable(ctx context.Context, w io.Writer, client *spanner.Client) error {
	stmt := spanner.Statement{SQL: `SELECT UserId, DocumentId, Timestamp, Contents FROM DocumentsWithTimestamp
		ORDER BY Timestamp DESC Limit 5`}
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
		var timestamp time.Time
		var contents string
		if err := row.Columns(&userID, &documentID, &timestamp, &contents); err != nil {
			return err
		}
		fmt.Fprintf(w, "%d %d %s %s\n", userID, documentID, timestamp, contents)
	}
}

func writeWithHistory(ctx context.Context, w io.Writer, client *spanner.Client) error {
	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		documentsColumns := []string{"UserId", "DocumentId", "Contents"}
		documentHistoryColumns := []string{"UserId", "DocumentId", "Timestamp", "PreviousContents"}
		txn.BufferWrite([]*spanner.Mutation{
			spanner.InsertOrUpdate("Documents", documentsColumns,
				[]interface{}{1, 1, "Hello World 1"}),
			spanner.InsertOrUpdate("Documents", documentsColumns,
				[]interface{}{1, 2, "Hello World 2"}),
			spanner.InsertOrUpdate("Documents", documentsColumns,
				[]interface{}{1, 3, "Hello World 3"}),
			spanner.InsertOrUpdate("Documents", documentsColumns,
				[]interface{}{2, 4, "Hello World 4"}),
			spanner.InsertOrUpdate("Documents", documentsColumns,
				[]interface{}{2, 5, "Hello World 5"}),
			spanner.InsertOrUpdate("DocumentHistory", documentHistoryColumns,
				[]interface{}{1, 1, spanner.CommitTimestamp, "Hello World 1"}),
			spanner.InsertOrUpdate("DocumentHistory", documentHistoryColumns,
				[]interface{}{1, 2, spanner.CommitTimestamp, "Hello World 2"}),
			spanner.InsertOrUpdate("DocumentHistory", documentHistoryColumns,
				[]interface{}{1, 3, spanner.CommitTimestamp, "Hello World 3"}),
			spanner.InsertOrUpdate("DocumentHistory", documentHistoryColumns,
				[]interface{}{2, 4, spanner.CommitTimestamp, "Hello World 4"}),
			spanner.InsertOrUpdate("DocumentHistory", documentHistoryColumns,
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
		documentsColumns := []string{"UserId", "DocumentId", "Contents"}
		documentHistoryColumns := []string{"UserId", "DocumentId", "Timestamp", "PreviousContents"}
		// Get row's Contents before updating.
		previousContents, err := getContents(spanner.Key{1, 1})
		if err != nil {
			return err
		}
		// Update row's Contents while saving previous Contents in DocumentHistory table.
		txn.BufferWrite([]*spanner.Mutation{
			spanner.InsertOrUpdate("Documents", documentsColumns,
				[]interface{}{1, 1, "Hello World 1 Updated"}),
			spanner.InsertOrUpdate("DocumentHistory", documentHistoryColumns,
				[]interface{}{1, 1, spanner.CommitTimestamp, previousContents}),
		})
		previousContents, err = getContents(spanner.Key{1, 3})
		if err != nil {
			return err
		}
		txn.BufferWrite([]*spanner.Mutation{
			spanner.InsertOrUpdate("Documents", documentsColumns,
				[]interface{}{1, 3, "Hello World 3 Updated"}),
			spanner.InsertOrUpdate("DocumentHistory", documentHistoryColumns,
				[]interface{}{1, 3, spanner.CommitTimestamp, previousContents}),
		})
		previousContents, err = getContents(spanner.Key{2, 5})
		if err != nil {
			return err
		}
		txn.BufferWrite([]*spanner.Mutation{
			spanner.InsertOrUpdate("Documents", documentsColumns,
				[]interface{}{2, 5, "Hello World 5 Updated"}),
			spanner.InsertOrUpdate("DocumentHistory", documentHistoryColumns,
				[]interface{}{2, 5, spanner.CommitTimestamp, previousContents}),
		})
		return nil
	})
	return err
}

func queryWithHistory(ctx context.Context, w io.Writer, client *spanner.Client) error {
	stmt := spanner.Statement{
		SQL: `SELECT d.UserId, d.DocumentId, d.Contents, dh.Timestamp, dh.PreviousContents
				FROM Documents d JOIN DocumentHistory dh
				ON dh.UserId = d.UserId AND dh.DocumentId = d.DocumentId
				ORDER BY dh.Timestamp DESC LIMIT 3`}
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
		var timestamp time.Time
		var contents, previousContents string
		if err := row.Columns(&userID, &documentID, &contents, &timestamp, &previousContents); err != nil {
			return err
		}
		fmt.Fprintf(w, "%d %d %s %s %s\n", userID, documentID, contents, timestamp, previousContents)
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
		fmt.Fprintf(os.Stderr, `Usage: spanner_snippets <command> <database_name>

	Command can be one of: createdatabase, write, query, read, update,
		writetransaction, addnewcolumn, querynewcolumn, addindex, queryindex, readindex,
		addstoringindex, readstoringindex, readonlytransaction, readstaledata, readbatchdata,
		addcommittimestamp, updatewithtimestamp, querywithtimestamp, createtablewithtimestamp,
		writewithtimestamp, querynewtable, createtabledocswithtimestamp, writetodocstable,
		updatedocstable, querydocstable, createtabledocswithhistorytable, writewithhistory,
		updatewithhistory, querywithhistory, writestructdata, querywithstruct, querywitharrayofstruct,
		querywithstructfield, querywithnestedstructfield

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
