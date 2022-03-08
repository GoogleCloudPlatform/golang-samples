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
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
)

func TestComputeCreateInstanceSnippets(t *testing.T) {
	ctx := context.Background()
	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	tc := testutil.SystemTest(t)
	zone := "europe-central2-b"
	instanceName := "test-" + fmt.Sprint(seededRand.Int())
	networkName := "global/networks/default"
	subnetworkName := "regions/europe-central2/subnetworks/default"
	expectedResult := "Instance created"

	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		t.Fatalf("NewInstancesRESTClient: %v", err)
	}
	defer instancesClient.Close()

	imagesClient, err := compute.NewImagesRESTClient(ctx)
	if err != nil {
		t.Fatalf("NewImagesRESTClient: %v", err)
	}
	defer imagesClient.Close()

	newestDebianReq := &computepb.GetFromFamilyImageRequest{
		Project: "debian-cloud",
		Family:  "debian-11",
	}

	newestDebian, err := imagesClient.GetFromFamily(ctx, newestDebianReq)
	if err != nil {
		t.Errorf("unable to get image from family: %v", err)
	}

	buf := &bytes.Buffer{}

	if err := createInstanceFromCustomImage(buf, tc.ProjectID, zone, instanceName, *newestDebian.SelfLink); err != nil {
		t.Errorf("createInstanceFromCustomImage got err: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("createInstanceFromCustomImage got %q, want %q", got, expectedResult)
	}

	err = deleteInstance(ctx, tc.ProjectID, zone, instanceName)
	if err != nil {
		t.Errorf("deleteInstance got err: %v", err)
	}

	buf.Reset()

	if err := createWithAdditionalDisk(buf, tc.ProjectID, zone, instanceName); err != nil {
		t.Errorf("createWithAdditionalDisk got err: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("createWithAdditionalDisk got %q, want %q", got, expectedResult)
	}
	err = deleteInstance(ctx, tc.ProjectID, zone, instanceName)
	if err != nil {
		t.Errorf("deleteInstance got err: %v", err)
	}

	buf.Reset()

	if err := createInstanceWithSubnet(buf, tc.ProjectID, zone, instanceName, networkName, subnetworkName); err != nil {
		t.Errorf("createInstanceWithSubnet got err: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("createInstanceWithSubnet got %q, want %q", got, expectedResult)
	}
	err = deleteInstance(ctx, tc.ProjectID, zone, instanceName)
	if err != nil {
		t.Errorf("deleteInstance got err: %v", err)
	}
}
