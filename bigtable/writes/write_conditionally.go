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

// [START bigtable_writes_conditional]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/bigtable"
)

func writeConditionally(w io.Writer, projectID, instanceID string, tableName string) error {
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
	timestamp := bigtable.Now()

	mut := bigtable.NewMutation()
	mut.Set(columnFamilyName, "os_name", timestamp, []byte("android"))

	filter := bigtable.ChainFilters(
		bigtable.FamilyFilter(columnFamilyName),
		bigtable.ColumnFilter("os_build"),
		bigtable.ValueFilter("PQ2A\\..*"))
	conditionalMutation := bigtable.NewCondMutation(filter, mut, nil)

	rowKey := "phone#4c410523#20190501"
	if err := tbl.Apply(ctx, rowKey, conditionalMutation); err != nil {
		return fmt.Errorf("Apply: %v", err)
	}

	fmt.Fprintln(w, "Successfully updated row's os_name")
	return nil
}

// [END bigtable_writes_conditional]
