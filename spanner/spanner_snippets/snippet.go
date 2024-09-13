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

// Command spanner_snippets contains runnable snippet code for Cloud Spanner.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"

	"cloud.google.com/go/iam/apiv1/iampb"
	database "cloud.google.com/go/spanner/admin/database/apiv1"
	adminpb "cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
	expr "google.golang.org/genproto/googleapis/type/expr"
)

type command func(ctx context.Context, w io.Writer, client *spanner.Client) error
type adminCommand func(ctx context.Context, w io.Writer, adminClient *database.DatabaseAdminClient, database string) error

var (
	commands = map[string]command{
		"write":                    write,
		"read":                     read,
		"readdatawithdatabaserole": read,
		"query":                    query,
		"update":                   update,
		"querynewcolumn":           queryNewColumn,
		"pgquerynewcolumn":         pgQueryNewColumn,
		"querywithparameter":       queryWithParameter,
		"pgqueryparameter":         pgQueryParameter,
		"dmlwrite":                 writeUsingDML,
		"pgdmlwrite":               pgWriteUsingDML,
		"dmlwritetxn":              writeWithTransactionUsingDML,
		"pgdmlwritetxn":            pgWriteWithTransactionUsingDML,
		"readindex":                readUsingIndex,
		"readstoringindex":         readStoringIndex,
		"readonlytransaction":      readOnlyTransaction,
	}

	adminCommands = map[string]adminCommand{
		"createdatabase":         createDatabase,
		"addnewcolumn":           addNewColumn,
		"pgaddnewcolumn":         pgAddNewColumn,
		"addstoringindex":        addStoringIndex,
		"pgaddstoringindex":      pgAddStoringIndex,
		"pgcreatedatabase":       pgCreateDatabase,
		"addanddropdatabaserole": addAndDropDatabaseRole,
		"listdatabaseroles":      listDatabaseRoles,
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

// [START spanner_dml_getting_started_insert]

func writeUsingDML(ctx context.Context, w io.Writer, client *spanner.Client) error {
	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		stmt := spanner.Statement{
			SQL: `INSERT Singers (SingerId, FirstName, LastName) VALUES
				(12, 'Melissa', 'Garcia'),
				(13, 'Russell', 'Morales'),
				(14, 'Jacqueline', 'Long'),
				(15, 'Dylan', 'Shaw')`,
		}
		rowCount, err := txn.Update(ctx, stmt)
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "%d record(s) inserted.\n", rowCount)
		return err
	})
	return err
}

// [END spanner_dml_getting_started_insert]

// [START spanner_query_with_parameter]

func queryWithParameter(ctx context.Context, w io.Writer, client *spanner.Client) error {
	stmt := spanner.Statement{
		SQL: `SELECT SingerId, FirstName, LastName FROM Singers
			WHERE LastName = @lastName`,
		Params: map[string]interface{}{
			"lastName": "Garcia",
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
		var singerID int64
		var firstName, lastName string
		if err := row.Columns(&singerID, &firstName, &lastName); err != nil {
			return err
		}
		fmt.Fprintf(w, "%d %s %s\n", singerID, firstName, lastName)
	}
}

// [END spanner_query_with_parameter]

// [START spanner_dml_getting_started_update]

func writeWithTransactionUsingDML(ctx context.Context, w io.Writer, client *spanner.Client) error {
	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		// getBudget returns the budget for a record with a given albumId and singerId.
		getBudget := func(albumID, singerID int64) (int64, error) {
			key := spanner.Key{albumID, singerID}
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
		// updateBudget updates the budget for a record with a given albumId and singerId.
		updateBudget := func(singerID, albumID, albumBudget int64) error {
			stmt := spanner.Statement{
				SQL: `UPDATE Albums
					SET MarketingBudget = @AlbumBudget
					WHERE SingerId = @SingerId and AlbumId = @AlbumId`,
				Params: map[string]interface{}{
					"SingerId":    singerID,
					"AlbumId":     albumID,
					"AlbumBudget": albumBudget,
				},
			}
			_, err := txn.Update(ctx, stmt)
			return err
		}

		// Transfer the marketing budget from one album to another. By keeping the actions
		// in a single transaction, it ensures the movement is atomic.
		const transferAmt = 200000
		album2Budget, err := getBudget(2, 2)
		if err != nil {
			return err
		}
		// The transaction will only be committed if this condition still holds at the time
		// of commit. Otherwise it will be aborted and the callable will be rerun by the
		// client library.
		if album2Budget >= transferAmt {
			album1Budget, err := getBudget(1, 1)
			if err != nil {
				return err
			}
			if err = updateBudget(1, 1, album1Budget+transferAmt); err != nil {
				return err
			}
			if err = updateBudget(2, 2, album2Budget-transferAmt); err != nil {
				return err
			}
			fmt.Fprintf(w, "Moved %d from Album2's MarketingBudget to Album1's.", transferAmt)
		}
		return nil
	})
	return err
}

// [END spanner_dml_getting_started_update]

func pgCreateDatabase(ctx context.Context, w io.Writer, adminClient *database.DatabaseAdminClient, db string) error {
	matches := regexp.MustCompile("^(.*)/databases/(.*)$").FindStringSubmatch(db)
	if matches == nil || len(matches) != 3 {
		return fmt.Errorf("invalid database id %v", db)
	}
	// Databases with PostgreSQL dialect do not support extra DDL statements in the `CreateDatabase` call.
	op, err := adminClient.CreateDatabase(ctx, &adminpb.CreateDatabaseRequest{
		Parent:          matches[1],
		DatabaseDialect: adminpb.DatabaseDialect_POSTGRESQL,
		// Note that PostgreSQL uses double quotes for quoting identifiers. This also
		// includes database names in the CREATE DATABASE statement.
		CreateStatement: `CREATE DATABASE "` + matches[2] + `"`,
	})
	if err != nil {
		return err
	}
	if _, err := op.Wait(ctx); err != nil {
		return err
	}
	updateReq := &adminpb.UpdateDatabaseDdlRequest{
		Database: db,
		Statements: []string{
			`CREATE TABLE Singers (
				SingerId   bigint NOT NULL PRIMARY KEY,
				FirstName  varchar(1024),
				LastName   varchar(1024),
				SingerInfo bytea
			)`,
			`CREATE TABLE Albums (
				AlbumId      bigint NOT NULL,
				SingerId     bigint NOT NULL REFERENCES Singers (SingerId),
				AlbumTitle   text,
				PRIMARY KEY(SingerId, AlbumId)
			)`,
		},
	}
	opUpdate, err := adminClient.UpdateDatabaseDdl(ctx, updateReq)
	if err != nil {
		return err
	}
	if err := opUpdate.Wait(ctx); err != nil {
		return err
	}
	fmt.Fprintf(w, "Created database [%v]\n", db)
	return nil
}

func pgWriteUsingDML(ctx context.Context, w io.Writer, client *spanner.Client) error {
	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		stmt := spanner.Statement{
			SQL: `INSERT INTO Singers (SingerId, FirstName, LastName) VALUES
				(12, 'Melissa', 'Garcia'),
				(13, 'Russell', 'Morales'),
				(14, 'Jacqueline', 'Long'),
				(15, 'Dylan', 'Shaw')`,
		}
		rowCount, err := txn.Update(ctx, stmt)
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "%d record(s) inserted.\n", rowCount)
		return err
	})
	return err
}

func pgQueryParameter(ctx context.Context, w io.Writer, client *spanner.Client) error {
	stmt := spanner.Statement{
		SQL: `SELECT SingerId, FirstName, LastName FROM Singers
			WHERE LastName = $1`,
		Params: map[string]interface{}{
			"p1": "Garcia",
		},
	}
	type Singers struct {
		SingerID            int64
		FirstName, LastName string
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
		var val Singers
		if err := row.ToStruct(&val); err != nil {
			return err
		}
		fmt.Fprintf(w, "%d %s %s\n", val.SingerID, val.FirstName, val.LastName)
	}
}

func pgAddNewColumn(ctx context.Context, w io.Writer, adminClient *database.DatabaseAdminClient, database string) error {
	op, err := adminClient.UpdateDatabaseDdl(ctx, &adminpb.UpdateDatabaseDdlRequest{
		Database: database,
		Statements: []string{
			"ALTER TABLE Albums ADD COLUMN MarketingBudget bigint",
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

func pgQueryNewColumn(ctx context.Context, w io.Writer, client *spanner.Client) error {
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
		if err := row.ColumnByName("singerid", &singerID); err != nil {
			return err
		}
		if err := row.ColumnByName("albumid", &albumID); err != nil {
			return err
		}
		if err := row.ColumnByName("marketingbudget", &marketingBudget); err != nil {
			return err
		}
		budget := "NULL"
		if marketingBudget.Valid {
			budget = strconv.FormatInt(marketingBudget.Int64, 10)
		}
		fmt.Fprintf(w, "%d %d %s\n", singerID, albumID, budget)
	}
}

func pgWriteWithTransactionUsingDML(ctx context.Context, w io.Writer, client *spanner.Client) error {
	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		// getBudget returns the budget for a record with a given albumId and singerId.
		getBudget := func(albumID, singerID int64) (int64, error) {
			key := spanner.Key{albumID, singerID}
			row, err := txn.ReadRow(ctx, "Albums", key, []string{"MarketingBudget"})
			if err != nil {
				return 0, fmt.Errorf("error reading marketing budget for album_id=%v,singer_id=%v: %w",
					albumID, singerID, err)
			}
			var budget int64
			if err := row.Column(0, &budget); err != nil {
				return 0, fmt.Errorf("error decoding marketing budget for album_id=%v,singer_id=%v: %w",
					albumID, singerID, err)
			}
			return budget, nil
		}
		// updateBudget updates the budget for a record with a given albumId and singerId.
		updateBudget := func(singerID, albumID, albumBudget int64) error {
			stmt := spanner.Statement{
				SQL: `UPDATE Albums
					SET MarketingBudget = $1
					WHERE SingerId = $2 and AlbumId = $3`,
				Params: map[string]interface{}{
					"p1": albumBudget,
					"p2": singerID,
					"p3": albumID,
				},
			}
			_, err := txn.Update(ctx, stmt)
			return err
		}

		// Transfer the marketing budget from one album to another. By keeping the actions
		// in a single transaction, it ensures the movement is atomic.
		const transferAmt = 200000
		album2Budget, err := getBudget(2, 2)
		if err != nil {
			return err
		}
		// The transaction will only be committed if this condition still holds at the time
		// of commit. Otherwise it will be aborted and the callable will be rerun by the
		// client library.
		if album2Budget >= transferAmt {
			album1Budget, err := getBudget(1, 1)
			if err != nil {
				return err
			}
			if err = updateBudget(1, 1, album1Budget+transferAmt); err != nil {
				return err
			}
			if err = updateBudget(2, 2, album2Budget-transferAmt); err != nil {
				return err
			}
			fmt.Fprintf(w, "Moved %d from Album2's MarketingBudget to Album1's.", transferAmt)
		}
		return nil
	})
	return err
}

func pgAddStoringIndex(ctx context.Context, w io.Writer, adminClient *database.DatabaseAdminClient, database string) error {
	op, err := adminClient.UpdateDatabaseDdl(ctx, &adminpb.UpdateDatabaseDdlRequest{
		Database: database,
		Statements: []string{
			"CREATE INDEX AlbumsByAlbumTitle2 ON Albums(AlbumTitle) INCLUDE (MarketingBudget)",
		},
	})
	if err != nil {
		return fmt.Errorf("failed to execute spanner database DDL request: %w", err)
	}
	if err := op.Wait(ctx); err != nil {
		return fmt.Errorf("failed to complete spanner database DDL request: %w", err)
	}
	fmt.Fprintf(w, "Added storing index\n")
	return nil
}

func addAndDropDatabaseRole(ctx context.Context, w io.Writer, adminClient *database.DatabaseAdminClient, db string) error {
	op, err := adminClient.UpdateDatabaseDdl(ctx, &adminpb.UpdateDatabaseDdlRequest{
		Database: db,
		Statements: []string{
			"CREATE ROLE parent",
			"GRANT SELECT ON TABLE Albums TO ROLE parent",
			"CREATE ROLE child",
			"GRANT ROLE parent TO ROLE child",
		},
	})
	if err != nil {
		return err
	}
	if err := op.Wait(ctx); err != nil {
		return err
	}
	fmt.Fprintf(w, "Created roles parent and child and granted privileges\n")
	op, err = adminClient.UpdateDatabaseDdl(ctx, &adminpb.UpdateDatabaseDdlRequest{
		Database: db,
		Statements: []string{
			"REVOKE ROLE parent FROM ROLE child",
			"DROP ROLE child",
		},
	})
	if err != nil {
		return err
	}
	if err := op.Wait(ctx); err != nil {
		return err
	}
	fmt.Fprintf(w, "Revoked privileges and dropped role child\n")
	return nil
}

func enableFineGrainedAccess(ctx context.Context, w io.Writer, adminClient *database.DatabaseAdminClient, db string, iamMember string) error {
	if iamMember == "" {
		return errors.New("IAM member must be specified")
	}
	databaseRole := "parent"
	title := "condition title"
	policy, err := adminClient.GetIamPolicy(ctx, &iampb.GetIamPolicyRequest{
		Resource: db,
		Options: &iampb.GetPolicyOptions{
			// IAM conditions need at least version 3
			RequestedPolicyVersion: 3,
		},
	})
	if err != nil {
		return err
	}

	// IAM conditions need at least version 3
	if policy.Version < 3 {
		policy.Version = 3
	}
	policy.Bindings = append(policy.Bindings, []*iampb.Binding{
		{
			Role:    "roles/spanner.fineGrainedAccessUser",
			Members: []string{iamMember},
		},
		{
			Role:    "roles/spanner.databaseRoleUser",
			Members: []string{iamMember},
			Condition: &expr.Expr{
				Expression: fmt.Sprintf(`resource.name.endsWith("/databaseRoles/%s")`, databaseRole),
				Title:      title,
			},
		},
	}...)
	_, err = adminClient.SetIamPolicy(ctx, &iampb.SetIamPolicyRequest{
		Resource: db,
		Policy:   policy,
	})
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "Enabled fine-grained access in IAM.\n")
	return nil
}

func listDatabaseRoles(ctx context.Context, w io.Writer, adminClient *database.DatabaseAdminClient, db string) error {
	iter := adminClient.ListDatabaseRoles(ctx, &adminpb.ListDatabaseRolesRequest{
		Parent: db,
	})
	rolePrefix := db + "/databaseRoles/"
	for {
		role, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		if !strings.HasPrefix(role.Name, rolePrefix) {
			return fmt.Errorf("Role %v does not have prefix %v", role.Name, rolePrefix)
		}
		fmt.Fprintf(w, "%s\n", strings.TrimPrefix(role.Name, rolePrefix))
	}
	return nil
}

func run(ctx context.Context, w io.Writer, cmd string, db string, arg string) error {
	var databaseRole string
	if cmd == "readdatawithdatabaserole" {
		databaseRole = "parent"
	}

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

	if cmd == "enablefinegrainedaccess" {
		err := enableFineGrainedAccess(ctx, w, adminClient, db, arg)
		if err != nil {
			fmt.Fprintf(w, "%s failed with %v", cmd, err)
		}
		return err
	}

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
		fmt.Fprintf(os.Stderr, `Usage: spanner_snippets <command> <database_name> [iam_member]

	Command can be one of: write, read, readdatawithdatabaserole, query, update,
		querynewcolumn, querywithparameter, dmlwrite, dmlwritetxn, readindex,
		readstoringindex, readonlytransaction, createdatabase, addnewcolumn,
		addstoringindex, addanddropdatabaserole, enablefinegrainedaccess,
		listdatabaseroles, pgcreatedatabase, pgqueryparameter, pgdmlwrite,
		pgaddnewcolumn, pgquerynewcolumn, pgdmlwritetxn, pgaddstoringindex

Examples:
	spanner_snippets createdatabase projects/my-project/instances/my-instance/databases/example-db
	spanner_snippets write projects/my-project/instances/my-instance/databases/example-db
	spanner_snippets enablefinegrainedaccess projects/my-project/instances/my-instance/databases/example-db user:alice@example.com
`)
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
