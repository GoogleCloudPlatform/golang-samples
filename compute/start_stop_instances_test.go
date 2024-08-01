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
	b64 "encoding/base64"
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

func TestStartStopSnippets(t *testing.T) {
	ctx := context.Background()
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		t.Fatalf("NewInstancesRESTClient: %v", err)
	}
	defer instancesClient.Close()

	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	tc := testutil.SystemTest(t)
	zone := "europe-central2-b"
	instanceName := "test-instance-" + fmt.Sprint(seededRand.Int())
	instanceName2 := "test-instance-" + fmt.Sprint(seededRand.Int())
	machineType := "n1-standard-1"
	sourceImage := "projects/debian-cloud/global/images/family/debian-12"
	networkName := "global/networks/default"

	buf := &bytes.Buffer{}

	if err := createInstance(buf, tc.ProjectID, zone, instanceName, machineType, sourceImage, networkName); err != nil {
		t.Fatalf("createInstance got err: %v", err)
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

	if *instance.Status != computepb.Instance_RUNNING.String() {
		t.Errorf("Instance is not in running status")
	}

	buf.Reset()

	if err := stopInstance(buf, tc.ProjectID, zone, instanceName); err != nil {
		t.Errorf("stopInstance got err: %v", err)
	}

	expectedResult := "Instance stopped"
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("stopInstance got %q, want %q", got, expectedResult)
	}

	instance, err = instancesClient.Get(ctx, instanceReq)
	if err != nil {
		t.Errorf("unable to get instance: %v", err)
	}

	if *instance.Status != computepb.Instance_TERMINATED.String() {
		t.Errorf("Instance is not in terminated status")
	}

	buf.Reset()

	if err := startInstance(buf, tc.ProjectID, zone, instanceName); err != nil {
		t.Errorf("startInstance got err: %v", err)
	}

	expectedResult = "Instance started"
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("startInstance got %q, want %q", got, expectedResult)
	}

	instance, err = instancesClient.Get(ctx, instanceReq)
	if err != nil {
		t.Errorf("unable to get instance: %v", err)
	}

	if *instance.Status != computepb.Instance_RUNNING.String() {
		t.Errorf("Instance is not in running status")
	}

	buf.Reset()

	if err := resetInstance(buf, tc.ProjectID, zone, instanceName); err != nil {
		t.Errorf("resetInstance got err: %v", err)
	}

	expectedResult = "Instance reset"
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("resetInstance got %q, want %q", got, expectedResult)
	}

	if err := deleteInstance(buf, tc.ProjectID, zone, instanceName); err != nil {
		t.Errorf("deleteInstance got err: %v", err)
	}

	base64Key := b64.RawStdEncoding.EncodeToString([]byte("random-random-random-random-s123"))

	instanceReq = &computepb.GetInstanceRequest{
		Project:  tc.ProjectID,
		Zone:     zone,
		Instance: instanceName2,
	}

	req := &computepb.InsertInstanceRequest{
		Project: tc.ProjectID,
		Zone:    zone,
		InstanceResource: &computepb.Instance{
			Name: proto.String(instanceName2),
			Disks: []*computepb.AttachedDisk{
				{
					InitializeParams: &computepb.AttachedDiskInitializeParams{
						DiskSizeGb:  proto.Int64(10),
						SourceImage: proto.String(sourceImage),
					},
					AutoDelete: proto.Bool(true),
					Boot:       proto.Bool(true),
					Type:       proto.String(computepb.AttachedDisk_PERSISTENT.String()),
					DiskEncryptionKey: &computepb.CustomerEncryptionKey{
						RawKey: proto.String(base64Key),
					},
				},
			},
			MachineType: proto.String(fmt.Sprintf("zones/%s/machineTypes/%s", zone, machineType)),
			NetworkInterfaces: []*computepb.NetworkInterface{
				{
					Name: proto.String(networkName),
				},
			},
		},
	}

	op, err := instancesClient.Insert(ctx, req)
	if err != nil {
		t.Fatalf("unable to create instance: %v", err)
	}

	err = op.Wait(ctx)
	if err != nil {
		t.Errorf("unable to wait for the operation: %v", err)
	}

	buf.Reset()

	if err := stopInstance(buf, tc.ProjectID, zone, instanceName2); err != nil {
		t.Errorf("stopInstance got err: %v", err)
	}

	instance, err = instancesClient.Get(ctx, instanceReq)
	if err != nil {
		t.Errorf("unable to get instance: %v", err)
	}

	if *instance.Status != computepb.Instance_TERMINATED.String() {
		t.Errorf("Instance is not in terminated status")
	}

	if err := startInstanceWithEncKey(buf, tc.ProjectID, zone, instanceName2, base64Key); err != nil {
		t.Errorf("startInstanceWithEncKey got err: %v", err)
	}

	expectedResult = "Instance with encryption key started"
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("rstartInstanceWithEncKey got %q, want %q", got, expectedResult)
	}

	instance, err = instancesClient.Get(ctx, instanceReq)
	if err != nil {
		t.Errorf("unable to get instance: %v", err)
	}

	if *instance.Status != computepb.Instance_RUNNING.String() {
		t.Errorf("Instance is not in running status")
	}

	if err := deleteInstance(buf, tc.ProjectID, zone, instanceName2); err != nil {
		t.Errorf("deleteInstance got err: %v", err)
	}

}
