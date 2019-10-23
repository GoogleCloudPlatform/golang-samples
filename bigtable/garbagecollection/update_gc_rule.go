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

package garbagecollection

// [START bigtable_update_gc_rule]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/bigtable"
)

func updateGCRule(w io.Writer, projectID, instanceID string, tableName string) error {
	// projectID := "my-project-id"
	// instanceID := "my-instance-id"
	// tableName := "my-table-name"

	ctx := context.Background()

	adminClient, err := bigtable.NewAdminClient(ctx, projectID, instanceID)
	if err != nil {
		return fmt.Errorf("bigtable.NewAdminClient: %v", err)
	}

	columnFamilyName := "cf1"
	// Update the column family cf1 to update the GC rule.
	policy := bigtable.MaxVersionsPolicy(1)
	if err := adminClient.SetGCPolicy(ctx, tableName, columnFamilyName, policy); err != nil {
		return fmt.Errorf("SetGCPolicy(%s): %v", policy, err)
	}

	fmt.Fprintf(w, "Updated column family %s GC rule with policy: %v\n", columnFamilyName, policy)
	return nil
}

// [END bigtable_update_gc_rule]
