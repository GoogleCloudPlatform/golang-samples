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

// [START spanner_restore_backup]

import (
	"context"
	"fmt"
	"io"
	"regexp"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
	adminpb "cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
)

func restoreBackup(ctx context.Context, w io.Writer, db, backupID string) error {
	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return err
	}
	defer adminClient.Close()

	matches := regexp.MustCompile("^(.*)/databases/(.*)$").FindStringSubmatch(db)
	if matches == nil || len(matches) != 3 {
		return fmt.Errorf("Invalid database id %s", db)
	}
	instanceName := matches[1]
	databaseID := matches[2]
	backupName := instanceName + "/backups/" + backupID

	// Start restoring backup to a new database.
	restoreOp, err := adminClient.RestoreDatabase(ctx, &adminpb.RestoreDatabaseRequest{
		Parent:     instanceName,
		DatabaseId: databaseID,
		Source: &adminpb.RestoreDatabaseRequest_Backup{
			Backup: backupName,
		},
	})
	if err != nil {
		return err
	}
	// Wait for restore operation to complete.
	dbObj, err := restoreOp.Wait(ctx)
	if err != nil {
		return err
	}
	// Newly created database has restore information.
	backupInfo := dbObj.RestoreInfo.GetBackupInfo()
	if backupInfo != nil {
		fmt.Fprintf(w, "Source database %s restored from backup %s\n", backupInfo.SourceDatabase, backupInfo.Backup)
	}

	return nil
}

// [END spanner_restore_backup]
