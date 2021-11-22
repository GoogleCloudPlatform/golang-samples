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
	"google.golang.org/protobuf/proto"
)

func TestCreateInstanceFromTemplateSnippets(t *testing.T) {
	ctx := context.Background()
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		t.Fatalf("NewInstancesRESTClient: %v", err)
	}
	defer instancesClient.Close()

	instanceTemplatesClient, err := compute.NewInstanceTemplatesRESTClient(ctx)
	if err != nil {
		t.Fatalf("NewInstanceTemplatesRESTClient: %v", err)
	}
	defer instancesClient.Close()

	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	tc := testutil.SystemTest(t)
	zone := "europe-central2-b"
	instanceName := "test-instance-" + fmt.Sprint(seededRand.Int())
	instanceName2 := "test-instance-" + fmt.Sprint(seededRand.Int())
	instanceTemplateName := "test-instance-template" + fmt.Sprint(seededRand.Int())
	machineType := "n1-standard-1"
	sourceImage := "projects/debian-cloud/global/images/family/debian-10"
	networkName := "global/networks/default"

	insertTemplateReq := &computepb.InsertInstanceTemplateRequest{
		Project: tc.ProjectID,
		InstanceTemplateResource: &computepb.InstanceTemplate{
			Name: &instanceTemplateName,
			Properties: &computepb.InstanceProperties{
				MachineType: proto.String(machineType),
				Disks: []*computepb.AttachedDisk{
					{
						InitializeParams: &computepb.AttachedDiskInitializeParams{
							DiskSizeGb:  proto.Int64(10),
							SourceImage: proto.String(sourceImage),
						},
						AutoDelete: proto.Bool(true),
						Boot:       proto.Bool(true),
						Type:       computepb.AttachedDisk_PERSISTENT.Enum(),
					},
				},
				NetworkInterfaces: []*computepb.NetworkInterface{
					{
						Name: proto.String(networkName),
					},
				},
			},
		},
	}

	op, err := instanceTemplatesClient.Insert(ctx, insertTemplateReq)
	if err != nil {
		t.Fatalf("unable to create instance template: %v", err)
	}

	globalOperationsClient, err := compute.NewGlobalOperationsRESTClient(ctx)
	if err != nil {
		t.Errorf("NewGlobalOperationsRESTClient: %v", err)
	}
	defer globalOperationsClient.Close()

	for {
		waitReq := &computepb.WaitGlobalOperationRequest{
			Operation: op.Proto().GetName(),
			Project:   tc.ProjectID,
		}
		globalOp, err := globalOperationsClient.Wait(ctx, waitReq)
		if err != nil {
			t.Errorf("unable to wait for the operation: %v", err)
		}

		if *globalOp.Status.Enum() == computepb.Operation_DONE {
			break
		}
	}

	buf := &bytes.Buffer{}

	if err := createInstanceFromTemplate(buf, tc.ProjectID, zone, instanceName, fmt.Sprintf("global/instanceTemplates/%s", instanceTemplateName)); err != nil {
		t.Errorf("createInstanceFromTemplate: %v", err)
	}

	expectedResult := "Instance created"
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("createInstanceFromTemplate got %q, want %q", got, expectedResult)
	}

	if err := deleteInstance(buf, tc.ProjectID, zone, instanceName); err != nil {
		t.Errorf("deleteInstance got err: %v", err)
	}

	buf.Reset()

	if err := createInstanceFromTemplateWithOverrides(buf, tc.ProjectID, zone, instanceName2, instanceTemplateName, machineType, sourceImage); err != nil {
		t.Errorf("createInstanceFromTemplateWithOverrides: %v", err)
	}

	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("createInstanceFromTemplateWithOverrides got %q, want %q", got, expectedResult)
	}

	instanceReq := &computepb.GetInstanceRequest{
		Project:  tc.ProjectID,
		Zone:     zone,
		Instance: instanceName2,
	}

	instance, err := instancesClient.Get(ctx, instanceReq)
	if err != nil {
		t.Errorf("unable to get instance: %v", err)
	}

	if len(instance.GetDisks()) != 2 {
		t.Errorf("Instance must contain 2 disks")
	}

	if err := deleteInstance(buf, tc.ProjectID, zone, instanceName2); err != nil {
		t.Errorf("deleteInstance got err: %v", err)
	}

	deleteTemplateReq := &computepb.DeleteInstanceTemplateRequest{
		Project:          tc.ProjectID,
		InstanceTemplate: instanceTemplateName,
	}

	op, err = instanceTemplatesClient.Delete(ctx, deleteTemplateReq)
	if err != nil {
		t.Errorf("unable to delete instance template: %v", err)
	}

	for {
		waitReq := &computepb.WaitGlobalOperationRequest{
			Operation: op.Proto().GetName(),
			Project:   tc.ProjectID,
		}
		globalOp, err := globalOperationsClient.Wait(ctx, waitReq)
		if err != nil {
			t.Errorf("unable to wait for the operation: %v", err)
		}

		if *globalOp.Status.Enum() == computepb.Operation_DONE {
			break
		}
	}

}
