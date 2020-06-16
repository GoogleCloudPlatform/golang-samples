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

// +build ignore

// Command spanner_snippets contains runnable snippet code for Cloud Spanner.
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	. "github.com/GoogleCloudPlatform/golang-samples/spanner/spanner_snippets/spanner"
)

type command func(w io.Writer, db string) error
type backupCommand func(w io.Writer, db, backupID string) error

var (
	commands = map[string]command{
		"write":                           Write,
		"delete":                          Delete,
		"query":                           Query,
		"read":                            Read,
		"update":                          Update,
		"writetransaction":                WriteWithTransaction,
		"querynewcolumn":                  QueryNewColumn,
		"queryindex":                      QueryUsingIndex,
		"readindex":                       ReadUsingIndex,
		"readstoringindex":                ReadStoringIndex,
		"readonlytransaction":             ReadOnlyTransaction,
		"readstaledata":                   ReadStaleData,
		"readbatchdata":                   ReadBatchData,
		"updatewithtimestamp":             UpdateWithTimestamp,
		"querywithtimestamp":              QueryWithTimestamp,
		"writewithtimestamp":              WriteWithTimestamp,
		"querynewtable":                   QueryNewTable,
		"writetodocstable":                WriteToDocumentsTable,
		"updatedocstable":                 UpdateDocumentsTable,
		"querydocstable":                  QueryDocumentsTable,
		"writewithhistory":                WriteWithHistory,
		"updatewithhistory":               UpdateWithHistory,
		"querywithhistory":                QueryWithHistory,
		"writestructdata":                 WriteStructData,
		"querywithstruct":                 QueryWithStruct,
		"querywitharrayofstruct":          QueryWithArrayOfStruct,
		"querywithstructfield":            QueryWithStructField,
		"querywithnestedstructfield":      QueryWithNestedStructField,
		"dmlinsert":                       InsertUsingDML,
		"dmlupdate":                       UpdateUsingDML,
		"dmldelete":                       DeleteUsingDML,
		"dmlwithtimestamp":                UpdateUsingDMLWithTimestamp,
		"dmlwriteread":                    WriteAndReadUsingDML,
		"dmlupdatestruct":                 UpdateUsingDMLStruct,
		"dmlwrite":                        WriteUsingDML,
		"querywithparameter":              QueryWithParameter,
		"dmlwritetxn":                     WriteWithTransactionUsingDML,
		"dmlupdatepart":                   UpdateUsingPartitionedDML,
		"dmldeletepart":                   DeleteUsingPartitionedDML,
		"dmlbatchupdate":                  UpdateUsingBatchDML,
		"writedatatypesdata":              WriteDatatypesData,
		"querywitharray":                  QueryWithArray,
		"querywithbool":                   QueryWithBool,
		"querywithbytes":                  QueryWithBytes,
		"querywithdate":                   QueryWithDate,
		"querywithfloat":                  QueryWithFloat,
		"querywithint":                    QueryWithInt,
		"querywithstring":                 QueryWithString,
		"querywithtimestampparameter":     QueryWithTimestampParameter,
		"querywithqueryoptions":           QueryWithQueryOptions,
		"createclientwithqueryoptions":    CreateClientWithQueryOptions,
		"createdatabase":                  CreateDatabase,
		"addnewcolumn":                    AddNewColumn,
		"addindex":                        AddIndex,
		"addstoringindex":                 AddStoringIndex,
		"addcommittimestamp":              AddCommitTimestamp,
		"createtablewithtimestamp":        CreateTableWithTimestamp,
		"createtablewithdatatypes":        CreateTableWithDatatypes,
		"createtabledocswithtimestamp":    CreateTableDocumentsWithTimestamp,
		"createtabledocswithhistorytable": CreateTableDocumentsWithHistoryTable,
		"listbackupoperations":            ListBackupOperations,
		"listdatabaseoperations":          ListDatabaseOperations,
	}

	backupCommands = map[string]backupCommand{
		"createbackup":  CreateBackup,
		"cancelbackup":  CancelBackup,
		"listbackups":   ListBackups,
		"updatebackup":  UpdateBackup,
		"deletebackup":  DeleteBackup,
		"restorebackup": RestoreBackup,
	}
)

func run(ctx context.Context, w io.Writer, cmd string, db string, backupID string) error {
	// Command that needs a backup ID.
	if backupCmdFn := backupCommands[cmd]; backupCmdFn != nil {
		err := backupCmdFn(w, db, backupID)
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
	err := cmdFn(w, db)
	if err != nil {
		fmt.Fprintf(w, "%s failed with %v", cmd, err)
	}
	return err
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: go run snippet.go <command> <database_name> <backup_id>

	Command can be one of: createdatabase, write, query, read, update,
		writetransaction, addnewcolumn, querynewcolumn, addindex, queryindex, readindex,
		addstoringindex, readstoringindex, readonlytransaction, readstaledata, readbatchdata,
		addcommittimestamp, updatewithtimestamp, querywithtimestamp, createtablewithtimestamp,
		writewithtimestamp, querynewtable, createtabledocswithtimestamp, writetodocstable,
		updatedocstable, querydocstable, createtabledocswithhistorytable, writewithhistory,
		updatewithhistory, querywithhistory, writestructdata, querywithstruct, querywitharrayofstruct,
		querywithstructfield, querywithnestedstructfield, dmlinsert, dmlupdate, dmldelete,
		dmlwithtimestamp, dmlwriteread, dmlwrite, dmlwritetxn, querywithparameter, dmlupdatepart,
		dmldeletepart, dmlbatchupdate, createtablewithdatatypes, writedatatypesdata, querywitharray,
		querywithbool, querywithbytes, querywithdate, querywithfloat, querywithint, querywithstring,
		querywithtimestampparameter, createbackup, listbackups, updatebackup, deletebackup, restorebackup,
		listbackupoperations, listdatabaseoperations, querywithtimestampparameter, querywithqueryoptions,
		createclientwithqueryoptions

Examples:
	go run snippet.go createdatabase projects/my-project/instances/my-instance/databases/example-db
	go run snippet.go write projects/my-project/instances/my-instance/databases/example-db
	go run snippet.go createbackup projects/my-project/instances/my-instance/databases/example-db my-backup
`)
	}

	flag.Parse()
	if len(flag.Args()) < 2 {
		flag.Usage()
		os.Exit(2)
	}

	cmd, db, backupID := flag.Arg(0), flag.Arg(1), flag.Arg(2)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	if err := run(ctx, os.Stdout, cmd, db, backupID); err != nil {
		os.Exit(1)
	}
}
