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
	"google.golang.org/protobuf/proto"
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

func TestSuspendResumeSnippets(t *testing.T) {
	ctx := context.Background()
	var r *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	tc := testutil.SystemTest(t)
	zone := "europe-central2-b"
	instanceName := fmt.Sprintf("test-vm-%v-%v", time.Now().Format("01-02-2006"), r.Int())

	buf := &bytes.Buffer{}

	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		t.Fatalf("NewInstancesRESTClient: %v", err)
	}
	defer instancesClient.Close()

	req := &computepb.InsertInstanceRequest{
		Project: tc.ProjectID,
		Zone:    zone,
		InstanceResource: &computepb.Instance{
			Name: proto.String(instanceName),
			Disks: []*computepb.AttachedDisk{
				{
					InitializeParams: &computepb.AttachedDiskInitializeParams{
						DiskSizeGb: proto.Int64(64),
						SourceImage: proto.String(
							"projects/ubuntu-os-cloud/global/images/family/ubuntu-2004-lts",
						),
					},
					AutoDelete: proto.Bool(true),
					Boot:       proto.Bool(true),
				},
			},
			MachineType: proto.String(fmt.Sprintf("zones/%s/machineTypes/n1-standard-1", zone)),
			NetworkInterfaces: []*computepb.NetworkInterface{
				{
					Name: proto.String("default"),
				},
			},
		},
	}

	op, err := instancesClient.Insert(ctx, req)
	if err != nil {
		t.Fatalf("unable to create instance: %v", err)
	}

	if err = op.Wait(ctx); err != nil {
		t.Errorf("unable to wait for the operation: %v", err)
	}

	// Once the machine is running, give it some time to fully start all processes
	// before trying to suspend it
	time.Sleep(45 * time.Second)

	err = suspendInstance(buf, tc.ProjectID, zone, instanceName)
	if err != nil {
		t.Errorf("suspendInstance got err: %v", err)
	}

	want := "Instance suspended"
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("suspendInstance got %q, want %q", got, want)
	}

	getInstanceReq := &computepb.GetInstanceRequest{
		Project:  tc.ProjectID,
		Zone:     zone,
		Instance: instanceName,
	}

	instance, err := instancesClient.Get(ctx, getInstanceReq)
	if err != nil {
		t.Errorf("unable to get instance: %v", err)
	}

	for instance.GetStatus() == "SUSPENDING" {
		instance, err = instancesClient.Get(ctx, getInstanceReq)
		if err != nil {
			t.Errorf("unable to get instance: %v", err)
		}
		time.Sleep(5 * time.Second)
	}

	instance, err = instancesClient.Get(ctx, getInstanceReq)
	if err != nil {
		t.Errorf("unable to get instance: %v", err)
	}

	if instance.GetStatus() != "SUSPENDED" {
		t.Errorf("incorrect instance status got %q, want SUSPENDED", instance.GetStatus())
	}

	buf.Reset()

	err = resumeInstance(buf, tc.ProjectID, zone, instanceName)
	if err != nil {
		t.Errorf("resumeInstance got err: %v", err)
	}

	want = "Instance resumed"
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("resumeInstance got %q, want %q", got, want)
	}

	instance, err = instancesClient.Get(ctx, getInstanceReq)
	if err != nil {
		t.Errorf("unable to get instance: %v", err)
	}

	if instance.GetStatus() != "RUNNING" {
		t.Errorf("incorrect instance status got %q, want RUNNING", instance.GetStatus())
	}

	err = deleteInstance(ctx, tc.ProjectID, zone, instanceName)
	if err != nil {
		t.Errorf("deleteInstance got err: %v", err)
	}

}
