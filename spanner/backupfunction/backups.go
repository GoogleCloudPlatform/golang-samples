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

// [START spanner_functions_backup_util]

// Package backupfunction is a Cloud function that can be periodically triggered to create a backup
// for the specified database.
package backupfunction

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
)

// client is a global Spanner client, to avoid initializing a new client for
// every request.
var client *database.DatabaseAdminClient
var clientOnce sync.Once

// PubSubMessage is the payload of a Pub/Sub event.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// Backupfunction is intended to be called by a scheduled cloud event that generates a PubSub message.
// The PubSubMessage can contain the parameters to call the backup function if required
func Backupfunction(ctx context.Context, m PubSubMessage) error {
	clientOnce.Do(func() {
		// Declare a separate err variable to avoid shadowing client.
		var err error
		client, err = database.NewDatabaseAdminClient(context.Background())
		if err != nil {
			log.Printf("database.NewDatabaseAdminClient: %v", err)
			return
		}
		// Alternatively parameters defined below can be passed as part of the event, and can be parsed from PubSubMessage
		databaseName := "projects/my-project/instances/my-instance/databases/example-db" // set the name of the database to be backed up
		// DEFAULT Set expire to be the minimum expire duration of 6 hours
		expire := 6 * time.Hour             //The time.Duration after which the backup will expire
		backupPrefix := "example-db-backup" //Prefix for backup name, where backup name will be prefix+timestamp
		// Logging can be customised as per	https://cloud.google.com/functions/docs/monitoring/logging
		_, err = CreateBackup(ctx, os.Stdout, client, databaseName, expire, backupPrefix)
		if err != nil {
			log.Printf("CreateBackup encountered error: %v", err)
			return
		}
	})
	return nil
}

// CreateBackup calls StartBackupOperation on behalf of main (CLI) or Backupfunction which is executed by Cloud Function
func CreateBackup(ctx context.Context, w io.Writer, adminClient *database.DatabaseAdminClient, databaseName string, expiry time.Duration, backupPrefix string) (lrop *database.CreateBackupOperation, err error) {
	timeNow := time.Now()
	backupID := backupPrefix + strings.Replace(timeNow.Format("20060102150405.000000000"), ".", "", -1)
	fmt.Fprintf(w, "backupID = %s\n", backupID)
	expires := timeNow.Add(expiry)
	op, err := adminClient.StartBackupOperation(ctx, backupID, databaseName, expires)
	if err != nil {
		log.SetOutput(w)
		log.Printf("Create backup operation FAILED for database [%s], set to expire at [%s], backupID = %s\n with ERROR=%s", databaseName, expires.Format(time.RFC3339), backupID, err)
		return nil, err
	}
	fmt.Fprintf(w, "Create backup operation [%s] started for database [%s], set to expire at [%s], backupID = %s\n", op.Name(), databaseName, expires.Format(time.RFC3339), backupID)
	return op, nil
}
