// Copyright 2024 Google LLC
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

// [START spanner_update_backup_schedule]

import (
	"context"
	"fmt"
	"io"
	"time"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
	"cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func updateBackupSchedule(w io.Writer, dbName string, scheduleId string) error {
	ctx := context.Background()

	client, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	// Update a schedule to create backups daily at 3:45 PM, using the database's
	// encryption config, and retained for 48 hours.
	req := databasepb.UpdateBackupScheduleRequest{
		BackupSchedule: &databasepb.BackupSchedule{
			Name: fmt.Sprintf("%s/backupSchedules/%s", dbName, scheduleId),
			Spec: &databasepb.BackupScheduleSpec{
				ScheduleSpec: &databasepb.BackupScheduleSpec_CronSpec{
					CronSpec: &databasepb.CrontabSpec{
						Text: "45 15 * * *",
					},
				},
			},
			RetentionDuration: durationpb.New(48 * time.Hour),
			EncryptionConfig: &databasepb.CreateBackupEncryptionConfig{
				EncryptionType: databasepb.CreateBackupEncryptionConfig_USE_DATABASE_ENCRYPTION,
			},
		},
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: []string{
				"spec.cron_spec.text",
				"retention_duration",
				"encryption_config",
			},
		},
	}

	res, err := client.UpdateBackupSchedule(ctx, &req)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "Updated backup schedule: %s", res)
	return nil
}

// [END spanner_update_backup_schedule]
