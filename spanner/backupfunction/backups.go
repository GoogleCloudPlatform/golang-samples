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

// [START spanner_functions_backup_util]

// Package backupfunction is a Cloud function that can be periodically triggered
// to create a backup for the specified database.
package backupfunction

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"sync"
	"time"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
)

// client is a global Spanner client, to avoid initializing a new client for
// every request.
var client *database.DatabaseAdminClient
var clientOnce sync.Once

var (
	validDBPattern = regexp.MustCompile("^projects/(?P<project>[^/]+)/instances/(?P<instance>[^/]+)/databases/(?P<database>[^/]+)$")
)

func parseDatabaseName(db string) (project, instance, database string, err error) {
	matches := validDBPattern.FindStringSubmatch(db)
	if len(matches) == 0 {
		return "", "", "", fmt.Errorf("Failed to parse database name from %q according to pattern %q",
			db, validDBPattern.String())
	}
	return matches[1], matches[2], matches[3], nil
}

// PubSubMessage is the payload of a Pub/Sub event.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// Meta is the payload of the `Data` field.
type Meta struct {
	BackupID string `json:"backupId"`
	Database string `json:"database"`
	Expire   string `json:"expire"`
}

// SpannerCreateBackup is intended to be called by a scheduled cloud event that
// is passed in as a PubSub message. The PubSubMessage can contain the
// parameters to call the backup function if required.
func SpannerCreateBackup(ctx context.Context, m PubSubMessage) error {
	clientOnce.Do(func() {
		// Declare a separate err variable to avoid shadowing client.
		var err error
		client, err = database.NewDatabaseAdminClient(context.Background())
		if err != nil {
			log.Printf("Failed to create an instance of DatabaseAdminClient: %v", err)
			return
		}
	})
	if client == nil {
		return fmt.Errorf("Client should not be nil")
	}

	var meta Meta
	err := json.Unmarshal(m.Data, &meta)
	if err != nil {
		return fmt.Errorf("Failed to parse data %s: %v", string(m.Data), err)
	}
	expire, err := time.ParseDuration(meta.Expire)
	if err != nil {
		return fmt.Errorf("Failed to parse expire duration %s: %v", meta.Expire, err)
	}
	_, err = createBackup(ctx, meta.BackupID, meta.Database, expire)
	if err != nil {
		return err
	}
	return nil
}

// createBackup starts a backup operation but not waiting for its completion.
func createBackup(ctx context.Context, backupID, dbName string, expire time.Duration) (*database.CreateBackupOperation, error) {
	now := time.Now()
	if backupID == "" {
		_, _, dbID, err := parseDatabaseName(dbName)
		if err != nil {
			return nil, fmt.Errorf("Failed to start a backup operation for database [%s]: %v", dbName, err)
		}
		backupID = fmt.Sprintf("schedule-%s-%d", dbID, now.UTC().Unix())
	}
	expireTime := now.Add(expire)
	op, err := client.StartBackupOperation(ctx, backupID, dbName, expireTime)
	if err != nil {
		return nil, fmt.Errorf("Failed to start a backup operation for database [%s], expire time [%s], backupID = [%s] with error = %v", dbName, expireTime.Format(time.RFC3339), backupID, err)
	}
	log.Printf("Create backup operation [%s] started for database [%s], expire time [%s], backupID = [%s]", op.Name(), dbName, expireTime.Format(time.RFC3339), backupID)
	return op, nil
}

// [END spanner_functions_backup_util]
