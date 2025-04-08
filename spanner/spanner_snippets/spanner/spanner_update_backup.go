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

// [START spanner_update_backup]

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"time"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
	adminpb "cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
	"google.golang.org/genproto/protobuf/field_mask"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// updateBackup updates the expiration time of a pending or completed backup.
func updateBackup(w io.Writer, db string, backupID string) error {
	// db := "projects/my-project/instances/my-instance/databases/my-database"
	// backupID := "my-backup"

	// Add timeout to context.
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return err
	}
	defer adminClient.Close()

	matches := regexp.MustCompile("^(.*)/databases/(.*)$").FindStringSubmatch(db)
	if matches == nil || len(matches) != 3 {
		return fmt.Errorf("invalid database id %s", db)
	}
	backupName := matches[1] + "/backups/" + backupID

	// Get the backup instance.
	backup, err := adminClient.GetBackup(ctx, &adminpb.GetBackupRequest{Name: backupName})
	if err != nil {
		return err
	}

	// Expire time must be within 366 days of the create time of the backup.
	maxExpireTime := time.Unix(backup.MaxExpireTime.Seconds, int64(backup.MaxExpireTime.Nanos))
	expireTime := time.Unix(backup.ExpireTime.Seconds, int64(backup.ExpireTime.Nanos)).AddDate(0, 0, 30)

	// Ensure that new expire time is less than the max expire time.
	if expireTime.After(maxExpireTime) {
		expireTime = maxExpireTime
	}
	expireTimepb := timestamppb.New(expireTime)

	// Make the update backup request.
	_, err = adminClient.UpdateBackup(ctx, &adminpb.UpdateBackupRequest{
		Backup: &adminpb.Backup{
			Name:       backupName,
			ExpireTime: expireTimepb,
		},
		UpdateMask: &field_mask.FieldMask{Paths: []string{"expire_time"}},
	})
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "Updated backup %s with expire time %s\n", backupName, expireTime)

	return nil
}

// [END spanner_update_backup]
