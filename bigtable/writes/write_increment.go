// Copyright 2019 Google LLC
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

package writes

// [START bigtable_writes_increment]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/bigtable"
)

func writeIncrement(w io.Writer, projectID, instanceID string, tableName string) error {
	// projectID := "my-project-id"
	// instanceID := "my-instance-id"
	// tableName := "mobile-time-series"

	ctx := context.Background()
	client, err := bigtable.NewClient(ctx, projectID, instanceID)
	if err != nil {
		return fmt.Errorf("bigtable.NewAdminClient: %v", err)
	}
	defer client.Close()
	tbl := client.Open(tableName)
	columnFamilyName := "stats_summary"

	increment := bigtable.NewReadModifyWrite()
	increment.Increment(columnFamilyName, "connected_wifi", -1)

	rowKey := "phone#4c410523#20190501"
	if _, err := tbl.ApplyReadModifyWrite(ctx, rowKey, increment); err != nil {
		return fmt.Errorf("ApplyReadModifyWrite: %v", err)
	}

	fmt.Fprintf(w, "Successfully updated row: %s\n", rowKey)
	return nil
}

// [END bigtable_writes_increment]
