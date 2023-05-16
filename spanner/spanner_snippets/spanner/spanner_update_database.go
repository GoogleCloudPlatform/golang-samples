// Copyright 2023 Google LLC
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

// [START spanner_update_database]
import (
	"context"
	"fmt"
	"io"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
	adminpb "cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
	"google.golang.org/genproto/protobuf/field_mask"
)

func updateDatabase(ctx context.Context, w io.Writer, db string) error {
	// db = `projects/<project>/instances/<instance-id>/database/<database-id>`
	// Instantiate database admin client.
	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return fmt.Errorf("updateDatabase.NewDatabaseAdminClient: %v", err)
	}
	defer adminClient.Close()

	// Instantiate the request for performing update database operation.
	op, err := adminClient.UpdateDatabase(ctx, &adminpb.UpdateDatabaseRequest{
		Database: &adminpb.Database{
			Name:                 db,
			EnableDropProtection: true,
		},
		UpdateMask: &field_mask.FieldMask{
			Paths: []string{"enable_drop_protection"},
		},
	})
	if err != nil {
		return fmt.Errorf("updateDatabase.UpdateDatabase: %v", err)
	}

	// Wait for update database operation to complete.
	fmt.Fprintf(w, "Waiting for update database operation to complete [%s]\n", db)
	if _, err := op.Wait(ctx); err != nil {
		return fmt.Errorf("updateDatabase.Wait: %v", err)
	}
	fmt.Fprintf(w, "Updated database [%s]\n", db)
	return nil
}

// [END spanner_update_database]
