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

func TestPreventingAccidentalVMDeletionSnippets(t *testing.T) {
	ctx := context.Background()
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		t.Fatalf("NewInstancesRESTClient: %v", err)
	}
	defer instancesClient.Close()

	var r *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	tc := testutil.SystemTest(t)
	zone := "europe-central2-b"
	instanceName := fmt.Sprintf("test-vm-%v-%v", time.Now().Format("01-02-2006"), r.Int())

	buf := &bytes.Buffer{}

	if err := createInstance(buf, tc.ProjectID, zone, instanceName, true); err != nil {
		t.Fatalf("createInstance got err: %v", err)
	}

	want := "Instance created"
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("createInstance got %q, want %q", got, want)
	}

	buf.Reset()

	if err := getDeleteProtection(buf, tc.ProjectID, zone, instanceName); err != nil {
		t.Errorf("getDeleteProtection got err: %v", err)
	}

	want = fmt.Sprintf("Instance %s has DeleteProtection value: %v", instanceName, true)
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("getDeleteProtection got %q, want %q", got, want)
	}

	buf.Reset()

	if err := setDeleteProtection(buf, tc.ProjectID, zone, instanceName, false); err != nil {
		t.Errorf("setDeleteProtection got err: %v", err)
	}

	want = "Instance updated"
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("setDeleteProtection got %q, want %q", got, want)
	}

	buf.Reset()

	if err := getDeleteProtection(buf, tc.ProjectID, zone, instanceName); err != nil {
		t.Errorf("getDeleteProtection got err: %v", err)
	}

	want = fmt.Sprintf("Instance %s has DeleteProtection value: %v", instanceName, false)
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("getDeleteProtection got %q, want %q", got, want)
	}

	buf.Reset()

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
