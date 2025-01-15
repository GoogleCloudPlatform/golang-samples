// Copyright 2021 Google LLC
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

// [START spanner_create_backup_with_encryption_key]

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

func createBackupWithCustomerManagedEncryptionKey(ctx context.Context, w io.Writer, db, backupID, kmsKeyName string) error {
	// db = `projects/<project>/instances/<instance-id>/database/<database-id>`
	// backupID = `my-backup-id`
	// kmsKeyName = `projects/<project>/locations/<location>/keyRings/<key_ring>/cryptoKeys/<kms_key_name>`
	matches := regexp.MustCompile("^(.+)/databases/(.+)$").FindStringSubmatch(db)
	if matches == nil || len(matches) != 3 {
		return fmt.Errorf("createBackupWithCustomerManagedEncryptionKey: invalid database id %q", db)
	}
	instanceName := matches[1]

	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return fmt.Errorf("createBackupWithCustomerManagedEncryptionKey.NewDatabaseAdminClient: %w", err)
	}
	defer adminClient.Close()

	expireTime := time.Now().AddDate(0, 0, 14)
	// Create a backup for a database using a Customer Managed Encryption Key
	req := adminpb.CreateBackupRequest{
		Parent:   instanceName,
		BackupId: backupID,
		Backup: &adminpb.Backup{
			Database:   db,
			ExpireTime: &pbt.Timestamp{Seconds: expireTime.Unix(), Nanos: int32(expireTime.Nanosecond())},
		},
		EncryptionConfig: &adminpb.CreateBackupEncryptionConfig{
			KmsKeyName:     kmsKeyName,
			EncryptionType: adminpb.CreateBackupEncryptionConfig_CUSTOMER_MANAGED_ENCRYPTION,
		},
	}
	op, err := adminClient.CreateBackup(ctx, &req)
	if err != nil {
		return fmt.Errorf("createBackupWithCustomerManagedEncryptionKey.CreateBackup: %w", err)
	}
	// Wait for backup operation to complete.
	backup, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("createBackupWithCustomerManagedEncryptionKey.Wait: %w", err)
	}

	// Get the name, create time, backup size and encryption key from the backup.
	backupCreateTime := time.Unix(backup.CreateTime.Seconds, int64(backup.CreateTime.Nanos))
	fmt.Fprintf(w,
		"Backup %s of size %d bytes was created at %s using encryption key %s\n",
		backup.Name,
		backup.SizeBytes,
		backupCreateTime.Format(time.RFC3339),
		kmsKeyName)
	return nil
}

// [END spanner_create_backup_with_encryption_key]
