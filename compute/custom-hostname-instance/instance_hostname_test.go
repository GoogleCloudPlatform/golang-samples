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
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestInstanceHostnameSnippets(t *testing.T) {
	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	tc := testutil.SystemTest(t)
	ctx := context.Background()
	zone := "europe-central2-b"
	customHostname := "host.domain.com"
	instanceName := "test-" + fmt.Sprint(seededRand.Int())
	machineType := "n1-standard-1"
	sourceImage := "projects/debian-cloud/global/images/family/debian-12"
	networkName := "global/networks/default"
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		t.Fatalf("NewInstancesRESTClient: %v", err)
	}
	defer instancesClient.Close()

	buf := &bytes.Buffer{}

	if err := createInstanceWithCustomHostname(buf, tc.ProjectID, zone, instanceName, customHostname, machineType, sourceImage, networkName); err != nil {
		t.Errorf("createInstanceWithCustomHostname got err: %v", err)
	}

	expectedResult := "Instance created"
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("createInstanceWithCustomHostname got %q, want %q", got, expectedResult)
	}

	instanceReq := &computepb.GetInstanceRequest{
		Project:  tc.ProjectID,
		Zone:     zone,
		Instance: instanceName,
	}

	instance, err := instancesClient.Get(ctx, instanceReq)
	if err != nil {
		t.Errorf("unable to get instance: %v", err)
	}

	if instance.GetHostname() != customHostname {
		t.Errorf("Instance has incorrect hostname got %q, want %q", instance.GetHostname(), customHostname)
	}

	buf.Reset()

	if err := getInstanceHostname(buf, tc.ProjectID, zone, instanceName); err != nil {
		t.Errorf("getInstanceHostname got err: %v", err)
	}

	expectedResult = fmt.Sprintf("Instance %v has hostname: %v", instanceName, instance.GetHostname())
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("getInstanceHostname got %q, want %q", got, expectedResult)
	}

	deleteReq := &computepb.DeleteInstanceRequest{
		Project:  tc.ProjectID,
		Zone:     zone,
		Instance: instanceName,
	}

	op, err := instancesClient.Delete(ctx, deleteReq)
	if err != nil {
		t.Errorf("unable to delete instance: %v", err)
	}

	if err = op.Wait(ctx); err != nil {
		t.Errorf("unable to wait for the operation: %v", err)
	}
}
