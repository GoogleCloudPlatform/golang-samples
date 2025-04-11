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

func createDisk(ctx context.Context, projectId, zone, diskName, sourceImage string) error {
	disksClient, err := compute.NewDisksRESTClient(ctx)
	if err != nil {
		return err
	}
	defer disksClient.Close()
	req := &computepb.InsertDiskRequest{
		Project: projectId,
		Zone:    zone,
		DiskResource: &computepb.Disk{
			Name:        proto.String(diskName),
			SourceImage: proto.String(sourceImage),
		},
	}

	op, err := disksClient.Insert(ctx, req)
	if err != nil {
		return err
	}

	return op.Wait(ctx)
}

func deleteDisk(ctx context.Context, projectId, zone, diskName string) error {
	disksClient, err := compute.NewDisksRESTClient(ctx)
	if err != nil {
		return err
	}
	defer disksClient.Close()
	req := &computepb.DeleteDiskRequest{
		Project: projectId,
		Zone:    zone,
		Disk:    diskName,
	}

	op, err := disksClient.Delete(ctx, req)
	if err != nil {
		return err
	}

	return op.Wait(ctx)
}

func TestComputeCreateInstanceSnippets(t *testing.T) {
	ctx := context.Background()
	var r *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	tc := testutil.SystemTest(t)
	zone := "europe-central2-b"
	instanceName := fmt.Sprintf("test-instance-%v-%v", time.Now().Format("01-02-2006"), r.Int())
	bootDiskName := fmt.Sprintf("test-disk-%v-%v", time.Now().Format("01-02-2006"), r.Int())
	diskName2 := fmt.Sprintf("test-disk-%v-%v", time.Now().Format("01-02-2006"), r.Int())
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

	if err := createWithLocalSSD(buf, tc.ProjectID, zone, instanceName); err != nil {
		t.Errorf("createWithLocalSSD got err: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("createWithLocalSSD got %q, want %q", got, expectedResult)
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

	buf.Reset()

	err = createDisk(ctx, tc.ProjectID, zone, diskName2, *newestDebian.SelfLink)
	if err != nil {
		t.Fatalf("createDisk got err: %v", err)
	}

	err = createDisk(ctx, tc.ProjectID, zone, bootDiskName, *newestDebian.SelfLink)
	if err != nil {
		t.Fatalf("createDisk got err: %v", err)
	}

	diskNames := []string{
		bootDiskName,
		diskName2,
	}

	if err := createWithExistingDisks(buf, tc.ProjectID, zone, instanceName, diskNames); err != nil {
		t.Fatalf("createInstanceWithSubnet got err: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("createInstanceWithSubnet got %q, want %q", got, expectedResult)
	}

	err = deleteInstance(ctx, tc.ProjectID, zone, instanceName)
	if err != nil {
		t.Errorf("deleteInstance got err: %v", err)
	}

	err = deleteDisk(ctx, tc.ProjectID, zone, bootDiskName)
	if err != nil {
		t.Errorf("deleteDisk got err: %v", err)
	}

	err = deleteDisk(ctx, tc.ProjectID, zone, diskName2)
	if err != nil {
		t.Errorf("deleteDisk got err: %v", err)
	}
}

func TestComputeBulkCreateInstanceSnippets(t *testing.T) {
	ctx := context.Background()

	instanceTemplatesClient, err := compute.NewInstanceTemplatesRESTClient(ctx)
	if err != nil {
		t.Fatalf("NewInstanceTemplatesRESTClient: %v", err)
	}
	defer instanceTemplatesClient.Close()

	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	tc := testutil.SystemTest(t)
	zone := "europe-central2-b"
	instanceTemplateName := "test-instance-template" + fmt.Sprint(seededRand.Int())
	machineType := "n1-standard-1"
	sourceImage := "projects/debian-cloud/global/images/family/debian-12"
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
						Type:       proto.String(computepb.AttachedDisk_PERSISTENT.String()),
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

	if err = op.Wait(ctx); err != nil {
		t.Errorf("unable to wait for the operation: %v", err)
	}

	namePattern := "i-#-" + fmt.Sprint(seededRand.Int())[0:5]
	buf := &bytes.Buffer{}

	instances, err := createFiveInstances(buf, tc.ProjectID, zone, instanceTemplateName, namePattern)
	if err != nil {
		t.Errorf("createFiveInstances got err: %v", err)
	}

	for _, instance := range instances {
		if err := deleteInstance(ctx, tc.ProjectID, zone, *instance.Name); err != nil {
			t.Errorf("deleteInstance got err: %v", err)
		}
	}

	deleteTemplateReq := &computepb.DeleteInstanceTemplateRequest{
		Project:          tc.ProjectID,
		InstanceTemplate: instanceTemplateName,
	}

	op, err = instanceTemplatesClient.Delete(ctx, deleteTemplateReq)
	if err != nil {
		t.Errorf("unable to delete instance template: %v", err)
	}
}

func TestCreateWithReplica(t *testing.T) {
	ctx := context.Background()

	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	tc := testutil.SystemTest(t)
	region := "europe-central2"
	zone := "europe-central2-b"
	replicaZones := []string{"europe-central2-a", "europe-central2-b"}
	instanceName := "test-instance" + fmt.Sprint(seededRand.Int())
	diskName := "test-disk" + fmt.Sprint(seededRand.Int())
	replicatedDiskName := "test-replicated-disk" + fmt.Sprint(seededRand.Int())
	snapshotName := "test-snapshot" + fmt.Sprint(seededRand.Int())
	snapshotLink := fmt.Sprintf("projects/%s/global/snapshots/%s", tc.ProjectID, snapshotName)
	var buf bytes.Buffer

	imagesClient, err := compute.NewImagesRESTClient(ctx)
	if err != nil {
		t.Fatalf("NewImagesRESTClient: %v", err)
	}
	defer imagesClient.Close()

	snapshotsClient, err := compute.NewSnapshotsRESTClient(ctx)
	if err != nil {
		t.Fatalf("NewSnapshotsRESTClient: %v", err)
	}
	defer snapshotsClient.Close()

	disksClient, err := compute.NewDisksRESTClient(ctx)
	if err != nil {
		t.Fatalf("NewDisksRESTClient: %v", err)
	}
	defer disksClient.Close()

	newestDebianReq := &computepb.GetFromFamilyImageRequest{
		Project: "debian-cloud",
		Family:  "debian-11",
	}

	newestDebian, err := imagesClient.GetFromFamily(ctx, newestDebianReq)
	if err != nil {
		t.Errorf("unable to get image from family: %v", err)
	}

	insertDiskReq := &computepb.InsertDiskRequest{
		Project: tc.ProjectID,
		Zone:    zone,
		DiskResource: &computepb.Disk{
			SourceImage: newestDebian.SelfLink,
			Name:        proto.String(diskName),
		},
	}

	op, err := disksClient.Insert(ctx, insertDiskReq)
	if err != nil {
		t.Errorf("unable to create disk: %v", err)
	}

	if err = op.Wait(ctx); err != nil {
		t.Errorf("unable to wait for the operation: %v", err)
	}
	defer deleteDisk(ctx, tc.ProjectID, zone, diskName)

	diskSnaphotReq := &computepb.CreateSnapshotDiskRequest{
		Project: tc.ProjectID,
		Zone:    zone,
		Disk:    diskName,
		SnapshotResource: &computepb.Snapshot{
			Name: proto.String(snapshotName),
		},
	}

	op, err = disksClient.CreateSnapshot(ctx, diskSnaphotReq)
	if err != nil {
		t.Errorf("unable to create disk snapshot: %v", err)
	}

	if err = op.Wait(ctx); err != nil {
		t.Errorf("unable to wait for the operation: %v", err)
	}

	if err := createReplicatedBootDisk(&buf, tc.ProjectID, zone, snapshotLink, replicatedDiskName, instanceName, replicaZones); err != nil {
		t.Errorf("createReplicatedBootDisk failed: %v", err)
	}

	// Cleanup
	err = deleteInstance(ctx, tc.ProjectID, zone, instanceName)
	if err != nil {
		t.Errorf("deleteInstance got err: %v", err)
	}

	regionDisksClient, err := compute.NewRegionDisksRESTClient(ctx)
	if err != nil {
		t.Errorf("NewDisksRESTClient: %v", err)
	}
	defer disksClient.Close()

	req := &computepb.DeleteRegionDiskRequest{
		Project: tc.ProjectID,
		Region:  region,
		Disk:    replicatedDiskName,
	}

	op, err = regionDisksClient.Delete(ctx, req)
	if err != nil {
		t.Errorf("unable to delete disk: %v", err)
	}

	if err = op.Wait(ctx); err != nil {
		t.Errorf("unable to wait for the operation: %v", err)
	}

	delReq := &computepb.DeleteSnapshotRequest{
		Project:  tc.ProjectID,
		Snapshot: snapshotName,
	}

	op, err = snapshotsClient.Delete(ctx, delReq)
	if err != nil {
		t.Errorf("unable to delete snapshot: %v", err)
	}

	if err = op.Wait(ctx); err != nil {
		t.Errorf("unable to wait for the operation: %v", err)
	}
}
