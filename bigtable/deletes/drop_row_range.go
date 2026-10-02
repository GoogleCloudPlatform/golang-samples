// Copyright 2026 Google LLC
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

package deletes

// [START bigtable_drop_row_range]

import (
	"context"
	"fmt"
	"io"

	admin "cloud.google.com/go/bigtable/admin/apiv2"
	"cloud.google.com/go/bigtable/admin/apiv2/adminpb"
)

func dropRowRange(w io.Writer, projectID, instanceID, tableName string) error {
	ctx := context.Background()
	adminClient, err := admin.NewBigtableTableAdminClient(ctx)
	if err != nil {
		return fmt.Errorf("admin.NewBigtableTableAdminClient: %w", err)
	}
	defer adminClient.Close()

	// Use a specific row key prefix to drop.
	prefix := "phone#4c410523"
	req := &adminpb.DropRowRangeRequest{
		Name:   fmt.Sprintf("projects/%s/instances/%s/tables/%s", projectID, instanceID, tableName),
		Target: &adminpb.DropRowRangeRequest_RowKeyPrefix{RowKeyPrefix: []byte(prefix)},
	}
	if err := adminClient.DropRowRange(ctx, req); err != nil {
		return fmt.Errorf("adminClient.DropRowRange: %w", err)
	}

	fmt.Fprintf(w, "Successfully dropped row range with prefix: %s\n", prefix)
	return nil
}

// [END bigtable_drop_row_range]
