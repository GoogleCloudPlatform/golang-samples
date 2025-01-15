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

// [START compute_firewall_create]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/protobuf/proto"
)

// createFirewallRule creates a firewall rule allowing for incoming HTTP and HTTPS access from the entire Internet.
func createFirewallRule(w io.Writer, projectID, firewallRuleName, networkName string) error {
	// projectID := "your_project_id"
	// firewallRuleName := "europe-central2-b"
	// networkName := "global/networks/default"

	ctx := context.Background()
	firewallsClient, err := compute.NewFirewallsRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewInstancesRESTClient: %w", err)
	}
	defer firewallsClient.Close()

	firewallRule := &computepb.Firewall{
		Allowed: []*computepb.Allowed{
			{
				IPProtocol: proto.String("tcp"),
				Ports:      []string{"80", "443"},
			},
		},
		Direction: proto.String(computepb.Firewall_INGRESS.String()),
		Name:      &firewallRuleName,
		TargetTags: []string{
			"web",
		},
		Network:     &networkName,
		Description: proto.String("Allowing TCP traffic on port 80 and 443 from Internet."),
	}

	// Note that the default value of priority for the firewall API is 1000.
	// If you check the value of `firewallRule.GetPriority()` at this point it
	// will be equal to 0, however it is not treated as "set" by the library and thus
	// the default will be applied to the new rule. If you want to create a rule that
	// has priority == 0, you need to explicitly set it so:

	// firewallRule.Priority = proto.Int32(0)

	req := &computepb.InsertFirewallRequest{
		Project:          projectID,
		FirewallResource: firewallRule,
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

// [END compute_firewall_create]
