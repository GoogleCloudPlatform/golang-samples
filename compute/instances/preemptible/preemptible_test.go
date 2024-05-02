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

func TestPreemptibleSnippets(t *testing.T) {
	ctx := context.Background()
	var r *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	tc := testutil.SystemTest(t)
	zone := "europe-central2-b"
	instanceName := fmt.Sprintf("test-vm-%v-%v", time.Now().Format("01-02-2006"), r.Int())

	buf := &bytes.Buffer{}

	if err := createPreemtibleInstance(buf, tc.ProjectID, zone, instanceName); err != nil {
		t.Fatalf("createPreemtibleInstance got err: %v", err)
	}

	want := "Instance created"
	if got := buf.String(); !strings.Contains(got, want) {
		t.Fatalf("createPreemtibleInstance got %q, want %q", got, want)
	}

	buf.Reset()

	if err := printPreemtible(buf, tc.ProjectID, zone, instanceName); err != nil {
		t.Errorf("printPreemtible got err: %v", err)
	}

	want = "Is instance preemptible: true"
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("printPreemtible got %q, want %q", got, want)
	}

	buf.Reset()

	customFilter := fmt.Sprintf(`targetLink="https://www.googleapis.com/compute/v1/projects/%s/zones/%s/instances/%s"`, tc.ProjectID, zone, instanceName)

	if err := preemptionHisory(buf, tc.ProjectID, zone, instanceName, customFilter); err != nil {
		t.Errorf("preemptionHisory got err: %v", err)
	}

	want = fmt.Sprintf("- %s", instanceName)
	if got := buf.String(); !strings.Contains(got, want) {
		t.Fatalf("preemptionHisory got %q, want %q", got, want)
	}

	if err := deleteInstance(ctx, tc.ProjectID, zone, instanceName); err != nil {
		t.Errorf("deleteInstance got err: %v", err)
	}
}
