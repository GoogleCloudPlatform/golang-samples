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

// [START compute_firewall_list]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/api/iterator"
)

// listFirewallRules prints the list of firewall names and their descriptions in specified project
func listFirewallRules(w io.Writer, projectID string) error {
	// projectID := "your_project_id"

	ctx := context.Background()
	firewallsClient, err := compute.NewFirewallsRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewInstancesRESTClient: %w", err)
	}
	defer firewallsClient.Close()

	req := &computepb.ListFirewallsRequest{
		Project: projectID,
	}

	it := firewallsClient.List(ctx, req)
	for {
		firewallRule, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "- %s: %s\n", firewallRule.GetName(), firewallRule.GetDescription())
	}

	return nil
}

// [END compute_firewall_list]
