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
	ctx := context.Background()
	var r *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	tc := testutil.SystemTest(t)
	zone := "europe-central2-b"
	instanceName := fmt.Sprintf("test-instance-%v-%v", time.Now().Format("01-02-2006"), r.Int())
	diskName := fmt.Sprintf("test-disk-%v-%v", time.Now().Format("01-02-2006"), r.Int())
	diskName2 := fmt.Sprintf("test-disk-%v-%v", time.Now().Format("01-02-2006"), r.Int())
	snapshotName := fmt.Sprintf("test-snapshot-%v-%v", time.Now().Format("01-02-2006"), r.Int())
	sourceImage := "projects/debian-cloud/global/images/family/debian-11"
	diskType := fmt.Sprintf("zones/%s/diskTypes/pd-ssd", zone)
	diskSnapshotLink := fmt.Sprintf("projects/%s/global/snapshots/%s", tc.ProjectID, snapshotName)
	want := "Disk created"

	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		t.Fatalf("NewInstancesRESTClient: %v", err)
	}

	defer instancesClient.Close()

	buf := &bytes.Buffer{}

	t.Run("createDiskSnapshot and deleteDisk", func(t *testing.T) {
		err := createDisk(ctx, tc.ProjectID, zone, diskName, sourceImage)
		if err != nil {
			t.Fatalf("createDisk got err: %v", err)
		}

		err = createDiskSnapshot(ctx, tc.ProjectID, zone, diskName, snapshotName)
		if err != nil {
			t.Fatalf("createDiskSnapshot got err: %v", err)
		}

		if err := createDiskFromSnapshot(buf, tc.ProjectID, zone, diskName2, diskType, diskSnapshotLink); err != nil {
			t.Errorf("createDiskFromSnapshot got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, want) {
			t.Errorf("createDiskFromSnapshot got %q, want %q", got, want)
		}

		buf.Reset()
		want = "Disk deleted"

		if err := deleteDisk(buf, tc.ProjectID, zone, diskName2); err != nil {
			t.Errorf("deleteDisk got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, want) {
			t.Errorf("deleteDisk got %q, want %q", got, want)
		}

		err = deleteDiskSnapshot(ctx, tc.ProjectID, snapshotName)
		if err != nil {
			t.Errorf("deleteDiskSnapshot got err: %v", err)
		}

		err = deleteDisk(buf, tc.ProjectID, zone, diskName)
		if err != nil {
			t.Errorf("deleteDisk got err: %v", err)
		}
	})

	t.Run("createEmptyDisk", func(t *testing.T) {
		buf.Reset()
		want = "Disk created"

		if err := createEmptyDisk(buf, tc.ProjectID, zone, diskName, diskType); err != nil {
			t.Fatalf("createEmptyDisk got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, want) {
			t.Errorf("createEmptyDisk got %q, want %q", got, want)
		}

		err = deleteDisk(buf, tc.ProjectID, zone, diskName)
		if err != nil {
			t.Errorf("deleteDisk got err: %v", err)
		}
	})

	t.Run("setDiskAutoDelete", func(t *testing.T) {
		buf.Reset()
		want = "disk autoDelete field updated."

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
			t.Fatalf("unable to create instance: %v", err)
		}

		if err = op.Wait(ctx); err != nil {
			t.Fatalf("unable to wait for the operation: %v", err)
		}

		if err := setDiskAutoDelete(buf, tc.ProjectID, zone, instanceName, diskName); err != nil {
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

		err = deleteInstance(ctx, tc.ProjectID, zone, instanceName)
		if err != nil {
			t.Errorf("deleteInstance got err: %v", err)
		}
	})
}
