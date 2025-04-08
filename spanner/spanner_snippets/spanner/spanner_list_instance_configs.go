// Copyright 2021 Google LLC
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

// [START spanner_list_instance_configs]
import (
	"context"
	"fmt"
	"io"

	instance "cloud.google.com/go/spanner/admin/instance/apiv1"
	"cloud.google.com/go/spanner/admin/instance/apiv1/instancepb"
	"google.golang.org/api/iterator"
)

// istInstanceConfigs gets available leader options for all instances
func listInstanceConfigs(w io.Writer, projectName string) error {
	// projectName = `projects/<project>
	ctx := context.Background()
	instanceAdmin, err := instance.NewInstanceAdminClient(ctx)
	if err != nil {
		return err
	}
	defer instanceAdmin.Close()

	request := &instancepb.ListInstanceConfigsRequest{
		Parent: projectName,
	}
	for {
		iter := instanceAdmin.ListInstanceConfigs(ctx, request)
		for {
			ic, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return err
			}
			fmt.Fprintf(w, "Available leader options for instance config %s: %v\n", ic.Name, ic.LeaderOptions)
		}
		pageToken := iter.PageInfo().Token
		if pageToken == "" {
			break
		} else {
			request.PageToken = pageToken
		}
	}

	return nil
}

// [END spanner_list_instance_configs]
