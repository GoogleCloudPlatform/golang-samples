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
	"flag"
	"fmt"
	"io"
	"os"

	"cloud.google.com/go/spanner"
)

type command func(w io.Writer, client *spanner.Client) error
type backupCommand func(w io.Writer, database, backupID string) error

var (
	commands = map[string]command{
		"write":                           write,
		"delete":                          delete,
		"query":                           query,
		"read":                            read,
		"update":                          update,
		"writetransaction":                writeWithTransaction,
		"querynewcolumn":                  queryNewColumn,
		"queryindex":                      queryUsingIndex,
		"readindex":                       readUsingIndex,
		"readstoringindex":                readStoringIndex,
		"readonlytransaction":             readOnlyTransaction,
		"readstaledata":                   readStaleData,
		"readbatchdata":                   readBatchData,
		"updatewithtimestamp":             updateWithTimestamp,
		"querywithtimestamp":              queryWithTimestamp,
		"writewithtimestamp":              writeWithTimestamp,
		"querynewtable":                   queryNewTable,
		"writetodocstable":                writeToDocumentsTable,
		"updatedocstable":                 updateDocumentsTable,
		"querydocstable":                  queryDocumentsTable,
		"writewithhistory":                writeWithHistory,
		"updatewithhistory":               updateWithHistory,
		"querywithhistory":                queryWithHistory,
		"writestructdata":                 writeStructData,
		"querywithstruct":                 queryWithStruct,
		"querywitharrayofstruct":          queryWithArrayOfStruct,
		"querywithstructfield":            queryWithStructField,
		"querywithnestedstructfield":      queryWithNestedStructField,
		"dmlinsert":                       insertUsingDML,
		"dmlupdate":                       updateUsingDML,
		"dmldelete":                       deleteUsingDML,
		"dmlwithtimestamp":                updateUsingDMLWithTimestamp,
		"dmlwriteread":                    writeAndReadUsingDML,
		"dmlupdatestruct":                 updateUsingDMLStruct,
		"dmlwrite":                        writeUsingDML,
		"querywithparameter":              queryWithParameter,
		"dmlwritetxn":                     writeWithTransactionUsingDML,
		"dmlupdatepart":                   updateUsingPartitionedDML,
		"dmldeletepart":                   deleteUsingPartitionedDML,
		"dmlbatchupdate":                  updateUsingBatchDML,
		"writedatatypesdata":              writeDatatypesData,
		"querywitharray":                  queryWithArray,
		"querywithbool":                   queryWithBool,
		"querywithbytes":                  queryWithBytes,
		"querywithdate":                   queryWithDate,
		"querywithfloat":                  queryWithFloat,
		"querywithint":                    queryWithInt,
		"querywithstring":                 queryWithString,
		"querywithtimestampparameter":     queryWithTimestampParameter,
		"querywithqueryoptions":           queryWithQueryOptions,
		"createclientwithqueryoptions":    createClientWithQueryOptions,
		"createdatabase":                  createDatabase,
		"addnewcolumn":                    addNewColumn,
		"addindex":                        addIndex,
		"addstoringindex":                 addStoringIndex,
		"addcommittimestamp":              addCommitTimestamp,
		"createtablewithtimestamp":        createTableWithTimestamp,
		"createtablewithdatatypes":        createTableWithDatatypes,
		"createtabledocswithtimestamp":    createTableDocumentsWithTimestamp,
		"createtabledocswithhistorytable": createTableDocumentsWithHistoryTable,
		"listbackupoperations":            listBackupOperations,
		"listdatabaseoperations":          listDatabaseOperations,
	}

	backupCommands = map[string]backupCommand{
		"createbackup":  createBackup,
		"cancelbackup":  cancelBackup,
		"listbackups":   listBackups,
		"updatebackup":  updateBackup,
		"deletebackup":  deleteBackup,
		"restorebackup": restoreBackup,
	}
)

func run(w io.Writer, cmd string, db string, backupID string) error {
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
	err := cmdFn(w, dataClient)
	if err != nil {
		fmt.Fprintf(w, "%s failed with %v", cmd, err)
	}
	return err
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: spanner_snippets <command> <database_name> <backup_id>

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
	spanner_snippets createdatabase projects/my-project/instances/my-instance/databases/example-db
	spanner_snippets write projects/my-project/instances/my-instance/databases/example-db
	spanner_snippets createbackup projects/my-project/instances/my-instance/databases/example-db my-backup
`)
	}

	flag.Parse()
	if len(flag.Args()) < 2 {
		flag.Usage()
		os.Exit(2)
	}

	cmd, db, backupID := flag.Arg(0), flag.Arg(1), flag.Arg(2)
	if err := run(os.Stdout, cmd, db, backupID); err != nil {
		os.Exit(1)
	}
}
