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

package spanner

// [START spanner_create_backup]

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"time"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
	pbt "github.com/golang/protobuf/ptypes/timestamp"
	adminpb "google.golang.org/genproto/googleapis/spanner/admin/database/v1"
)

func createBackup(w io.Writer, db, backupID string) error {
	matches := regexp.MustCompile("^(.+)/databases/(.+)$").FindStringSubmatch(db)
	if matches == nil || len(matches) != 3 {
		return fmt.Errorf("createBackup: invalid database id %q", db)
	}

	ctx := context.Background()
	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return fmt.Errorf("createBackup.NewDatabaseAdminClient: %v", err)
	}
	defer adminClient.Close()
	dbMetadata, err := adminClient.GetDatabase(ctx, &adminpb.GetDatabaseRequest{Name: db})
	if err != nil {
		return fmt.Errorf("createBackup.GetDatabase: %v", err)
	}

	expireTime := time.Now().AddDate(0, 0, 14)
	// Create a backup.
	req := adminpb.CreateBackupRequest{
		Parent:   matches[1],
		BackupId: backupID,
		Backup: &adminpb.Backup{
			Database:    db,
			ExpireTime:  &pbt.Timestamp{Seconds: expireTime.Unix(), Nanos: int32(expireTime.Nanosecond())},
			VersionTime: dbMetadata.EarliestVersionTime,
		},
	}
	op, err := adminClient.CreateBackup(ctx, &req)
	if err != nil {
		return fmt.Errorf("createBackup.CreateBackup: %v", err)
	}
	// Wait for backup operation to complete.
	backup, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("createBackup.Wait: %v", err)
	}

	// Get the name, create time, version time and backup size.
	createTime := time.Unix(backup.CreateTime.Seconds, int64(backup.CreateTime.Nanos))
	versionTime := time.Unix(dbMetadata.EarliestVersionTime.Seconds, int64(dbMetadata.EarliestVersionTime.Nanos))
	fmt.Fprintf(w,
		"Backup %s of size %d bytes was created at %s with version time %s\n",
		backup.Name,
		backup.SizeBytes,
		createTime.Format(time.RFC3339),
		versionTime.Format(time.RFC3339))
	return nil
}

// [END spanner_create_backup]
