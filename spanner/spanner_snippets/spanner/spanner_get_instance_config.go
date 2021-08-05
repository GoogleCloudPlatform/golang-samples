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
	instancepb "google.golang.org/genproto/googleapis/spanner/admin/instance/v1"
)

func getInstanceConfig(ctx context.Context, w io.Writer, instanceConfigName string) error {
	adminClient, err := instance.NewInstanceAdminClient(ctx)
	if err != nil {
		return err
	}
	defer adminClient.Close()

	op, err := adminClient.GetInstanceConfig(ctx, &instancepb.GetInstanceConfigRequest{
		Name: instanceConfigName,
	})

	if err != nil {
		return err
	}

	fmt.Fprintf(w, "Available leader options for instance config [%s]: [%s]", op.GetDisplayName(), op.GetLeaderOptions())

	return nil
}
