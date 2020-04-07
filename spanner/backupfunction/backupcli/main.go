// Copyright 2020 Google LLC
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

// backupcli can be built to run as a CLI function that is intended to be run as a cron job
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
	"github.com/GoogleCloudPlatform/golang-samples/spanner/backupfunction"
)

func main() {
	databaseName := flag.String("databaseName", "", "projects/my-project/instances/my-instance/databases/example-db")
	// Set expire to be 30 days
	expire := flag.Duration("expire", 720*time.Hour, "The time.Duration after which the backup will expire")
	backupPrefix := flag.String("backupPrefix", "backup", "Prefix for backup name, where backup name will be prefix+timestamp")
	awaitCompletion := flag.Bool("awaitCompletion", false, "Boolean: await completion of backup")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: backupfunction -databaseName=<database_name> -expire=<expire duration> -backupPrefix=<backup prefix> -awaitCompletion=<await completion>
	Examples:
	backupfunction -databaseName=projects/my-project/instances/my-instance/databases/example-db 
	backupfunction -databaseName=projects/my-project/instances/my-instance/databases/example-db -expire=6h -backupPrefix=example-backup -awaitCompletion=true`)
		fmt.Println("")
		flag.PrintDefaults()
	}
	flag.Parse()
	// Check that a database has been supplied
	if *databaseName == "" {
		log.Print("ERROR: databaseName cannot be empty")
		flag.Usage()
		os.Exit(2)
	}
	ctx := context.Background()
	client, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		log.Printf("database.NewDatabaseAdminClient: %v", err)
		os.Exit(1)
	}
	op, backupErr := backupfunction.CreateBackup(ctx, os.Stdout, client, *databaseName, *expire, *backupPrefix)
	if backupErr != nil {
		log.Printf("database.NewDatabaseAdminClient: %v", backupErr)
		os.Exit(1)
	}
	if *awaitCompletion {
		log.Println("Awaiting Backup completion")
		backup, completionErr := op.Wait(ctx)
		if completionErr != nil {
			log.Printf("*database.CreateBackupOperation: %v", completionErr)
			os.Exit(1)
		}
		log.Printf("*database.CreateBackupOperation: backup=%v , state= %v", backup.Name, backup.State)
	}
}
