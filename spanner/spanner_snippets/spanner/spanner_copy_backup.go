// Copyright 2022 Google LLC
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

// [START spanner_copy_backup]

import (
	"context"
	"fmt"
	"io"
	"time"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
	pbt "github.com/golang/protobuf/ptypes/timestamp"
	adminpb "google.golang.org/genproto/googleapis/spanner/admin/database/v1"
)

// copyBackup copies an existing backup to a given instance in same or different region, or in same or different project.
func copyBackup(w io.Writer, instancePath string, copyBackupId string, sourceBackupPath string) error {
	// instancePath := "projects/my-project/instances/my-instance"
	// copyBackupId := "my-copy-backup"
	// sourceBackupPath := "projects/my-project/instances/my-instance/backups/my-source-backup"

	// Add timeout to context.
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	// Instantiate database admin client.
	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return fmt.Errorf("database.NewDatabaseAdminClient: %v", err)
	}
	defer adminClient.Close()

	expireTime := time.Now().AddDate(0, 0, 14)

	// Instantiate the request for performing copy backup operation.
	copyBackupReq := adminpb.CopyBackupRequest{
		Parent:       instancePath,
		BackupId:     copyBackupId,
		SourceBackup: sourceBackupPath,
		ExpireTime:   &pbt.Timestamp{Seconds: expireTime.Unix(), Nanos: int32(expireTime.Nanosecond())},
	}

	// Start copying the backup.
	copyBackupOp, err := adminClient.CopyBackup(ctx, &copyBackupReq)
	if err != nil {
		return fmt.Errorf("adminClient.CopyBackup: %v", err)
	}

	// Wait for copy backup operation to complete.
	fmt.Fprintf(w, "Waiting for backup copy %s/backups/%s to complete...\n", instancePath, copyBackupId)
	copyBackup, err := copyBackupOp.Wait(ctx)
	if err != nil {
		return fmt.Errorf("copyBackup.Wait: %v", err)
	}

	// Check if long-running copyBackup operation is completed.
	if !copyBackupOp.Done() {
		return fmt.Errorf("backup %v could not be copied to %v", sourceBackupPath, copyBackupId)
	}

	// Get the name, create time, version time and backup size.
	copyBackupCreateTime := time.Unix(copyBackup.CreateTime.Seconds, int64(copyBackup.CreateTime.Nanos))
	copyBackupVersionTime := time.Unix(copyBackup.VersionTime.Seconds, int64(copyBackup.VersionTime.Nanos))
	fmt.Fprintf(w,
		"Backup %s of size %d bytes was created at %s with version time %s\n",
		copyBackup.Name,
		copyBackup.SizeBytes,
		copyBackupCreateTime.Format(time.RFC3339),
		copyBackupVersionTime.Format(time.RFC3339))

	return nil
}

// [END spanner_copy_backup]
