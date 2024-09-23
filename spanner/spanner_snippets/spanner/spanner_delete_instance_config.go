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

// [START spanner_delete_instance_config]

import (
	"context"
	"fmt"
	"io"

	instance "cloud.google.com/go/spanner/admin/instance/apiv1"
	"cloud.google.com/go/spanner/admin/instance/apiv1/instancepb"
)

// deleteInstanceConfig deletes the custom spanner instance config
func deleteInstanceConfig(w io.Writer, projectID, userConfigID string) error {
	// projectID := "my-project-id"
	// userConfigID := "custom-config", custom config names must start with the prefix “custom-”.

	ctx := context.Background()
	adminClient, err := instance.NewInstanceAdminClient(ctx)
	if err != nil {
		return err
	}
	defer adminClient.Close()
	err = adminClient.DeleteInstanceConfig(ctx, &instancepb.DeleteInstanceConfigRequest{
		Name: fmt.Sprintf("projects/%s/instanceConfigs/%s", projectID, userConfigID),
	})
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Deleted instance configuration [%s]\n", userConfigID)
	return nil
}

// [END spanner_delete_instance_config]
