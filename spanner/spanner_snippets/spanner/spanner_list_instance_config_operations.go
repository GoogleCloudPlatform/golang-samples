// Copyright 2022 Google LLC
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

// [START spanner_list_instance_config_operations]

import (
	"context"
	"fmt"
	"io"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/api/iterator"

	instance "cloud.google.com/go/spanner/admin/instance/apiv1"
	"cloud.google.com/go/spanner/admin/instance/apiv1/instancepb"
)

// listInstanceConfigOperations lists all the custom instance config operations
func listInstanceConfigOperations(w io.Writer, projectID string) error {
	// projectID := "my-project-id"

	ctx := context.Background()
	adminClient, err := instance.NewInstanceAdminClient(ctx)
	if err != nil {
		return err
	}
	defer adminClient.Close()
	iter := adminClient.ListInstanceConfigOperations(ctx, &instancepb.ListInstanceConfigOperationsRequest{
		Parent: fmt.Sprintf("projects/%s", projectID),
		Filter: "(metadata.@type=type.googleapis.com/google.spanner.admin.instance.v1.CreateInstanceConfigMetadata)",
	})
	for {
		resp, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			panic(err)
		}
		metadata := &instancepb.CreateInstanceConfigMetadata{}
		if err := ptypes.UnmarshalAny(resp.Metadata, metadata); err != nil {
			return err
		}
		fmt.Fprintf(w, "List instance config operations %s is %d%% completed.\n",
			metadata.InstanceConfig.Name,
			metadata.Progress.ProgressPercent,
		)
	}
	return nil
}

// [END spanner_list_instance_config_operations]
