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
	pbts "github.com/golang/protobuf/ptypes"
	adminpb "google.golang.org/genproto/googleapis/spanner/admin/database/v1"
	"google.golang.org/genproto/protobuf/field_mask"
)

func updateBackup(w io.Writer, db, backupID string) error {
	ctx := context.Background()
	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return err
	}
	defer adminClient.Close()

	matches := regexp.MustCompile("^(.*)/databases/(.*)$").FindStringSubmatch(db)
	if matches == nil || len(matches) != 3 {
		return fmt.Errorf("Invalid database id %s", db)
	}
	backupName := matches[1] + "/backups/" + backupID

	backup, err := adminClient.GetBackup(ctx, &adminpb.GetBackupRequest{Name: backupName})
	if err != nil {
		return err
	}

	// Expire time must be within 366 days of the create time of the backup.
	expireTime := time.Unix(backup.CreateTime.Seconds, int64(backup.CreateTime.Nanos)).AddDate(0, 0, 30)
	expirespb, err := pbts.TimestampProto(expireTime)
	if err != nil {
		return err
	}

	_, err = adminClient.UpdateBackup(ctx, &adminpb.UpdateBackupRequest{
		Backup: &adminpb.Backup{
			Name:       backupName,
			ExpireTime: expirespb,
		},
		UpdateMask: &field_mask.FieldMask{Paths: []string{"expire_time"}},
	})
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "Updated backup %s\n", backupID)
	return nil
}

// [END spanner_update_backup]
