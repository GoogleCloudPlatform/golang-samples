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

// [START compute_create_egress_rule_windows_activation]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/protobuf/proto"
)

// createFirewallRuleForWindowsActivationHost creates an egress firewall rule with
// the highest priority for host kms.windows.googlecloud.com (35.190.247.13)
// for Windows activation.
func createFirewallRuleForWindowsActivationHost(
	w io.Writer,
	projectID, firewallRuleName, networkName string,
) error {
	// projectID := "your_project_id"
	// firewallRuleName := "your_firewall_rule_name"
	// networkName := "global/networks/default"

	ctx := context.Background()
	firewallsClient, err := compute.NewFirewallsRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewFirewallsRESTClient: %w", err)
	}
	defer firewallsClient.Close()

	req := &computepb.InsertFirewallRequest{
		Project: projectID,
		FirewallResource: &computepb.Firewall{
			Name: proto.String(firewallRuleName),
			Allowed: []*computepb.Allowed{
				{
					IPProtocol: proto.String("tcp"),
					Ports:      []string{"1688"},
				},
			},
			Direction:         proto.String("EGRESS"),
			Network:           proto.String(networkName),
			DestinationRanges: []string{"35.190.247.13/32"},
			Priority:          proto.Int32(0),
		},
	}

	op, err := firewallsClient.Insert(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to create firewall rule: %w", err)
	}

	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %w", err)
	}

	fmt.Fprintf(w, "Firewall rule created\n")

	return nil
}

// [END compute_create_egress_rule_windows_activation]
