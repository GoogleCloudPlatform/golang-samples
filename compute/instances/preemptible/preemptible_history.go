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

package snippets

// [START compute_preemptible_history]
import (
	"context"
	"fmt"
	"io"
	"strings"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/api/iterator"
	"google.golang.org/protobuf/proto"
)

// preemptionHisory gets a list of preemption operations from given zone in a project.
// Optionally limit the results to instance name.
func preemptionHisory(w io.Writer, projectID, zone, instanceName, customFilter string) error {
	// projectID := "your_project_id"
	// zone := "europe-central2-b"
	// instanceName := "your_instance_name"
	// customFilter := "operationType=\"compute.instances.preempted\""

	ctx := context.Background()
	operationsClient, err := compute.NewZoneOperationsRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewZoneOperationsRESTClient: %w", err)
	}
	defer operationsClient.Close()

	filter := ""

	if customFilter != "" {
		filter = customFilter
	} else {
		filter = "operationType=\"compute.instances.preempted\""

		if instanceName != "" {
			filter += fmt.Sprintf(
				` AND targetLink="https://www.googleapis.com/compute/v1/projects/%s/zones/%s/instances/%s"`,
				projectID, zone, instanceName,
			)
		}
	}

	req := &computepb.ListZoneOperationsRequest{
		Project: projectID,
		Zone:    zone,
		Filter:  proto.String(filter),
	}

	it := operationsClient.List(ctx, req)
	for {
		operation, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}

		ss := strings.Split(operation.GetTargetLink(), "/")
		curInstName := ss[len(ss)-1]
		if curInstName == instanceName {
			// The filter used is not 100% accurate, it's `contains` not `equals`
			// So we need to check the name to make sure it's the one we want.
			fmt.Fprintf(w, "- %s %s\n", instanceName, operation.GetInsertTime())
		}
	}

	return nil
}

// [END compute_preemptible_history]
