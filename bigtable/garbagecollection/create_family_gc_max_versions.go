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

// [START bigtable_create_family_gc_max_versions]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/bigtable"
)

func createFamilyGCMaxVersions(w io.Writer, projectID, instanceID string, tableName string) error {
	// projectID := "my-project-id"
	// instanceID := "my-instance-id"
	// tableName := "my-table-name"

	ctx := context.Background()

	adminClient, err := bigtable.NewAdminClient(ctx, projectID, instanceID)
	if err != nil {
		return fmt.Errorf("bigtable.NewAdminClient: %v", err)
	}

	columnFamilyName := "cf2"
	if err := adminClient.CreateColumnFamily(ctx, tableName, columnFamilyName); err != nil {
		return fmt.Errorf("CreateColumnFamily(%s): %v", columnFamilyName, err)
	}

	// Set a garbage collection policy of 2 versions.
	policy := bigtable.MaxVersionsPolicy(2)
	if err := adminClient.SetGCPolicy(ctx, tableName, columnFamilyName, policy); err != nil {
		return fmt.Errorf("SetGCPolicy(%s): %v", policy, err)
	}

	fmt.Fprintf(w, "created column family %s with policy: %v\n", columnFamilyName, policy)
	return nil
}

// [END bigtable_create_family_gc_max_versions]
