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

func copyBackup(ctx context.Context, w io.Writer, instanceId string, copyBackupId string, sourceBackupId string) error {
	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return fmt.Errorf("database.NewDatabaseAdminClient: %v", err)
	}
	defer adminClient.Close()

	expireTime := time.Now().AddDate(0, 0, 14)

	// Instantiate the request for performing copy backup operation.
	copyBackupReq := adminpb.CopyBackupRequest{
		Parent:       instanceId,
		BackupId:     copyBackupId,
		SourceBackup: sourceBackupId,
		ExpireTime:   &pbt.Timestamp{Seconds: expireTime.Unix(), Nanos: int32(expireTime.Nanosecond())},
	}

	// Start copying the backup.
	copyBackupOp, err := adminClient.CopyBackup(ctx, &copyBackupReq)
	if err != nil {
		return fmt.Errorf("adminClient.CopyBackup: %v", err)
	}

	// Wait for copy backup operation to complete.
	copyBackup, err := copyBackupOp.Wait(ctx)
	if err != nil {
		return fmt.Errorf("copyBackup.Wait: %v", err)
	}

	// Check if long-running copyBackup operation is completed.
	if !copyBackupOp.Done() {
		return fmt.Errorf("backup %v could not be copied to %v", sourceBackupId, copyBackupId)
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
