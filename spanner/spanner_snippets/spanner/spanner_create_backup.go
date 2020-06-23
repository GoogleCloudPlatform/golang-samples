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
	"time"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
)

func createBackup(w io.Writer, db, backupID string) error {
	ctx := context.Background()
	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return err
	}
	defer adminClient.Close()

	expireTime := time.Now().AddDate(0, 0, 14)
	// Create a backup.
	op, err := adminClient.StartBackupOperation(ctx, backupID, db, expireTime)
	if err != nil {
		return err
	}
	// Wait for backup operation to complete.
	backup, err := op.Wait(ctx)
	if err != nil {
		return err
	}

	// Get the name, create time and backup size.
	createTime := time.Unix(backup.CreateTime.Seconds, int64(backup.CreateTime.Nanos))
	fmt.Fprintf(w, "Backup %s of size %d bytes was created at %s\n", backup.Name, backup.SizeBytes, createTime.Format(time.RFC3339))
	return nil
}

// [END spanner_create_backup]
