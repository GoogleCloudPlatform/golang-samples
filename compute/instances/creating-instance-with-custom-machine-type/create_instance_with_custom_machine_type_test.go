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

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	compute "cloud.google.com/go/compute/apiv1"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
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

func getInstance(ctx context.Context, projectID, zone, instanceName string) (*computepb.Instance, error) {
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return nil, err
	}
	reqInstance := &computepb.GetInstanceRequest{
		Project:  projectID,
		Zone:     zone,
		Instance: instanceName,
	}

	return instancesClient.Get(ctx, reqInstance)
}

func TestComputeCreateInstanceWithCustomMachineTypeSnippets(t *testing.T) {
	ctx := context.Background()
	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	tc := testutil.SystemTest(t)
	zone := "europe-central2-b"
	instanceName := "test-" + fmt.Sprint(seededRand.Int())
	expectedResult := "Instance created"

	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		t.Fatalf("NewInstancesRESTClient: %v", err)
	}
	defer instancesClient.Close()

	buf := &bytes.Buffer{}

	customMT := fmt.Sprintf("zones/%s/machineTypes/n2-custom-8-10240", zone)

	if err := createInstanceWithCustomMachineType(buf, tc.ProjectID, zone, instanceName, customMT); err != nil {
		t.Errorf("createInstanceWithCustomMachineType got err: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("createInstanceWithCustomMachineType got %q, want %q", got, expectedResult)
	}

	err = deleteInstance(ctx, tc.ProjectID, zone, instanceName)
	if err != nil {
		t.Errorf("deleteInstance got err: %v", err)
	}

	buf.Reset()

	if err := createInstanceWithCustomMachineTypeWithHelper(buf, tc.ProjectID, zone, instanceName, E2, 4, 8192, CPUSeries_E2_Limit); err != nil {
		t.Errorf("createInstanceWithCustomMachineTypeWithHelper got err: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("createInstanceWithCustomMachineTypeWithHelper got %q, want %q", got, expectedResult)
	}

	instance, err := getInstance(ctx, tc.ProjectID, zone, instanceName)
	if err != nil {
		t.Errorf("unable to get instance: %v", err)
	}

	expectedResult = fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/zones/%s/machineTypes/e2-custom-4-8192", tc.ProjectID, zone)
	if instance.GetMachineType() != expectedResult {
		t.Errorf("incorrect instance MachineType got %q, want %q", instance.GetMachineType(), expectedResult)
	}

	err = deleteInstance(ctx, tc.ProjectID, zone, instanceName)
	if err != nil {
		t.Errorf("deleteInstance got err: %v", err)
	}

	buf.Reset()

	if err := createInstanceWithCustomSharedCore(buf, tc.ProjectID, zone, instanceName, E2_MICRO, 2048, CPUSeries_E2_MICRO_Limit); err != nil {
		t.Errorf("createInstanceWithCustomSharedCore got err: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("createInstanceWithCustomSharedCore got %q, want %q", got, expectedResult)
	}

	instance, err = getInstance(ctx, tc.ProjectID, zone, instanceName)
	if err != nil {
		t.Errorf("unable to get instance: %v", err)
	}

	expectedResult = fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/zones/%s/machineTypes/e2-custom-micro-2048", tc.ProjectID, zone)
	if instance.GetMachineType() != expectedResult {
		t.Errorf("incorrect instance MachineType got %q, want %q", instance.GetMachineType(), expectedResult)
	}

	err = deleteInstance(ctx, tc.ProjectID, zone, instanceName)
	if err != nil {
		t.Errorf("deleteInstance got err: %v", err)
	}

	buf.Reset()

	expectedResult = "Instance updated"

	if err := createInstanceWithCustomMachineType(buf, tc.ProjectID, zone, instanceName, customMT); err != nil {
		t.Errorf("createInstanceWithCustomMachineType got err: %v", err)
	}

	if err := modifyInstanceWithExtendedMemory(buf, tc.ProjectID, zone, instanceName, 819200); err != nil {
		t.Errorf("modifyInstanceWithExtendedMemory got err: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("modifyInstanceWithExtendedMemory got %q, want %q", got, expectedResult)
	}

	instance, err = getInstance(ctx, tc.ProjectID, zone, instanceName)
	if err != nil {
		t.Errorf("unable to get instance: %v", err)
	}

	if !strings.HasSuffix(instance.GetMachineType(), "819200-ext") {
		t.Errorf("incorrect instance MachineType got %q, want suffix %q", instance.GetMachineType(), "819200-ext")
	}

	err = deleteInstance(ctx, tc.ProjectID, zone, instanceName)
	if err != nil {
		t.Errorf("deleteInstance got err: %v", err)
	}
}
