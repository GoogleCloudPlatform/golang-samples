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

// [START spanner_get_instance_config]
import (
	"context"
	"fmt"
	"io"

	instance "cloud.google.com/go/spanner/admin/instance/apiv1"
	"cloud.google.com/go/spanner/admin/instance/apiv1/instancepb"
)

// getInstanceConfig gets available leader options
func getInstanceConfig(w io.Writer, instanceConfigName string) error {
	// defaultLeader = `nam3`
	ctx := context.Background()
	instanceAdmin, err := instance.NewInstanceAdminClient(ctx)
	if err != nil {
		return err
	}
	defer instanceAdmin.Close()

	ic, err := instanceAdmin.GetInstanceConfig(ctx, &instancepb.GetInstanceConfigRequest{
		Name: instanceConfigName,
	})

	if err != nil {
		return fmt.Errorf("could not get instance config %s: %w", instanceConfigName, err)
	}

	fmt.Fprintf(w, "Available leader options for instance config %s: %v", instanceConfigName, ic.LeaderOptions)

	return nil
}

// [END spanner_get_instance_config]
