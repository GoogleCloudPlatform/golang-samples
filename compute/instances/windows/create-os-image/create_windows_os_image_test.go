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

func TestCreateWindowsOSImageSnippets(t *testing.T) {
	ctx := context.Background()
	var r *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	tc := testutil.SystemTest(t)
	zone := "europe-central2-b"
	instanceName := fmt.Sprintf("test-vm-%v-%v", time.Now().Format("01-02-2006"), r.Int())
	diskName := "test-" + fmt.Sprint(r.Int())
	imageName := "test-" + fmt.Sprint(r.Int())
	storageLocation := "eu"

	buf := &bytes.Buffer{}

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

	req := &computepb.InsertInstanceRequest{
		Project: tc.ProjectID,
		Zone:    zone,
		InstanceResource: &computepb.Instance{
			Name: proto.String(instanceName),
			Disks: []*computepb.AttachedDisk{
				{
					// Describe the size and source image of the boot disk to attach to the instance.
					InitializeParams: &computepb.AttachedDiskInitializeParams{
						DiskName:   proto.String(diskName),
						DiskSizeGb: proto.Int64(64),
						SourceImage: proto.String(
							"projects/windows-cloud/global/images/windows-server-2022-dc-core-v20231011",
						),
					},
					DeviceName: proto.String(diskName),
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
		t.Errorf("unable to create instance: %v", err)
	}

	if err = op.Wait(ctx); err != nil {
		t.Errorf("unable to wait for the operation: %v", err)
	}

	want := fmt.Sprintf("instance %s should be stopped.", instanceName)
	err = createWindowsOSImage(buf, tc.ProjectID, zone, diskName, imageName, storageLocation, false)
	if err != nil {
		if got := fmt.Sprint(err); !strings.Contains(got, want) {
			t.Errorf("createWindowsOSImage error got %q, want %q", got, want)
		}
	} else {
		t.Errorf("createWindowsOSImage should return an error")
	}

	buf.Reset()

	err = createWindowsOSImage(buf, tc.ProjectID, zone, diskName, imageName, storageLocation, true)
	if err != nil {
		t.Errorf("createWindowsOSImage got err: %v", err)
	}

	want = "Image created"
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("createWindowsOSImage got %q, want %q", got, want)
	}

	deleteImageReq := &computepb.DeleteImageRequest{
		Project: tc.ProjectID,
		Image:   imageName,
	}

	op, err = imagesClient.Delete(ctx, deleteImageReq)
	if err != nil {
		t.Errorf("unable to delete image: %v", err)
	}

	if err = op.Wait(ctx); err != nil {
		t.Errorf("unable to wait for the operation: %v", err)
	}

	err = deleteInstance(ctx, tc.ProjectID, zone, instanceName)
	if err != nil {
		t.Errorf("deleteInstance got err: %v", err)
	}
}
