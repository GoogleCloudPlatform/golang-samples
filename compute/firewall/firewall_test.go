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

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestFirewallSnippets(t *testing.T) {
	ctx := context.Background()
	firewallsClient, err := compute.NewFirewallsRESTClient(ctx)
	if err != nil {
		t.Fatalf("NewInstancesRESTClient: %v", err)
	}
	defer firewallsClient.Close()

	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	tc := testutil.SystemTest(t)
	buf := &bytes.Buffer{}

	firewallRuleName := "test-firewall-rule" + fmt.Sprint(seededRand.Int())
	defaultNetwork := "global/networks/default"

	if err := createFirewallRule(buf, tc.ProjectID, firewallRuleName, defaultNetwork); err != nil {
		t.Fatalf("createFirewallRule got err: %v", err)
	}

	expectedResult := "Firewall rule created"
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("createFirewallRule got %q, want %q", got, expectedResult)
	}

	buf.Reset()

	req := &computepb.GetFirewallRequest{
		Project:  tc.ProjectID,
		Firewall: firewallRuleName,
	}

	firewallBefore, err := firewallsClient.Get(ctx, req)
	if err != nil {
		t.Errorf("unable to get firewall rule: %v", err)
	}

	if firewallBefore.GetPriority() != 1000 {
		t.Errorf(fmt.Sprintf("Got: %q; want %q", firewallBefore.GetPriority(), 1000))
	}

	var newFirewallPriority int32 = 500
	if err := patchFirewallPriority(buf, tc.ProjectID, firewallRuleName, newFirewallPriority); err != nil {
		t.Fatalf("patchFirewallPriority got err: %v", err)
	}

	expectedResult = "Firewall rule updated"
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("patchFirewallPriority got %q, want %q", got, expectedResult)
	}

	firewallAfter, err := firewallsClient.Get(ctx, req)
	if err != nil {
		t.Errorf("unable to get firewall rule: %v", err)
	}

	if firewallAfter.GetPriority() != newFirewallPriority {
		t.Errorf(fmt.Sprintf("Got: %q; want %q", firewallAfter.GetPriority(), newFirewallPriority))
	}

	buf.Reset()

	if err := listFirewallRules(buf, tc.ProjectID); err != nil {
		t.Fatalf("listFirewallRules got err: %v", err)
	}

	expectedResult = fmt.Sprintf("- %s:", firewallRuleName)
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("listFirewallRules got %q, want %q", got, expectedResult)
	}

	buf.Reset()

	if err := deleteFirewallRule(buf, tc.ProjectID, firewallRuleName); err != nil {
		t.Errorf("deleteFirewallRule got err: %v", err)
	}

	expectedResult = "Firewall rule deleted"
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("deleteFirewallRule got %q, want %q", got, expectedResult)
	}
}
