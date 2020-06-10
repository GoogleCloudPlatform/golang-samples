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

// [START spanner_cancel_backup_create]

import (
	"context"
	"fmt"
	"io"
	"time"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
	"google.golang.org/genproto/googleapis/longrunning"
	adminpb "google.golang.org/genproto/googleapis/spanner/admin/database/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func cancelBackup(w io.Writer, db, backupID string) error {
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

	// Cancel backup creation.
	err = adminClient.LROClient.CancelOperation(ctx, &longrunning.CancelOperationRequest{Name: op.Name()})
	if err != nil {
		return err
	}

	// Cancel operations are best effort so either it will complete or be
	// cancelled.
	backup, err := op.Wait(ctx)
	if err != nil {
		if waitStatus, ok := status.FromError(err); !ok || waitStatus.Code() != codes.Canceled {
			return err
		}
	} else {
		// Backup was completed before it could be cancelled so delete the
		// unwanted backup.
		err = adminClient.DeleteBackup(ctx, &adminpb.DeleteBackupRequest{Name: backup.Name})
		if err != nil {
			return err
		}
	}

	fmt.Fprintf(w, "Backup cancelled.\n")
	return nil
}

// [END spanner_cancel_backup_create]
