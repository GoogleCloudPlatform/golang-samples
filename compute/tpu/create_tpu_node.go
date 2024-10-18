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

// [START tpu_vm_create]
import (
	"context"
	"fmt"
	"io"

	tpu "cloud.google.com/go/tpu/apiv1"
	"cloud.google.com/go/tpu/apiv1/tpupb"
)

// createTpuNode creates TPU node in given location
func createTPUNode(w io.Writer, projectID, location, nodeName string) error {
	// projectID := "your_project_id"
	// location := "europe-central2-b"
	// nodeName := "your_instance_name"

	ctx := context.Background()
	client, err := tpu.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewTpuClient: %w", err)
	}
	defer client.Close()

	req := &tpupb.CreateNodeRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, location),
		NodeId: nodeName,
		Node: &tpupb.Node{
			// Format: v<TPU version>-<number of chips>
			// TPU is available only in certain locations.
			// Check out the list: https://cloud.google.com/tpu/docs/regions-zones
			AcceleratorType: "v2-8",
			// Specifies an image, being installed on a VM.
			// Apart from OS, image contains also libraries for training models.
			// Read more: https://cloud.google.com/tpu/docs/runtimes
			TensorflowVersion: "tpu-vm-tf-2.14.1",
		},
	}

	op, err := client.CreateNode(ctx, req)
	if err != nil {
		return err
	}

	node, err := op.Wait(ctx)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "Node created: %s", node.GetName())

	return nil
}

// [END tpu_vm_create]