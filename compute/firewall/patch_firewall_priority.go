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

// [START compute_firewall_patch]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/protobuf/proto"
)

// patchFirewallPriority modifies the priority of a given firewall rule.
func patchFirewallPriority(w io.Writer, projectID, firewallRuleName string, priority int32) error {
	// projectID := "your_project_id"
	// firewallRuleName := "europe-central2-b"
	// priority := 10

	ctx := context.Background()
	firewallsClient, err := compute.NewFirewallsRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewInstancesRESTClient: %w", err)
	}
	defer firewallsClient.Close()

	firewallRule := &computepb.Firewall{
		Priority: proto.Int32(priority),
	}

	req := &computepb.PatchFirewallRequest{
		Project:          projectID,
		Firewall:         firewallRuleName,
		FirewallResource: firewallRule,
	}

	// The patch operation doesn't require the full definition of a Firewall interface. It will only update
	// the values that were set in it, in this case it will only change the priority.
	op, err := firewallsClient.Patch(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to patch firewall rule: %w", err)
	}

	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %w", err)
	}

	fmt.Fprintf(w, "Firewall rule updated\n")

	return nil
}

// [END compute_firewall_patch]
