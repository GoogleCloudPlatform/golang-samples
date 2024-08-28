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

// [START spanner_list_backups]

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"time"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
	adminpb "cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
	"google.golang.org/api/iterator"
)

func listBackups(ctx context.Context, w io.Writer, db, backupID string) error {
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

	printBackups := func(iter *database.BackupIterator) error {
		for {
			resp, err := iter.Next()
			if err == iterator.Done {
				return nil
			}
			if err != nil {
				return err
			}
			fmt.Fprintf(w, "Backup %s\n", resp.Name)
		}
	}

	var iter *database.BackupIterator
	var filter string
	// List all backups.
	iter = adminClient.ListBackups(ctx, &adminpb.ListBackupsRequest{
		Parent: instanceName,
	})
	if err := printBackups(iter); err != nil {
		return err
	}

	// List all backups that contain a name.
	iter = adminClient.ListBackups(ctx, &adminpb.ListBackupsRequest{
		Parent: instanceName,
		Filter: "name:" + backupID,
	})
	if err := printBackups(iter); err != nil {
		return err
	}

	// List all backups that expire before a timestamp.
	expireTime := time.Now().AddDate(0, 0, 30)
	filter = fmt.Sprintf(`expire_time < "%s"`, expireTime.Format(time.RFC3339))
	iter = adminClient.ListBackups(ctx, &adminpb.ListBackupsRequest{
		Parent: instanceName,
		Filter: filter,
	})
	if err := printBackups(iter); err != nil {
		return err
	}

	// List all backups for a database that contains a name.
	iter = adminClient.ListBackups(ctx, &adminpb.ListBackupsRequest{
		Parent: instanceName,
		Filter: "database:" + db,
	})
	if err := printBackups(iter); err != nil {
		return err
	}

	// List all backups with a size greater than some bytes.
	iter = adminClient.ListBackups(ctx, &adminpb.ListBackupsRequest{
		Parent: instanceName,
		Filter: "size_bytes > 100",
	})
	if err := printBackups(iter); err != nil {
		return err
	}

	// List backups that were created after a timestamp that are also ready.
	createTime := time.Now().AddDate(0, 0, -1)
	filter = fmt.Sprintf(
		`create_time >= "%s" AND state:READY`,
		createTime.Format(time.RFC3339),
	)
	iter = adminClient.ListBackups(ctx, &adminpb.ListBackupsRequest{
		Parent: instanceName,
		Filter: filter,
	})
	if err := printBackups(iter); err != nil {
		return err
	}

	// List backups with pagination.
	request := &adminpb.ListBackupsRequest{
		Parent:   instanceName,
		PageSize: 10,
	}
	for {
		iter = adminClient.ListBackups(ctx, request)
		if err := printBackups(iter); err != nil {
			return err
		}
		pageToken := iter.PageInfo().Token
		if pageToken == "" {
			break
		} else {
			request.PageToken = pageToken
		}
	}

	fmt.Fprintf(w, "Backups listed.\n")
	return nil
}

// [END spanner_list_backups]
