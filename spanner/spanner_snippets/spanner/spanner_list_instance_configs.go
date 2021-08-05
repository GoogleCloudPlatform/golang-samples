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

// [START spanner_create_database]
import (
	"context"
	"fmt"
	"io"

	instance "cloud.google.com/go/spanner/admin/instance/apiv1"
	"google.golang.org/api/iterator"
	instancepb "google.golang.org/genproto/googleapis/spanner/admin/instance/v1"
)

func listInstanceConfigs(ctx context.Context, w io.Writer, projectName string) error {
	instanceAdmin, err := instance.NewInstanceAdminClient(ctx)
	if err != nil {
		return err
	}
	defer instanceAdmin.Close()

	request := &instancepb.ListInstanceConfigsRequest{
		Parent:   projectName,
		PageSize: 10,
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
			fmt.Printf("Available leader options for instance config %s: %v\n", ic.Name, ic.LeaderOptions)
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
