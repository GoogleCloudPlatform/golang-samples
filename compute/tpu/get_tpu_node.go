// Copyright 2024 Google LLC
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

package snippets

// [START tpu_vm_get]
import (
	"context"
	"fmt"
	"io"

	tpu "cloud.google.com/go/tpu/apiv1"
	"cloud.google.com/go/tpu/apiv1/tpupb"
)

// getTpuNode gets TPU node by given name and location within project
func getTPUNode(w io.Writer, projectID, location, nodeName string) error {
	// projectID := "your_project_id"
	// location := "europe-central2-b"
	// nodeName := "your_instance_name"

	ctx := context.Background()
	client, err := tpu.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewTpuClient: %w", err)
	}
	defer client.Close()

	req := &tpupb.GetNodeRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/nodes/%s", projectID, location, nodeName),
	}

	node, err := client.GetNode(ctx, req)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Got node: %s", node.GetName())

	return nil
}

// [END tpu_vm_get]
