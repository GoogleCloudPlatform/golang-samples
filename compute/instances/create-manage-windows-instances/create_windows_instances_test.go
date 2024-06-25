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

func deleteInstance(ctx context.Context, projectId, zone, instanceName string) error {
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return err
	}
	req := &computepb.DeleteInstanceRequest{
		Project:  projectId,
		Zone:     zone,
		Instance: instanceName,
	}

	op, err := instancesClient.Delete(ctx, req)
	if err != nil {
		return err
	}

	return op.Wait(ctx)
}

func TestComputeCreateWindowsInstancesSnippets(t *testing.T) {
	ctx := context.Background()
	var r *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	tc := testutil.SystemTest(t)
	zone := "europe-central2-b"
	instanceName := fmt.Sprintf("test-vm-%v-%v", time.Now().Format("01-02-2006"), r.Int())
	firewallRuleName := "test-" + fmt.Sprint(r.Int())
	routeName := "test-" + fmt.Sprint(r.Int())
	machineType := "n1-standard-1"
	networkLink := "global/networks/default"
	subnetworkLink := "regions/europe-central2/subnetworks/default"
	sourceImageFamily := "windows-2022"
	want := "Instance created"

	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		t.Fatalf("NewInstancesRESTClient: %v", err)
	}
	defer instancesClient.Close()

	firewallsClient, err := compute.NewFirewallsRESTClient(ctx)
	if err != nil {
		t.Fatalf("NewFirewallsRESTClient: %v", err)
	}
	defer instancesClient.Close()

	routesClient, err := compute.NewRoutesRESTClient(ctx)
	if err != nil {
		t.Fatalf("NewRoutesRESTClient: %v", err)
	}
	defer routesClient.Close()

	buf := &bytes.Buffer{}

	if err := createWndowsServerInstanceExternalIP(buf, tc.ProjectID, zone, instanceName, machineType, sourceImageFamily); err != nil {
		t.Errorf("createWndowsServerInstanceExternalIP got err: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("createWndowsServerInstanceExternalIP got %q, want %q", got, want)
	}

	err = deleteInstance(ctx, tc.ProjectID, zone, instanceName)
	if err != nil {
		t.Errorf("deleteInstance got err: %v", err)
	}

	buf.Reset()

	if err := createWndowsServerInstanceInternalIP(buf, tc.ProjectID, zone, instanceName, machineType, sourceImageFamily, networkLink, subnetworkLink); err != nil {
		t.Errorf("createWndowsServerInstanceInternalIP got err: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("createWndowsServerInstanceInternalIP got %q, want %q", got, want)
	}

	buf.Reset()
	want = "Firewall rule created"

	if err := createFirewallRuleForWindowsActivationHost(buf, tc.ProjectID, firewallRuleName, networkLink); err != nil {
		t.Errorf("createFirewallRuleForWindowsActivationHost got err: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("createFirewallRuleForWindowsActivationHost got %q, want %q", got, want)
	}

	buf.Reset()
	want = "Route created"

	if err := createRouteToWindowsActivationHost(buf, tc.ProjectID, routeName, networkLink); err != nil {
		t.Errorf("createRouteToWindowsActivationHost got err: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("createRouteToWindowsActivationHost got %q, want %q", got, want)
	}

	// Delete route
	deleteRoutereq := &computepb.DeleteRouteRequest{
		Project: tc.ProjectID,
		Route:   routeName,
	}

	op, err := routesClient.Delete(ctx, deleteRoutereq)
	if err != nil {
		t.Errorf("unable to delete route: %v", err)
	}

	if err = op.Wait(ctx); err != nil {
		t.Errorf("unable to wait for the operation: %v", err)
	}

	// Delete firewall rule
	deleteFirewallReq := &computepb.DeleteFirewallRequest{
		Project:  tc.ProjectID,
		Firewall: firewallRuleName,
	}

	op, err = firewallsClient.Delete(ctx, deleteFirewallReq)
	if err != nil {
		t.Errorf("unable to delete firewall rule: %v", err)
	}

	if err = op.Wait(ctx); err != nil {
		t.Errorf("unable to wait for the operation: %v", err)
	}

	err = deleteInstance(ctx, tc.ProjectID, zone, instanceName)
	if err != nil {
		t.Errorf("deleteInstance got err: %v", err)
	}
}
