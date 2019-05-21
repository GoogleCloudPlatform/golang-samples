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

// [START bigtable_create_family_gc_nested]
import (
	"cloud.google.com/go/bigtable"
	"context"
	"fmt"
	"io"
	"time"
)

func createFamilyGCNested(w io.Writer, projectID, instanceID string, tableName string) error {
	// projectID := "my-project-id"
	// instanceID := "my-instance-id"
	// tableName := "my-table-name"

	ctx := context.Background()

	adminClient, err := bigtable.NewAdminClient(ctx, projectID, instanceID)
	if err != nil {
		return fmt.Errorf("could not create admin client: %v", err)
	}

	columnFamilyName := "cf5"
	if err := adminClient.CreateColumnFamily(ctx, tableName, columnFamilyName); err != nil {
		return fmt.Errorf("could not create column family %s: %v", columnFamilyName, err)
	}

	// Create a nested GC rule:
	// Drop cells that are either older than the 10 recent versions
	// OR
	// Drop cells that are older than a month AND older than the 2 recent versions
	maxAge := time.Hour * 24 * 5
	maxAgePolicy := bigtable.MaxAgePolicy(maxAge)
	nestedPolicy := bigtable.UnionPolicy(bigtable.MaxVersionsPolicy(10),
		bigtable.IntersectionPolicy(bigtable.MaxVersionsPolicy(2),
			maxAgePolicy))
	if err := adminClient.SetGCPolicy(ctx, tableName, columnFamilyName, nestedPolicy); err != nil {
		return fmt.Errorf("could not set garbage collection policy: %v", err)
	}

	fmt.Fprintf(w, "created column family %s with policy: %v\n", columnFamilyName, nestedPolicy)
	return nil
}

// [END bigtable_create_family_gc_nested]
