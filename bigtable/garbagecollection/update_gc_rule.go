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

	admin "cloud.google.com/go/bigtable/admin/apiv2"
	"cloud.google.com/go/bigtable/admin/apiv2/adminpb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func updateGCRule(w io.Writer, projectID, instanceID string, tableName string) error {
	// projectID := "my-project-id"
	// instanceID := "my-instance-id"
	// tableName := "my-table-name"

	ctx := context.Background()

	adminClient, err := admin.NewBigtableTableAdminClient(ctx)
	if err != nil {
		return fmt.Errorf("admin.NewBigtableTableAdminClient: %w", err)
	}
	defer adminClient.Close()

	columnFamilyName := "cf1"
	// Update the column family cf1 to update the GC rule.
	policy := &adminpb.GcRule{
		Rule: &adminpb.GcRule_MaxNumVersions{
			MaxNumVersions: 1,
		},
	}
	req := &adminpb.ModifyColumnFamiliesRequest{
		Name: fmt.Sprintf("projects/%s/instances/%s/tables/%s", projectID, instanceID, tableName),
		Modifications: []*adminpb.ModifyColumnFamiliesRequest_Modification{
			{
				Id: columnFamilyName,
				Mod: &adminpb.ModifyColumnFamiliesRequest_Modification_Update{
					Update: &adminpb.ColumnFamily{
						GcRule: policy,
					},
				},
				UpdateMask: &fieldmaskpb.FieldMask{
					Paths: []string{"gc_rule"},
				},
			},
		},
	}
	if _, err := adminClient.ModifyColumnFamilies(ctx, req); err != nil {
		return fmt.Errorf("ModifyColumnFamilies(%s): %w", columnFamilyName, err)
	}

	fmt.Fprintf(w, "Updated column family %s GC rule with policy: %v\n", columnFamilyName, policy)
	return nil
}

// [END bigtable_update_gc_rule]
