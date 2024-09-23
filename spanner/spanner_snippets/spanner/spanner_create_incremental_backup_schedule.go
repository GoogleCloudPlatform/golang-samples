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

// [START spanner_create_incremental_backup_schedule]

import (
	"context"
	"fmt"
	"io"
	"time"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
	"cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
	"google.golang.org/protobuf/types/known/durationpb"
)

func createIncrementalBackupSchedule(w io.Writer, dbName string, scheduleId string) error {
	ctx := context.Background()

	client, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	// Create a schedule to create incremental backups daily at 12:30 AM, using
	// Google-managed encryption, and retained for 24 hours.
	req := databasepb.CreateBackupScheduleRequest{
		Parent:           dbName,
		BackupScheduleId: scheduleId,
		BackupSchedule: &databasepb.BackupSchedule{
			Spec: &databasepb.BackupScheduleSpec{
				ScheduleSpec: &databasepb.BackupScheduleSpec_CronSpec{
					CronSpec: &databasepb.CrontabSpec{
						Text: "30 12 * * *",
					},
				},
			},
			RetentionDuration: durationpb.New(24 * time.Hour),
			EncryptionConfig: &databasepb.CreateBackupEncryptionConfig{
				EncryptionType: databasepb.CreateBackupEncryptionConfig_GOOGLE_DEFAULT_ENCRYPTION,
			},
			BackupTypeSpec: &databasepb.BackupSchedule_IncrementalBackupSpec{},
		},
	}

	res, err := client.CreateBackupSchedule(ctx, &req)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "Created incremental backup schedule: %s", res)
	return nil
}

// [END spanner_create_incremental_backup_schedule]
