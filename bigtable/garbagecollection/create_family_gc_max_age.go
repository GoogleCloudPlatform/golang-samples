// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package garbagecollection

// [START bigtable_create_family_gc_max_age]
import (
	"context"
	"fmt"
	"io"
	"time"

	admin "cloud.google.com/go/bigtable/admin/apiv2"
	"cloud.google.com/go/bigtable/admin/apiv2/adminpb"
	"google.golang.org/protobuf/types/known/durationpb"
)

func createFamilyGCMaxAge(w io.Writer, projectID, instanceID string, tableName string) error {
	// projectID := "my-project-id"
	// instanceID := "my-instance-id"
	// tableName := "my-table-name"

	ctx := context.Background()

	adminClient, err := admin.NewBigtableTableAdminClient(ctx)
	if err != nil {
		return fmt.Errorf("admin.NewBigtableTableAdminClient: %w", err)
	}
	defer adminClient.Close()

	// Set a garbage collection policy of 5 days.
	maxAge := time.Hour * 24 * 5
	policy := &adminpb.GcRule{
		Rule: &adminpb.GcRule_MaxAge{
			MaxAge: durationpb.New(maxAge),
		},
	}

	columnFamilyName := "cf1"
	req := &adminpb.ModifyColumnFamiliesRequest{
		Name: fmt.Sprintf("projects/%s/instances/%s/tables/%s", projectID, instanceID, tableName),
		Modifications: []*adminpb.ModifyColumnFamiliesRequest_Modification{
			{
				Id: columnFamilyName,
				Mod: &adminpb.ModifyColumnFamiliesRequest_Modification_Create{
					Create: &adminpb.ColumnFamily{
						GcRule: policy,
					},
				},
			},
		},
	}
	if _, err := adminClient.ModifyColumnFamilies(ctx, req); err != nil {
		return fmt.Errorf("ModifyColumnFamilies(%s): %w", columnFamilyName, err)
	}

	fmt.Fprintf(w, "created column family %s with policy: %v\n", columnFamilyName, policy)
	return nil
}

// [END bigtable_create_family_gc_max_age]
