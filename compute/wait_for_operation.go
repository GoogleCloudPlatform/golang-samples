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

package snippets

// [START compute_instances_operation_check]
import (
	"context"
	"fmt"
	"io"
	"strings"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
)

// waitForOperation waits for an operation to be completed. Calling this function will block until the operation is finished.
func waitForOperation(w io.Writer, op *computepb.Operation, projectID string) error {
	// projectID := "your_project_id"
	ctx := context.Background()
	zoneArr := strings.Split(op.GetZone(), "/")

	zoneOperationsClient, err := compute.NewZoneOperationsRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewZoneOperationsRESTClient: %v", err)
	}
	defer zoneOperationsClient.Close()

	for {
		waitReq := &computepb.WaitZoneOperationRequest{
			Operation: op.GetName(),
			Project:   projectID,
			Zone:      zoneArr[len(zoneArr)-1],
		}

		// Waits for the specified Operation resource to return as DONE or for the request to approach the 2 minute deadline.
		op, err = zoneOperationsClient.Wait(ctx, waitReq)
		if err != nil {
			return fmt.Errorf("unable to wait for the operation: %v", err)
		}

		if op.GetStatus() == computepb.Operation_DONE {
			fmt.Fprintf(w, "Operation finished\n")
			break
		}
	}

	return nil
}

// [END compute_instances_operation_check]
