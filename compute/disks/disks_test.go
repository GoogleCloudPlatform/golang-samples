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
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"google.golang.org/api/googleapi"
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

func deleteInstance(t *testing.T, ctx context.Context, projectId, zone, instanceName string) {
	t.Helper()
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		t.Fatalf("NewInstancesRESTClient: %v", err)
	}
	defer instancesClient.Close()

	// Get the instance to set all disks to autodelte with intance
	instance, err := getInstance(ctx, projectId, zone, instanceName)
	if err != nil {
		t.Error("getInstance err", err)
	}

	for _, disk := range instance.GetDisks() {
		req := &computepb.SetDiskAutoDeleteInstanceRequest{
			Project:    projectId,
			Zone:       zone,
			Instance:   instanceName,
			DeviceName: disk.GetDeviceName(),
			AutoDelete: true,
		}

		_, err := instancesClient.SetDiskAutoDelete(ctx, req)
		if err != nil {
			t.Errorf("unable to set disk autodelete field: %v", err)
		}
	}

	req := &computepb.DeleteInstanceRequest{
		Project:  projectId,
		Zone:     zone,
		Instance: instanceName,
	}

	op, err := instancesClient.Delete(ctx, req)
	if err != nil {
		t.Errorf("instanceClient.Delete: %v", err)
	}

	err = op.Wait(ctx)
	if err != nil {
		t.Errorf("instanceClient.Delete: %v", err)
	}
}

// Produces a test error only in case it was NOT due to a 404. This avoids
// flackyness which may result from the ripper cleaning up resources.
func errorIfNot404(t *testing.T, msg string, err error) {
	var gerr *googleapi.Error
	if errors.As(err, &gerr) {
		t.Log(gerr.Message, " - ", gerr.Code)
		if gerr.Code == 404 {
			t.Skip(msg + " skipped due to a Not Found error (404)")
		} else {
			t.Errorf(msg+" got err: %v", err)
		}
	}
}

func TestComputeDisksSnippets(t *testing.T) {
	ctx := context.Background()
	var r *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	tc := testutil.SystemTest(t)
	zone := "europe-west2-b"
	region := "europe-west2"
	replicaZones := []string{"europe-west2-a", "europe-west2-b"}
	instanceName := fmt.Sprintf("test-instance-%v-%v", time.Now().Format("01-02-2006"), r.Int())
	diskName := fmt.Sprintf("test-disk-%v-%v", time.Now().Format("01-02-2006"), r.Int())
	instanceDiskName := fmt.Sprintf("test-instance-disk-%v-%v", time.Now().Format("01-02-2006"), r.Int())
	snapshotName := fmt.Sprintf("test-snapshot-%v-%v", time.Now().Format("01-02-2006"), r.Int())
	sourceImage := "projects/debian-cloud/global/images/family/debian-11"
	diskType := fmt.Sprintf("zones/%s/diskTypes/pd-ssd", zone)
	diskSnapshotLink := fmt.Sprintf("projects/%s/global/snapshots/%s", tc.ProjectID, snapshotName)

	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		t.Fatalf("NewInstancesRESTClient: %v", err)
	}
	defer instancesClient.Close()

	// Create a snapshot before we run the actual tests
	var buf bytes.Buffer
	err = createDiskFromImage(&buf, tc.ProjectID, zone, diskName, diskType, sourceImage, 50)
	if err != nil {
		t.Fatalf("createDiskFromImage got err: %v", err)
	}
	defer deleteDisk(&buf, tc.ProjectID, zone, diskName)

	err = createDiskSnapshot(ctx, tc.ProjectID, zone, diskName, snapshotName)
	if err != nil {
		t.Fatalf("createDiskSnapshot got err: %v", err)
	}
	defer deleteDiskSnapshot(ctx, tc.ProjectID, snapshotName)

	// Create a VM instance to attach disks to
	err = createInstance(ctx, tc.ProjectID, zone, instanceName, sourceImage, instanceDiskName)
	if err != nil {
		t.Fatalf("unable to create instance: %v", err)
	}
	defer deleteInstance(t, ctx, tc.ProjectID, zone, instanceName)

	t.Run("Create and delete zonal disk from a snapshot", func(t *testing.T) {
		zonalDiskName := fmt.Sprintf("test-zonal-disk-%v-%v", time.Now().Format("01-02-2006"), r.Int())
		var buf bytes.Buffer
		want := "Disk created"

		if err := createDiskFromSnapshot(&buf, tc.ProjectID, zone, zonalDiskName, diskType, diskSnapshotLink, 50); err != nil {
			t.Errorf("createDiskFromSnapshot got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, want) {
			t.Errorf("createDiskFromSnapshot got %q, want %q", got, want)
		}

		buf.Reset()
		want = "Disk deleted"

		if err := deleteDisk(&buf, tc.ProjectID, zone, zonalDiskName); err != nil {
			errorIfNot404(t, "deleteDisk", err)
		}
		if got := buf.String(); !strings.Contains(got, want) {
			t.Errorf("deleteDisk got %q, want %q", got, want)
		}
	})

	t.Run("Create and delete a regional disk from a snapshot", func(t *testing.T) {
		regionalDiskName := fmt.Sprintf("test-regional-disk-%v-%v", time.Now().Format("01-02-2006"), r.Int())
		var buf bytes.Buffer
		want := "Disk created"

		if err := createRegionalDiskFromSnapshot(&buf, tc.ProjectID, region, replicaZones, regionalDiskName, diskType, diskSnapshotLink, 50); err != nil {
			t.Errorf("createRegionalDiskFromSnapshot got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, want) {
			t.Errorf("createRegionalDiskFromSnapshot got %q, want %q", got, want)
		}

		buf.Reset()
		want = "Disk deleted"

		err = deleteRegionalDisk(&buf, tc.ProjectID, region, regionalDiskName)
		if err != nil {
			errorIfNot404(t, "deleteRegionalDisk", err)
		}
		if got := buf.String(); !strings.Contains(got, want) {
			t.Errorf("deleteRegionalDisk got %q, want %q", got, want)
		}
	})

	t.Run("Create and resize a regional disk", func(t *testing.T) {
		regionalDiskName := fmt.Sprintf("test-regional-resize-disk-%v-%v", time.Now().Format("01-02-2006"), r.Int())
		var buf bytes.Buffer
		want := "Disk created"

		if err := createRegionalDisk(&buf, tc.ProjectID, region, replicaZones, regionalDiskName, diskType, 20); err != nil {
			t.Errorf("createRegionalDisk got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, want) {
			t.Errorf("createRegionalDisk got %q, want %q", got, want)
		}

		buf.Reset()
		want = "Disk resized"

		resizeRegionalDisk(&buf, tc.ProjectID, region, regionalDiskName, 50)
		if err != nil {
			t.Errorf("resizeRegionalDisk got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, want) {
			t.Errorf("resizeRegionalDisk got %q, want %q", got, want)
		}

		buf.Reset()
		want = "Disk deleted"

		// clean up
		err = deleteRegionalDisk(&buf, tc.ProjectID, region, regionalDiskName)
		if err != nil {
			errorIfNot404(t, "deleteRegionalDisk", err)
		}
	})

	t.Run("createEmptyDisk and clone it into a regional disk", func(t *testing.T) {
		zonalDiskName := fmt.Sprintf("test-zonal-clone-disk-%v-%v", time.Now().Format("01-02-2006"), r.Int())
		sourceDisk := fmt.Sprintf("projects/%s/zones/europe-west2-b/disks/%s", tc.ProjectID, zonalDiskName)
		regionalDiskName := fmt.Sprintf("test-regional-clone-disk-%v-%v", time.Now().Format("01-02-2006"), r.Int())
		var buf bytes.Buffer
		want := "Disk created"

		if err := createEmptyDisk(&buf, tc.ProjectID, zone, zonalDiskName, diskType, 20); err != nil {
			t.Fatalf("createEmptyDisk got err: %v", err)
		}
		defer deleteDisk(&buf, tc.ProjectID, zone, zonalDiskName)

		if got := buf.String(); !strings.Contains(got, want) {
			t.Errorf("createEmptyDisk got %q, want %q", got, want)
		}

		if err := createRegionalDiskFromDisk(&buf, tc.ProjectID, region, replicaZones, regionalDiskName, diskType, sourceDisk, 30); err != nil {
			t.Fatalf("createRegionalDiskFromDisk got err: %v", err)
		}
		defer deleteRegionalDisk(&buf, tc.ProjectID, region, regionalDiskName)

		if got := buf.String(); !strings.Contains(got, want) {
			t.Errorf("createRegionalDiskFromDisk got %q, want %q", got, want)
		}
	})

	t.Run("create, clone and delete an encrypted disk", func(t *testing.T) {
		encDiskName1 := fmt.Sprintf("test-enc-disk1-%v-%v", time.Now().Format("01-02-2006"), r.Int())
		sourceDisk := fmt.Sprintf("projects/%s/zones/europe-west2-b/disks/%s", tc.ProjectID, encDiskName1)
		encDiskName2 := fmt.Sprintf("test-enc-disk2-%v-%v", time.Now().Format("01-02-2006"), r.Int())
		var buf bytes.Buffer
		want := "Disk created"

		if err := createEncryptedDisk(&buf, tc.ProjectID, zone, encDiskName1, diskType, 20, "SGVsbG8gZnJvbSBHb29nbGUgQ2xvdWQgUGxhdGZvcm0=", "", "", ""); err != nil {
			t.Fatalf("createEncryptedDisk got err: %v", err)
		}
		defer deleteDisk(&buf, tc.ProjectID, zone, encDiskName1)

		if got := buf.String(); !strings.Contains(got, want) {
			t.Errorf("createEncryptedDisk got %q, want %q", got, want)
		}

		if err := createDiskFromCustomerEncryptedDisk(&buf, tc.ProjectID, zone, encDiskName2, diskType, 20, sourceDisk, "SGVsbG8gZnJvbSBHb29nbGUgQ2xvdWQgUGxhdGZvcm0="); err != nil {
			t.Fatalf("createDiskFromCustomerEncryptedDisk got err: %v", err)
		}
		defer deleteDisk(&buf, tc.ProjectID, zone, encDiskName2)

		if got := buf.String(); !strings.Contains(got, want) {
			t.Errorf("createDiskFromCustomerEncryptedDisk got %q, want %q", got, want)
		}
	})

	t.Run("setDiskAutoDelete", func(t *testing.T) {
		buf.Reset()
		want := "disk autoDelete field updated."

		if err := setDiskAutoDelete(&buf, tc.ProjectID, zone, instanceName, instanceDiskName, true); err != nil {
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
		instanceRegionalDiskName := fmt.Sprintf("test-attach-rw-instance-disk-%v-%v", time.Now().Format("01-02-2006"), r.Int())

		if err := createRegionalDisk(&buf, tc.ProjectID, region, replicaZones, instanceRegionalDiskName, "regions/us-west3/diskTypes/pd-ssd", 20); err != nil {
			t.Fatalf("createRegionalDisk got err: %v", err)
		}

		buf.Reset()
		want := "Disk attached"

		diskUrl := fmt.Sprintf("projects/%s/regions/%s/disks/%s", tc.ProjectID, region, instanceRegionalDiskName)

		if err := attachRegionalDisk(&buf, tc.ProjectID, zone, instanceName, diskUrl); err != nil {
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
			if strings.Contains(*disk.Source, instanceRegionalDiskName) {
				foundDisk = true
			}
		}
		if !foundDisk {
			t.Errorf("The disk %s is not attached to the instance!", instanceRegionalDiskName)
		}

		// Cannot clean up the disk just yet because it must be done after the VM is terminated.
		// It will be done by deleteInstance function.
	})

	t.Run("Attach a read-only regional disk to VM", func(t *testing.T) {
		instanceRegionalDiskName := fmt.Sprintf("test-attach-ro-instance-disk-%v-%v", time.Now().Format("01-02-2006"), r.Int())

		if err := createRegionalDisk(&buf, tc.ProjectID, region, replicaZones, instanceRegionalDiskName, "regions/us-west3/diskTypes/pd-ssd", 20); err != nil {
			t.Fatalf("createRegionalDisk got err: %v", err)
		}

		buf.Reset()
		want := "Disk attached"

		diskUrl := fmt.Sprintf("projects/%s/regions/%s/disks/%s", tc.ProjectID, region, instanceRegionalDiskName)

		if err := attachRegionalDiskReadOnly(&buf, tc.ProjectID, zone, instanceName, diskUrl); err != nil {
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
			if strings.Contains(*disk.Source, instanceRegionalDiskName) {
				foundDisk = true
			}
		}
		if !foundDisk {
			t.Errorf("The disk %s is not attached to the instance!", instanceRegionalDiskName)
		}

		// Cannot clean up the disk just yet because it must be done after the VM is terminated.
		// It will be done by deleteInstance function.
	})
}
