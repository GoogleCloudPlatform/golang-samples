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

// [START tpu_vm_stop]
import (
	"context"
	"fmt"
	"io"

	tpu "cloud.google.com/go/tpu/apiv1"
	"cloud.google.com/go/tpu/apiv1/tpupb"
)

// stopTPUNode stops TPU node
func stopTPUNode(w io.Writer, nodeName string) error {
	// nodeName := "your_instance_name"

	ctx := context.Background()
	client, err := tpu.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewTpuClient: %w", err)
	}
	defer client.Close()

	req := &tpupb.StopNodeRequest{
		Name: nodeName,
	}

	op, err := client.StopNode(ctx, req)
	if err != nil {
		return err
	}

	node, err := op.Wait(ctx)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Node %s is stopped", node.GetName())

	return nil
}

// [END tpu_vm_stop]
