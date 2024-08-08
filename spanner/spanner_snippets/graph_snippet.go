// Copyright 2024 Google LLC
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

// Command spanner_snippets contains runnable snippet code for Cloud Spanner.
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"time"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
	adminpb "google.golang.org/genproto/googleapis/spanner/admin/database/v1"
)

type command func(ctx context.Context, w io.Writer, client *spanner.Client) error
type adminCommand func(ctx context.Context, w io.Writer, adminClient *database.DatabaseAdminClient, database string) error

var (
	commands = map[string]command{
		"insertdata":                    insert,
		"insertdatawithdml":             insertWithDml,
		"updatedatawithdml":             updateWithDml,
		"updatedatawithgraphqueryindml": updateWithGraphQueryInDml,
		"query":                         query,
		"querywithparameter":            queryWithParameter,
		"deletedata":                    delete,
		"deletedatawithdml":             deleteWithDml,
	}

	adminCommands = map[string]adminCommand{
		"createdatabase":              createDatabase,
		"updateallowcommittimestamps": updateAllowCommitTimestamps,
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
			`CREATE TABLE Person (
				id               INT64 NOT NULL,
				name             STRING(MAX),
				birthday         TIMESTAMP,
				country          STRING(MAX),
				city             STRING(MAX),
			) PRIMARY KEY (id)`,
			`CREATE TABLE Account (
				id               INT64 NOT NULL,
				create_time      TIMESTAMP,
				is_blocked       BOOL,
				nick_name        STRING(MAX),
			) PRIMARY KEY (id)`,
			`CREATE TABLE PersonOwnAccount (
				id               INT64 NOT NULL,
				account_id       INT64 NOT NULL,
				create_time      TIMESTAMP,
				FOREIGN KEY (account_id)
					REFERENCES Account (id)
			) PRIMARY KEY (id, account_id),
			INTERLEAVE IN PARENT Person ON DELETE CASCADE`,
			`CREATE TABLE AccountTransferAccount (
				id               INT64 NOT NULL,
				to_id            INT64 NOT NULL,
				amount           FLOAT64,
				create_time      TIMESTAMP NOT NULL,
				order_number     STRING(MAX),
				FOREIGN KEY (to_id) REFERENCES Account (id)
			) PRIMARY KEY (id, to_id, create_time),
			INTERLEAVE IN PARENT Account ON DELETE CASCADE`,
			`CREATE OR REPLACE PROPERTY GRAPH FinGraph
				NODE TABLES (Account, Person)
				EDGE TABLES (
					PersonOwnAccount
						SOURCE KEY(id) REFERENCES Person(id)
						DESTINATION KEY(account_id) REFERENCES Account(id)
						LABEL Owns,
					AccountTransferAccount
						SOURCE KEY(id) REFERENCES Account(id)
						DESTINATION KEY(to_id) REFERENCES Account(id)
						LABEL Transfers)`,
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

// [START spanner_update_allow_commit_timestamps]

func updateAllowCommitTimestamps(ctx context.Context, w io.Writer, adminClient *database.DatabaseAdminClient, db string) error {
	// List of DDL statements to be applied to the database.
	// Alter the AccountTransferAccount to allow commit timestamps to be set for the create_time column.
	ddl := []string{
		"ALTER TABLE AccountTransferAccount ALTER COLUMN create_time SET OPTIONS (allow_commit_timestamp = true)",
	}
	op, err := adminClient.UpdateDatabaseDdl(ctx, &adminpb.UpdateDatabaseDdlRequest{
		Database:   db,
		Statements: ddl,
	})
	if err != nil {
		return err
	}
	// Wait for the UpdateDatabaseDdl operation to finish.
	if err := op.Wait(ctx); err != nil {
		return fmt.Errorf("waiting for operation to complete failed: %w", err)
	}
	fmt.Fprintf(w, "Altered AccountTransferAccount table with allow_commit_timestamp option on database %v.\n", db)
	return err
}

// [END spanner_update_allow_commit_timestamps]

// [START spanner_insert_graph_data]

func parseTime(rfc3339Time string) time.Time {
	t, _ := time.Parse(time.RFC3339, rfc3339Time)
	return t
}

func insert(ctx context.Context, w io.Writer, client *spanner.Client) error {
	personColumns := []string{"id", "name", "birthday", "country", "city"}
	accountColumns := []string{"id", "create_time", "is_blocked", "nick_name"}
	ownColumns := []string{"id", "account_id", "create_time"}
	transferColumns := []string{"id", "to_id", "amount", "create_time", "order_number"}
	m := []*spanner.Mutation{
		spanner.Insert("Account", accountColumns,
			[]interface{}{7, parseTime("2020-01-10T06:22:20.12Z"), false, "Vacation Fund"}),
		spanner.Insert("Account", accountColumns,
			[]interface{}{16, parseTime("2020-01-27T17:55:09.12Z"), true, "Vacation Fund"}),
		spanner.Insert("Account", accountColumns,
			[]interface{}{20, parseTime("2020-02-18T05:44:20.12Z"), false, "Rainy Day Fund"}),
		spanner.Insert("Person", personColumns,
			[]interface{}{1, "Alex", parseTime("1991-12-21T00:00:00.12Z"), "Australia", " Adelaide"}),
		spanner.Insert("Person", personColumns,
			[]interface{}{2, "Dana", parseTime("1980-10-31T00:00:00.12Z"), "Czech_Republic", "Moravia"}),
		spanner.Insert("Person", personColumns,
			[]interface{}{3, "Lee", parseTime("1986-12-07T00:00:00.12Z"), "India", "Kollam"}),
		spanner.Insert("AccountTransferAccount", transferColumns,
			[]interface{}{7, 16, 300.0, parseTime("2020-08-29T15:28:58.12Z"), "304330008004315"}),
		spanner.Insert("AccountTransferAccount", transferColumns,
			[]interface{}{7, 16, 100.0, parseTime("2020-10-04T16:55:05.12Z"), "304120005529714"}),
		spanner.Insert("AccountTransferAccount", transferColumns,
			[]interface{}{16, 20, 300.0, parseTime("2020-09-25T02:36:14.12Z"), "103650009791820"}),
		spanner.Insert("AccountTransferAccount", transferColumns,
			[]interface{}{20, 7, 500.0, parseTime("2020-10-04T16:55:05.12Z"), "304120005529714"}),
		spanner.Insert("AccountTransferAccount", transferColumns,
			[]interface{}{20, 16, 200.0, parseTime("2020-10-17T03:59:40.12Z"), "302290001255747"}),
		spanner.Insert("PersonOwnAccount", ownColumns,
			[]interface{}{1, 7, parseTime("2020-01-10T06:22:20.12Z")}),
		spanner.Insert("PersonOwnAccount", ownColumns,
			[]interface{}{2, 20, parseTime("2020-01-27T17:55:09.12Z")}),
		spanner.Insert("PersonOwnAccount", ownColumns,
			[]interface{}{3, 16, parseTime("2020-02-18T05:44:20.12Z")}),
	}
	_, err := client.Apply(ctx, m)
	return err
}

// [END spanner_insert_graph_data]

// [START spanner_insert_graph_data_with_dml]

func insertWithDml(ctx context.Context, w io.Writer, client *spanner.Client) error {
	_, err1 := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		stmt := spanner.Statement{
			SQL: `INSERT INTO Account (id, create_time, is_blocked)
            		VALUES
            	    	(1, CAST('2000-08-10 08:18:48.463959-07:52' AS TIMESTAMP), false),
            			(2, CAST('2000-08-12 08:18:48.463959-07:52' AS TIMESTAMP), true)`,
		}
		rowCount, err := txn.Update(ctx, stmt)
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "%d record(s) inserted.\n", rowCount)
		return err
	})

	if err1 != nil {
		return err1
	}

	_, err2 := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		stmt := spanner.Statement{
			SQL: `INSERT INTO AccountTransferAccount (id, to_id, create_time, amount)
					VALUES
						(1, 2, PENDING_COMMIT_TIMESTAMP(), 100),
						(1, 1, PENDING_COMMIT_TIMESTAMP(), 200)`,
		}
		rowCount, err := txn.Update(ctx, stmt)
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "%d record(s) inserted.\n", rowCount)
		return err
	})

	return err2
}

// [END spanner_insert_graph_data_with_dml]

// [START spanner_update_graph_data_with_dml]

func updateWithDml(ctx context.Context, w io.Writer, client *spanner.Client) error {
	_, err1 := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		stmt := spanner.Statement{
			SQL: `UPDATE Account SET is_blocked = false WHERE id = 2`,
		}
		rowCount, err := txn.Update(ctx, stmt)
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "%d Account record(s) updated.\n", rowCount)
		return err
	})

	if err1 != nil {
		return err1
	}

	_, err2 := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		stmt := spanner.Statement{
			SQL: `UPDATE AccountTransferAccount SET amount = 300 WHERE id = 1 AND to_id = 2`,
		}
		rowCount, err := txn.Update(ctx, stmt)
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "%d AccountTransferAccount record(s) updated.\n", rowCount)
		return err
	})

	return err2
}

// [END spanner_update_graph_data_with_dml]

// [START spanner_update_graph_data_with_graph_query_in_dml]

func updateWithGraphQueryInDml(ctx context.Context, w io.Writer, client *spanner.Client) error {
	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		stmt := spanner.Statement{
			SQL: `UPDATE Account SET is_blocked = true 
            	  WHERE id IN {
            	  GRAPH FinGraph 
            	  MATCH (a:Account WHERE a.id = 1)-[:TRANSFERS]->{1,2}(b:Account)
            	  RETURN b.id}`,
		}
		rowCount, err := txn.Update(ctx, stmt)
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "%d Account record(s) updated.\n", rowCount)
		return err
	})

	return err
}

// [END spanner_update_graph_data_with_graph_query_in_dml]

// [START spanner_query_graph_data]

func query(ctx context.Context, w io.Writer, client *spanner.Client) error {
	stmt := spanner.Statement{SQL: `Graph FinGraph 
		 MATCH (a:Person)-[o:Owns]->()-[t:Transfers]->()<-[p:Owns]-(b:Person)
		 RETURN a.name AS sender, b.name AS receiver, t.amount, t.create_time AS transfer_at`}
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
		var sender, receiver string
		var amount float64
		var transfer_at time.Time
		if err := row.Columns(&sender, &receiver, &amount, &transfer_at); err != nil {
			return err
		}
		fmt.Fprintf(w, "%s %s %f %s\n", sender, receiver, amount, transfer_at.Format(time.RFC3339))
	}
}

// [END spanner_query_graph_data]

// [START spanner_query_graph_data_with_parameter]

func queryWithParameter(ctx context.Context, w io.Writer, client *spanner.Client) error {
	stmt := spanner.Statement{
		SQL: `Graph FinGraph 
			  MATCH (a:Person)-[o:Owns]->()-[t:Transfers]->()<-[p:Owns]-(b:Person)
			  WHERE t.amount >= @min
			  RETURN a.name AS sender, b.name AS receiver, t.amount, t.create_time AS transfer_at`,
		Params: map[string]interface{}{
			"min": 500,
		},
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
		var sender, receiver string
		var amount float64
		var transfer_at time.Time
		if err := row.Columns(&sender, &receiver, &amount, &transfer_at); err != nil {
			return err
		}
		fmt.Fprintf(w, "%s %s %f %s\n", sender, receiver, amount, transfer_at.Format(time.RFC3339))
	}
}

// [END spanner_query_graph_data_with_parameter]

// [START spanner_delete_graph_data_with_dml]

func deleteWithDml(ctx context.Context, w io.Writer, client *spanner.Client) error {
	_, err1 := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		stmt := spanner.Statement{SQL: `DELETE FROM AccountTransferAccount WHERE id = 1 AND to_id = 2`}
		rowCount, err := txn.Update(ctx, stmt)
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "%d record(s) deleted.\n", rowCount)
		return nil
	})

	if err1 != nil {
		return err1
	}

	_, err2 := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		stmt := spanner.Statement{SQL: `DELETE FROM Account WHERE id = 2`}
		rowCount, err := txn.Update(ctx, stmt)
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "%d record(s) deleted.\n", rowCount)
		return nil
	})

	return err2
}

// [END spanner_delete_graph_data_with_dml]

// [START spanner_delete_graph_data]

func delete(ctx context.Context, w io.Writer, client *spanner.Client) error {
	m := []*spanner.Mutation{
		// spanner.Key can be used to delete a specific set of rows.
		// Delete the PersonOwnAccount rows with the key values (1,7) and (2,20).
		spanner.Delete("PersonOwnAccount", spanner.Key{1, 7}),
		spanner.Delete("PersonOwnAccount", spanner.Key{2, 20}),

		// spanner.KeyRange can be used to delete rows with a key in a specific range.
		// Delete a range of rows where the key prefix is >=1 and <8
		spanner.Delete("AccountTransferAccount",
			spanner.KeyRange{Start: spanner.Key{1}, End: spanner.Key{8}, Kind: spanner.ClosedOpen}),

		// spanner.AllKeys can be used to delete all the rows in a table.
		// Delete all Account rows, which will also delete the remaining
		// AccountTransferAccount rows since it was defined with ON DELETE CASCADE.
		spanner.Delete("Account", spanner.AllKeys()),

		// Delete remaining Person rows, which will also delete the remaining
		// PersonOwnAccount rows since it was defined with ON DELETE CASCADE.
		spanner.Delete("Person", spanner.AllKeys()),
	}
	_, err := client.Apply(ctx, m)
	return err
}

// [END spanner_delete_graph_data]

func run(ctx context.Context, w io.Writer, cmd string, db string, arg string) error {
	var databaseRole string

	cfg := spanner.ClientConfig{
		DatabaseRole: databaseRole,
	}

	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer adminClient.Close()

	dataClient, err := spanner.NewClientWithConfig(ctx, db, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer dataClient.Close()

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
	err = cmdFn(ctx, w, dataClient)
	if err != nil {
		fmt.Fprintf(w, "%s failed with %v", cmd, err)
	}
	return err
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: graph_snippet <command> <database_name>

	Command can be one of: createdatabase, updateallowcommittimestamps, insertdata,
		insertdatawithdml, updatedatawithdml, updatedatawithgraphqueryindml
		query, querywithparameter, deletedatawithdml, deletedata

	Examples:
		graph_snippet createdatabase projects/my-project/instances/my-instance/databases/example-db
		graph_snippet insertdata projects/my-project/instances/my-instance/databases/example-db
		graph_snippet deletedata projects/my-project/instances/my-instance/databases/example-db`)
	}

	flag.Parse()
	if len(flag.Args()) < 2 || len(flag.Args()) > 3 {
		flag.Usage()
		os.Exit(2)
	}

	cmd, db, arg := flag.Arg(0), flag.Arg(1), flag.Arg(2)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	if err := run(ctx, os.Stdout, cmd, db, arg); err != nil {
		os.Exit(1)
	}
}
