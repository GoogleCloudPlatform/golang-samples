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

// [START tpu_vm_list]
import (
	"context"
	"fmt"
	"io"

	tpu "cloud.google.com/go/tpu/apiv1"
	"cloud.google.com/go/tpu/apiv1/tpupb"
	"google.golang.org/api/iterator"
)

// listTPUNodes gets list of TPU nodes by given location within project
func listTPUNodes(w io.Writer, projectID, location string) error {
	// projectID := "your_project_id"
	// location := "europe-central2-b"

	ctx := context.Background()
	client, err := tpu.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewTpuClient: %w", err)
	}
	defer client.Close()

	req := &tpupb.ListNodesRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, location),
	}

	it := client.ListNodes(ctx, req)
	for {
		node, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "- %s", node.GetName())
	}

	return nil
}

// [END tpu_vm_list]
