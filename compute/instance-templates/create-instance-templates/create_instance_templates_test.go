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
	"google.golang.org/protobuf/proto"
)

func TestCreateInstanceTemplatesSnippets(t *testing.T) {
	t.Skip("https://github.com/GoogleCloudPlatform/golang-samples/issues/2383")
	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	tc := testutil.SystemTest(t)
	zone := "europe-central2-b"
	instanceName := "test-instance-" + fmt.Sprint(seededRand.Int())
	templateName1 := "test-template-" + fmt.Sprint(seededRand.Int())
	templateName2 := "test-template-" + fmt.Sprint(seededRand.Int())
	templateName3 := "test-template-" + fmt.Sprint(seededRand.Int())
	machineType := "n1-standard-1"
	sourceImage := "projects/debian-cloud/global/images/family/debian-12"
	networkName := "global/networks/default-compute"
	subnetworkName := "regions/asia-east1/subnetworks/default-compute"

	ctx := context.Background()

	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		t.Fatalf("NewInstancesRESTClient: %v", err)
	}
	defer instancesClient.Close()

	zoneOperationsClient, err := compute.NewZoneOperationsRESTClient(ctx)
	if err != nil {
		t.Fatalf("NewZoneOperationsRESTClient: %v", err)
	}
	defer zoneOperationsClient.Close()

	buf := &bytes.Buffer{}

	if err := createTemplate(buf, tc.ProjectID, templateName1); err != nil {
		t.Errorf("createTemplate got err: %v", err)
	}

	expectedResult := "Instance template created"
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("createTemplate got %q, want %q", got, expectedResult)
	}

	buf.Reset()

	if err := listInstanceTemplates(buf, tc.ProjectID); err != nil {
		t.Errorf("listInstanceTemplates got err: %v", err)
	}

	expectedResult = fmt.Sprintf("- %s %s", templateName1, "e2-standard-4")
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("listInstanceTemplates got %q, want %q", got, expectedResult)
	}

	buf.Reset()

	req := &computepb.InsertInstanceRequest{
		Project: tc.ProjectID,
		Zone:    zone,
		InstanceResource: &computepb.Instance{
			Name: proto.String(instanceName),
			Disks: []*computepb.AttachedDisk{
				{
					InitializeParams: &computepb.AttachedDiskInitializeParams{
						DiskSizeGb:  proto.Int64(250),
						SourceImage: proto.String(sourceImage),
					},
					AutoDelete: proto.Bool(true),
					Boot:       proto.Bool(true),
					Type:       proto.String(computepb.AttachedDisk_PERSISTENT.String()),
					DeviceName: proto.String("disk-1"),
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
		t.Errorf("unable to create instance: %v", err)
	}

	if err = op.Wait(ctx); err != nil {
		t.Errorf("unable to wait for the operation: %v", err)
	}

	formattedInstanceName := fmt.Sprintf("projects/%s/zones/%s/instances/%s", tc.ProjectID, zone, instanceName)
	if err := createTemplateFromInstance(buf, tc.ProjectID, formattedInstanceName, templateName2); err != nil {
		t.Errorf("createTemplateFromInstance got err: %v", err)
	}

	expectedResult = "Instance template created"
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("createTemplateFromInstance got %q, want %q", got, expectedResult)
	}

	template, err := getInstanceTemplate(tc.ProjectID, templateName2)
	if err != nil {
		t.Errorf("getInstanceTemplate got err: %v", err)
	}

	got := template.GetName()
	if got != templateName2 {
		t.Errorf("template.GetName() got %q, want %q", got, templateName2)
	}

	got = template.GetProperties().GetMachineType()
	if got != machineType {
		t.Errorf("template.GetProperties().GetMachineType() got %q, want %q", got, machineType)
	}

	gotDiskSize := template.GetProperties().GetDisks()[0].GetDiskSizeGb()
	if gotDiskSize != 250 {
		t.Errorf("template.GetProperties().GetDisks()[0].GetDiskSizeGb() got %q, want %q", got, 250)
	}

	buf.Reset()

	if err := createTemplateWithSubnet(buf, tc.ProjectID, networkName, subnetworkName, templateName3); err != nil {
		t.Errorf("createTemplateWithSubnet got err: %v", err)
	}

	expectedResult = "Instance template created"
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("createTemplateFromInstance got %q, want %q", got, expectedResult)
	}

	buf.Reset()

	if err := deleteInstanceTemplate(buf, tc.ProjectID, templateName1); err != nil {
		t.Errorf("deleteInstanceTemplate got err: %v", err)
	}

	expectedResult = "Instance template deleted"
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("createTemplateFromInstance got %q, want %q", got, expectedResult)
	}

	buf.Reset()

	if err := deleteInstanceTemplate(buf, tc.ProjectID, templateName2); err != nil {
		t.Errorf("deleteInstanceTemplate got err: %v", err)
	}

	expectedResult = "Instance template deleted"
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("createTemplateFromInstance got %q, want %q", got, expectedResult)
	}

	buf.Reset()

	if err := deleteInstanceTemplate(buf, tc.ProjectID, templateName3); err != nil {
		t.Errorf("deleteInstanceTemplate got err: %v", err)
	}

	expectedResult = "Instance template deleted"
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("createTemplateFromInstance got %q, want %q", got, expectedResult)
	}

	buf.Reset()

	deleteReq := &computepb.DeleteInstanceRequest{
		Project:  tc.ProjectID,
		Zone:     zone,
		Instance: instanceName,
	}

	op, err = instancesClient.Delete(ctx, deleteReq)
	if err != nil {
		t.Errorf("unable to delete instance: %v", err)
	}

	if err = op.Wait(ctx); err != nil {
		t.Errorf("unable to wait for the operation: %v", err)
	}

	t.Run("regional template", func(t *testing.T) {
		buf.Reset()
		templateName := fmt.Sprintf("test-template-%d", seededRand.Int())
		region := "eu-central1"
		err := createRegionalTemplate(buf, tc.ProjectID, templateName, region)
		if err != nil {
			t.Errorf("createRegionalTemplate failed: %v", err)
		}

		expectedResult := "Instance template created"
		if got := buf.String(); !strings.Contains(got, expectedResult) {
			t.Errorf("createRegionalTemplate got %q, want %q", got, expectedResult)
		}

		template, err := getRegionalTemplate(tc.ProjectID, templateName, region)
		if err != nil {
			t.Errorf("getRegionalTemplate got err: %v", err)
		}

		got := template.GetName()
		if got != templateName {
			t.Errorf("template.GetName() got %q, want %q", got, templateName)
		}
		buf.Reset()

		err = deleteRegionalTemplate(buf, tc.ProjectID, templateName, region)
		if err != nil {
			t.Errorf("deleteRegionalTemplate failed: %v", err)
		}

		expectedResult = "Instance template deleted"
		if got := buf.String(); !strings.Contains(got, expectedResult) {
			t.Errorf("deleteRegionalTemplate got %q, want %q", got, expectedResult)
		}
	})
}
