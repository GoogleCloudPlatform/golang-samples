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

// [START spanner_restore_backup_with_encryption_key]

import (
	"context"
	"fmt"
	"io"
	"regexp"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
	adminpb "cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
)

func restoreBackupWithCustomerManagedEncryptionKey(ctx context.Context, w io.Writer, db, backupID, kmsKeyName string) error {
	// db = `projects/<project>/instances/<instance-id>/database/<database-id>`
	// backupID = `my-backup-id`
	// kmsKeyName = `projects/<project>/locations/<location>/keyRings/<key_ring>/cryptoKeys/<kms_key_name>`
	matches := regexp.MustCompile("^(.*)/databases/(.*)$").FindStringSubmatch(db)
	if matches == nil || len(matches) != 3 {
		return fmt.Errorf("restoreBackupWithCustomerManagedEncryptionKey: invalid database id %q", db)
	}
	instanceName := matches[1]
	databaseID := matches[2]
	backupName := instanceName + "/backups/" + backupID

	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return fmt.Errorf("restoreBackupWithCustomerManagedEncryptionKey.NewDatabaseAdminClient: %w", err)
	}
	defer adminClient.Close()

	// Restore a database from a backup using a Customer Managed Encryption Key.
	restoreOp, err := adminClient.RestoreDatabase(ctx, &adminpb.RestoreDatabaseRequest{
		Parent:     instanceName,
		DatabaseId: databaseID,
		Source: &adminpb.RestoreDatabaseRequest_Backup{
			Backup: backupName,
		},
		EncryptionConfig: &adminpb.RestoreDatabaseEncryptionConfig{
			EncryptionType: adminpb.RestoreDatabaseEncryptionConfig_CUSTOMER_MANAGED_ENCRYPTION,
			KmsKeyName:     kmsKeyName,
		},
	})
	if err != nil {
		return fmt.Errorf("restoreBackupWithCustomerManagedEncryptionKey.RestoreDatabase: %w", err)
	}
	// Wait for restore operation to complete.
	restoredDatabase, err := restoreOp.Wait(ctx)
	if err != nil {
		return fmt.Errorf("restoreBackupWithCustomerManagedEncryptionKey.Wait: %w", err)
	}
	// Get the information from the newly restored database.
	backupInfo := restoredDatabase.RestoreInfo.GetBackupInfo()
	fmt.Fprintf(w, "Database %s restored from backup %s using encryption key %s\n",
		backupInfo.SourceDatabase,
		backupInfo.Backup,
		restoredDatabase.EncryptionConfig.KmsKeyName)

	return nil
}

// [END spanner_restore_backup_with_encryption_key]
