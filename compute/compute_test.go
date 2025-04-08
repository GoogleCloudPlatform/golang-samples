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

func TestComputeSnippets(t *testing.T) {
	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	tc := testutil.SystemTest(t)
	zone := "us-central1-a"
	instanceName := "test-" + fmt.Sprint(seededRand.Int())
	instanceName2 := "test-" + fmt.Sprint(seededRand.Int())
	machineType := "n1-standard-1"
	sourceImage := "projects/debian-cloud/global/images/family/debian-12"
	networkName := "global/networks/default"

	buf := &bytes.Buffer{}

	if err := createInstance(buf, tc.ProjectID, zone, instanceName, machineType, sourceImage, networkName); err != nil {
		t.Errorf("createInstance got err: %v", err)
	}

	expectedResult := "Instance created"
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("createInstance got %q, want %q", got, expectedResult)
	}

	buf.Reset()

	if err := getInstance(buf, tc.ProjectID, zone, instanceName); err != nil {
		t.Errorf("getInstance got err: %v", err)
	}

	expectedResult = fmt.Sprintf("Instance: %s", instanceName)
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("getInstance got %q, want %q", got, expectedResult)
	}

	buf.Reset()

	if err := listInstances(buf, tc.ProjectID, zone); err != nil {
		t.Errorf("listInstances got err: %v", err)
	}

	expectedResult = "Instances found in zone"
	expectedResult2 := fmt.Sprintf("- %s", instanceName)
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("listInstances got %q, want %q", got, expectedResult)
	}
	if got := buf.String(); !strings.Contains(got, expectedResult2) {
		t.Errorf("listInstances got %q, want %q", got, expectedResult2)
	}

	buf.Reset()

	if err := listAllInstances(buf, tc.ProjectID); err != nil {
		t.Errorf("listAllInstances got err: %v", err)
	}

	expectedResult = "Instances found:"
	expectedResult2 = fmt.Sprintf("zones/%s\n", zone)
	expectedResult3 := fmt.Sprintf("- %s", instanceName)
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("listAllInstances got %q, want %q", got, expectedResult)
	}
	if got := buf.String(); !strings.Contains(got, expectedResult2) {
		t.Errorf("listAllInstances got %q, want %q", got, expectedResult2)
	}
	if got := buf.String(); !strings.Contains(got, expectedResult2) {
		t.Errorf("listAllInstances got %q, want %q", got, expectedResult3)
	}

	buf.Reset()

	if err := deleteInstance(buf, tc.ProjectID, zone, instanceName); err != nil {
		t.Errorf("deleteInstance got err: %v", err)
	}

	expectedResult = "Instance deleted"
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("deleteInstance got %q, want %q", got, expectedResult)
	}

	if err := createInstance(buf, tc.ProjectID, zone, instanceName2, machineType, sourceImage, networkName); err != nil {
		t.Errorf("createInstance got err: %v", err)
	}

	buf.Reset()

	if err := changeMachineType(buf, tc.ProjectID, zone, instanceName2, "e2-standard-2"); err != nil {
		t.Errorf("changeMachineType got err: %v", err)
	}

	expectedResult = "Instance updated"
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("waitForOperation got %q, want %q", got, expectedResult)
	}

	ctx := context.Background()
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		t.Errorf("NewInstancesRESTClient: %v", err)
	}
	defer instancesClient.Close()

	req := &computepb.DeleteInstanceRequest{
		Project:  tc.ProjectID,
		Zone:     zone,
		Instance: instanceName2,
	}

	op, err := instancesClient.Delete(ctx, req)
	if err != nil {
		t.Errorf("Delete instance request: %v", err)
	}

	buf.Reset()

	if err := waitForOperation(buf, tc.ProjectID, op); err != nil {
		t.Errorf("waitForOperation got err: %v", err)
	}

	expectedResult = "Operation finished"
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("waitForOperation got %q, want %q", got, expectedResult)
	}
}
