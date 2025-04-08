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
	adminpb "cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
	pbt "github.com/golang/protobuf/ptypes/timestamp"
)

func createBackup(ctx context.Context, w io.Writer, db, backupID string, versionTime time.Time) error {
	// versionTime := time.Now().AddDate(0, 0, -1) // one day ago
	matches := regexp.MustCompile("^(.+)/databases/(.+)$").FindStringSubmatch(db)
	if matches == nil || len(matches) != 3 {
		return fmt.Errorf("createBackup: invalid database id %q", db)
	}

	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return fmt.Errorf("createBackup.NewDatabaseAdminClient: %w", err)
	}
	defer adminClient.Close()

	expireTime := time.Now().AddDate(0, 0, 14)
	// Create a backup.
	req := adminpb.CreateBackupRequest{
		Parent:   matches[1],
		BackupId: backupID,
		Backup: &adminpb.Backup{
			Database:    db,
			ExpireTime:  &pbt.Timestamp{Seconds: expireTime.Unix(), Nanos: int32(expireTime.Nanosecond())},
			VersionTime: &pbt.Timestamp{Seconds: versionTime.Unix(), Nanos: int32(versionTime.Nanosecond())},
		},
	}
	op, err := adminClient.CreateBackup(ctx, &req)
	if err != nil {
		return fmt.Errorf("createBackup.CreateBackup: %w", err)
	}
	// Wait for backup operation to complete.
	backup, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("createBackup.Wait: %w", err)
	}

	// Get the name, create time, version time and backup size.
	backupCreateTime := time.Unix(backup.CreateTime.Seconds, int64(backup.CreateTime.Nanos))
	backupVersionTime := time.Unix(backup.VersionTime.Seconds, int64(backup.VersionTime.Nanos))
	fmt.Fprintf(w,
		"Backup %s of size %d bytes was created at %s with version time %s\n",
		backup.Name,
		backup.SizeBytes,
		backupCreateTime.Format(time.RFC3339),
		backupVersionTime.Format(time.RFC3339))
	return nil
}

// [END spanner_create_backup]
