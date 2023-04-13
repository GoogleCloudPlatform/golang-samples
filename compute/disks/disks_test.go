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
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
	"google.golang.org/protobuf/proto"
)

func createInstance(
	ctx context.Context,
	projectID, zone, instanceName, sourceImage, diskName string,
) error {
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return err
	}
	defer instancesClient.Close()
	req := &computepb.InsertInstanceRequest{
		Project: projectID,
		Zone:    zone,
		InstanceResource: &computepb.Instance{
			Name: proto.String(instanceName),
			Disks: []*computepb.AttachedDisk{
				{
					InitializeParams: &computepb.AttachedDiskInitializeParams{
						DiskSizeGb:  proto.Int64(25),
						SourceImage: proto.String(sourceImage),
						DiskName:    proto.String(diskName),
					},
					AutoDelete: proto.Bool(false),
					Boot:       proto.Bool(true),
					DeviceName: proto.String(diskName),
				},
			},
			MachineType: proto.String(fmt.Sprintf("zones/%s/machineTypes/n1-standard-1", zone)),
			NetworkInterfaces: []*computepb.NetworkInterface{
				{
					Name: proto.String("global/networks/default"),
				},
			},
		},
	}

	op, err := instancesClient.Insert(ctx, req)
	if err != nil {
		return err
	}

	return op.Wait(ctx)
}

func getInstance(
	ctx context.Context,
	projectID, zone, instanceName string,
) (*computepb.Instance, error) {
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return nil, err
	}
	defer instancesClient.Close()
	reqInstance := &computepb.GetInstanceRequest{
		Project:  projectID,
		Zone:     zone,
		Instance: instanceName,
	}

	return instancesClient.Get(ctx, reqInstance)
}

func createDiskSnapshot(ctx context.Context, projectId, zone, diskName, snapshotName string) error {
	disksClient, err := compute.NewDisksRESTClient(ctx)
	if err != nil {
		return err
	}
	defer disksClient.Close()
	req := &computepb.CreateSnapshotDiskRequest{
		Project: projectId,
		Zone:    zone,
		Disk:    diskName,
		SnapshotResource: &computepb.Snapshot{
			Name: proto.String(snapshotName),
		},
	}

	op, err := disksClient.CreateSnapshot(ctx, req)
	if err != nil {
		return err
	}

	return op.Wait(ctx)
}

func deleteDiskSnapshot(ctx context.Context, projectId, snapshotName string) error {
	snapshotsClient, err := compute.NewSnapshotsRESTClient(ctx)
	if err != nil {
		return err
	}
	defer snapshotsClient.Close()
	req := &computepb.DeleteSnapshotRequest{
		Project:  projectId,
		Snapshot: snapshotName,
	}

	op, err := snapshotsClient.Delete(ctx, req)
	if err != nil {
		return err
	}

	return op.Wait(ctx)
}

func deleteInstance(ctx context.Context, projectId, zone, instanceName string) error {
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return err
	}
	defer instancesClient.Close()
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

func TestComputeDisksSnippets(t *testing.T) {
	t.Skip("skipping for flakes. see googlecloudplatform/golang-samples#2929")
	ctx := context.Background()
	var r *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	tc := testutil.SystemTest(t)
	zone := "europe-central2-b"
	region := "europe-central2"
	replicaZones := []string{"europe-central2-a", "europe-central2-b"}
	instanceName := fmt.Sprintf("test-instance-%v-%v", time.Now().Format("01-02-2006"), r.Int())
	diskName := fmt.Sprintf("test-disk-%v-%v", time.Now().Format("01-02-2006"), r.Int())
	diskName2 := fmt.Sprintf("test-disk-%v-%v", time.Now().Format("01-02-2006"), r.Int())
	instanceDiskName := fmt.Sprintf("test-instance-disk-%v-%v", time.Now().Format("01-02-2006"), r.Int())
	instanceDiskName2 := fmt.Sprintf("test-instance-disk-%v-%v", time.Now().Format("01-02-2006"), r.Int())
	instanceDiskName3 := fmt.Sprintf("test-instance-disk-%v-%v", time.Now().Format("01-02-2006"), r.Int())
	snapshotName := fmt.Sprintf("test-snapshot-%v-%v", time.Now().Format("01-02-2006"), r.Int())
	sourceImage := "projects/debian-cloud/global/images/family/debian-11"
	sourceDisk := fmt.Sprintf("projects/%s/zones/europe-central2-b/disks/%s", tc.ProjectID, diskName)
	diskType := fmt.Sprintf("zones/%s/diskTypes/pd-ssd", zone)
	diskSnapshotLink := fmt.Sprintf("projects/%s/global/snapshots/%s", tc.ProjectID, snapshotName)

	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		t.Fatalf("NewInstancesRESTClient: %v", err)
	}

	defer instancesClient.Close()

	// Create a snapshot before we run the actual tests
	buf := &bytes.Buffer{}
	err = createDiskFromImage(buf, tc.ProjectID, zone, diskName, diskType, sourceImage, 50)
	if err != nil {
		t.Fatalf("createDiskFromImage got err: %v", err)
	}
	err = createDiskSnapshot(ctx, tc.ProjectID, zone, diskName, snapshotName)
	if err != nil {
		t.Fatalf("createDiskSnapshot got err: %v", err)
	}
	// Create a VM instance to attach disks to
	createInstance(ctx, tc.ProjectID, zone, instanceName, sourceImage, instanceDiskName)
	if err != nil {
		t.Fatalf("unable to create instance: %v", err)
	}

	t.Run("Create zonal disk from a snapshot", func(t *testing.T) {
		buf := &bytes.Buffer{}
		want := "Disk created"

		if err := createDiskFromSnapshot(buf, tc.ProjectID, zone, diskName2, diskType, diskSnapshotLink, 50); err != nil {
			t.Errorf("createDiskFromSnapshot got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, want) {
			t.Errorf("createDiskFromSnapshot got %q, want %q", got, want)
		}
	})

	t.Run("Delete a disk", func(t *testing.T) {
		buf := &bytes.Buffer{}
		want := "Disk deleted"

		if err := deleteDisk(buf, tc.ProjectID, zone, diskName2); err != nil {
			t.Errorf("deleteDisk got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, want) {
			t.Errorf("deleteDisk got %q, want %q", got, want)
		}
	})

	t.Run("Create a regional disk from a snapshot", func(t *testing.T) {
		buf := &bytes.Buffer{}
		want := "Disk created"

		if err := createRegionalDiskFromSnapshot(buf, tc.ProjectID, region, replicaZones, diskName2, diskType, diskSnapshotLink, 50); err != nil {
			t.Errorf("createRegionalDiskFromSnapshot got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, want) {
			t.Errorf("createRegionalDiskFromSnapshot got %q, want %q", got, want)
		}
	})

	t.Run("Delete a zonal disk", func(t *testing.T) {
		buf := &bytes.Buffer{}
		want := "Disk deleted"
		if err := deleteDisk(buf, tc.ProjectID, zone, diskName); err != nil {
			t.Errorf("deleteDisk got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, want) {
			t.Errorf("deleteRegionalDisk got %q, want %q", got, want)
		}
	})

	t.Run("Delete a regional disk", func(t *testing.T) {
		buf := &bytes.Buffer{}
		want := "Disk deleted"

		err = deleteRegionalDisk(buf, tc.ProjectID, region, diskName2)
		if err != nil {
			t.Errorf("deleteRegionalDisk got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, want) {
			t.Errorf("deleteRegionalDisk got %q, want %q", got, want)
		}
	})

	t.Run("Create and resize a regional disk", func(t *testing.T) {
		buf := &bytes.Buffer{}
		want := "Disk created"

		if err := createRegionalDisk(buf, tc.ProjectID, region, replicaZones, diskName2, diskType, 20); err != nil {
			t.Errorf("createRegionalDisk got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, want) {
			t.Errorf("createRegionalDisk got %q, want %q", got, want)
		}

		buf.Reset()
		want = "Disk resized"

		resizeRegionalDisk(buf, tc.ProjectID, region, diskName2, 50)
		if err != nil {
			t.Errorf("resizeRegionalDisk got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, want) {
			t.Errorf("resizeRegionalDisk got %q, want %q", got, want)
		}

		buf.Reset()
		want = "Disk deleted"

		// clean up
		err = deleteRegionalDisk(buf, tc.ProjectID, region, diskName2)
		if err != nil {
			t.Errorf("deleteRegionalDisk got err: %v", err)
		}
	})

	t.Run("createEmptyDisk and clone it into a regional disk", func(t *testing.T) {
		buf := &bytes.Buffer{}
		want := "Disk created"

		if err := createEmptyDisk(buf, tc.ProjectID, zone, diskName, diskType, 20); err != nil {
			t.Fatalf("createEmptyDisk got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, want) {
			t.Errorf("createEmptyDisk got %q, want %q", got, want)
		}

		if err := createRegionalDiskFromDisk(buf, tc.ProjectID, region, replicaZones, diskName2, diskType, sourceDisk, 30); err != nil {
			t.Errorf("createRegionalDiskFromDisk got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, want) {
			t.Errorf("createRegionalDiskFromDisk got %q, want %q", got, want)
		}

		// clean up
		err = deleteRegionalDisk(buf, tc.ProjectID, region, diskName2)
		if err != nil {
			t.Errorf("deleteRegionalDisk got err: %v", err)
		}

		err = deleteDisk(buf, tc.ProjectID, zone, diskName)
		if err != nil {
			t.Errorf("deleteDisk got err: %v", err)
		}
	})

	t.Run("create, clone and delete an encrypted disk", func(t *testing.T) {
		buf.Reset()
		want := "Disk created"

		if err := createEncryptedDisk(buf, tc.ProjectID, zone, diskName, diskType, 20, "SGVsbG8gZnJvbSBHb29nbGUgQ2xvdWQgUGxhdGZvcm0=", "", "", ""); err != nil {
			t.Fatalf("createEncryptedDisk got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, want) {
			t.Errorf("createEncryptedDisk got %q, want %q", got, want)
		}

		if err := createDiskFromCustomerEncryptedDisk(buf, tc.ProjectID, zone, diskName2, diskType, 20, sourceDisk, "SGVsbG8gZnJvbSBHb29nbGUgQ2xvdWQgUGxhdGZvcm0="); err != nil {
			t.Fatalf("createDiskFromCustomerEncryptedDisk got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, want) {
			t.Errorf("createDiskFromCustomerEncryptedDisk got %q, want %q", got, want)
		}

		// cleanup
		err = deleteDisk(buf, tc.ProjectID, zone, diskName)
		if err != nil {
			t.Errorf("deleteDisk got err: %v", err)
		}
	})

	t.Run("setDiskAutoDelete", func(t *testing.T) {
		buf.Reset()
		want := "disk autoDelete field updated."

		if err := setDiskAutoDelete(buf, tc.ProjectID, zone, instanceName, instanceDiskName, true); err != nil {
			t.Fatalf("setDiskAutodelete got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, want) {
			t.Fatalf("setDiskAutoDelete got %q, want %q", got, want)
		}

		instance, err := getInstance(ctx, tc.ProjectID, zone, instanceName)
		if err != nil {
			t.Fatalf("getInstance got err: %v", err)
		}

		if instance.GetDisks()[0].GetAutoDelete() != true {
			t.Errorf("instance got %t, want %t", instance.GetDisks()[0].GetAutoDelete(), true)
		}
	})

	t.Run("Attach a regional disk to VM", func(t *testing.T) {
		if err := createRegionalDisk(buf, tc.ProjectID, region, replicaZones, instanceDiskName2, "regions/us-west3/diskTypes/pd-ssd", 20); err != nil {
			t.Fatalf("createRegionalDisk got err: %v", err)
		}

		buf.Reset()
		want := "Disk attached"

		diskUrl := fmt.Sprintf("projects/%s/regions/%s/disks/%s", tc.ProjectID, region, instanceDiskName2)

		if err := attachRegionalDisk(buf, tc.ProjectID, zone, instanceName, diskUrl); err != nil {
			t.Fatalf("attachRegionalDisk got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, want) {
			t.Fatalf("attachRegionalDisk got %q, want %q", got, want)
		}

		instance, err := getInstance(ctx, tc.ProjectID, zone, instanceName)
		if err != nil {
			t.Fatalf("getInstance got err: %v", err)
		}

		foundDisk := false
		for _, disk := range instance.GetDisks() {
			if strings.Contains(*disk.Source, instanceDiskName2) {
				foundDisk = true
			}
		}
		if !foundDisk {
			t.Errorf("The disk %s is not attached to the instance!", instanceDiskName2)
		}

		// cannot clean up the disk just yet because it must be done after the VM is terminated
	})

	t.Run("Attach a read-only regional disk to VM", func(t *testing.T) {
		if err := createRegionalDisk(buf, tc.ProjectID, region, replicaZones, instanceDiskName3, "regions/us-west3/diskTypes/pd-ssd", 20); err != nil {
			t.Fatalf("createRegionalDisk got err: %v", err)
		}

		buf.Reset()
		want := "Disk attached"

		diskUrl := fmt.Sprintf("projects/%s/regions/%s/disks/%s", tc.ProjectID, region, instanceDiskName3)

		if err := attachRegionalDiskReadOnly(buf, tc.ProjectID, zone, instanceName, diskUrl); err != nil {
			t.Fatalf("attachRegionalDisk got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, want) {
			t.Fatalf("attachRegionalDisk got %q, want %q", got, want)
		}

		instance, err := getInstance(ctx, tc.ProjectID, zone, instanceName)
		if err != nil {
			t.Fatalf("getInstance got err: %v", err)
		}

		foundDisk := false
		for _, disk := range instance.GetDisks() {
			if strings.Contains(*disk.Source, instanceDiskName3) {
				foundDisk = true
			}
		}
		if !foundDisk {
			t.Errorf("The disk %s is not attached to the instance!", instanceDiskName3)
		}

		// cannot clean up the disk just yet because it must be done after the VM is terminated
	})

	// clean up
	err = deleteDiskSnapshot(ctx, tc.ProjectID, snapshotName)
	if err != nil {
		t.Errorf("deleteDiskSnapshot got err: %v", err)
	}

	err = deleteInstance(ctx, tc.ProjectID, zone, instanceName)
	if err != nil {
		t.Errorf("deleteInstance got err: %v", err)
	}

	// disks attached to a VM instance must be deleted after deleting the instance

	// disk 1 is set to auto-delete, no need to delete it explicitly

	if err := deleteRegionalDisk(buf, tc.ProjectID, region, instanceDiskName2); err != nil {
		t.Errorf("deleteRegionalDisk got err: %v", err)
	}

	if err := deleteRegionalDisk(buf, tc.ProjectID, region, instanceDiskName3); err != nil {
		t.Fatalf("deleteRegionalDisk got err: %v", err)
	}
}
