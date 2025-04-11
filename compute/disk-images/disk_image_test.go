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

func getDisk(ctx context.Context, projectID, diskName, zone string) (*computepb.Disk, error) {
	diskClient, err := compute.NewDisksRESTClient(ctx)
	if err != nil {
		return nil, err
	}
	defer diskClient.Close()

	req := &computepb.GetDiskRequest{
		Project: projectID,
		Disk:    diskName,
		Zone:    zone,
	}
	disk, err := diskClient.Get(ctx, req)
	if err != nil {
		return nil, err
	}
	return disk, nil
}

func createSnapshot(
	ctx context.Context,
	projectID, snapshotName string,
	disk *computepb.Disk,
	locations *[]string,
) error {
	snapshotsClient, err := compute.NewSnapshotsRESTClient(ctx)
	if err != nil {
		return err
	}
	defer snapshotsClient.Close()

	req := &computepb.InsertSnapshotRequest{
		Project: projectID,
		SnapshotResource: &computepb.Snapshot{
			Name:             proto.String(snapshotName),
			SourceDisk:       proto.String(disk.GetSelfLink()),
			StorageLocations: *locations,
		},
	}

	op, err := snapshotsClient.Insert(ctx, req)
	if err != nil {
		return err
	}

	return op.Wait(ctx)
}

func deleteSnapshot(ctx context.Context, projectID, snapshotName string) error {
	snapshotsClient, err := compute.NewSnapshotsRESTClient(ctx)
	if err != nil {
		return err
	}
	defer snapshotsClient.Close()

	req := &computepb.DeleteSnapshotRequest{
		Project:  projectID,
		Snapshot: snapshotName,
	}

	op, err := snapshotsClient.Delete(ctx, req)
	if err != nil {
		return err
	}

	return op.Wait(ctx)
}

func TestComputeDiskImageSnippets(t *testing.T) {
	ctx := context.Background()
	var r *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	tc := testutil.SystemTest(t)
	zone := "us-central1-a"
	imageName := fmt.Sprintf("test-image-go-%v-%v", time.Now().Format("01-02-2006"), r.Int())
	diskName := fmt.Sprintf("test-disk-go-%v-%v", time.Now().Format("01-02-2006"), r.Int())
	sourceImage := "projects/debian-cloud/global/images/family/debian-11"
	snapshotName := fmt.Sprintf("test-snapshot-go-%v-%v", time.Now().Format("01-02-2006"), r.Int())

	buf := &bytes.Buffer{}

	err := createDisk(ctx, tc.ProjectID, zone, diskName, sourceImage)
	if err != nil {
		t.Errorf("createDisk got err: %v", err)
	}

	t.Run("Test snapshot creation and deletion", func(t *testing.T) {
		want := "created"

		if err := createImageFromDisk(buf, tc.ProjectID, zone, diskName, imageName, []string{}, false); err != nil {
			t.Fatalf("createImageFromDisk got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, want) {
			t.Errorf("createImageFromDisk got %q, want %q", got, want)
		}

		want = "deprecated"
		buf.Reset()

		if err := deprecateDiskImage(buf, tc.ProjectID, imageName); err != nil {
			t.Errorf("deprecateDiskImage got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, want) {
			t.Errorf("deprecateDiskImage got %q, want %q", got, want)
		}

		buf.Reset()
		want = "was found"

		err = getDiskImage(buf, tc.ProjectID, imageName)
		if err != nil {
			t.Errorf("getDiskImage got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, want) {
			t.Errorf("getDiskImage got %q, want %q", got, want)
		}

		buf.Reset()
		want = "Newest disk image was found"

		_, err = getDiskImageFromFamily(buf, "debian-cloud", "debian-11")
		if err != nil {
			t.Errorf("getDiskImageFromFamily got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, want) {
			t.Errorf("getDiskImageFromFamily got %q, want %q", got, want)
		}

		buf.Reset()
		want = "deleted"

		if err := deleteDiskImage(buf, tc.ProjectID, imageName); err != nil {
			t.Errorf("deleteDiskImage got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, want) {
			t.Errorf("deleteDiskImage got %q, want %q", got, want)
		}
	})

	t.Run("Test image creation and deletion from snapshot", func(t *testing.T) {
		disk, err := getDisk(ctx, tc.ProjectID, diskName, zone)
		if err != nil {
			t.Fatalf("getDisk got err: %v", err)
		}

		err = createSnapshot(ctx, tc.ProjectID, snapshotName, disk, &[]string{})
		if err != nil {
			t.Fatalf("getDisk got err: %v", err)
		}

		want := "created"
		if err := createImageFromSnapshot(buf, tc.ProjectID, snapshotName, imageName); err != nil {
			t.Fatalf("createImageFromDisk got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, want) {
			t.Errorf("createImageFromDisk got %q, want %q", got, want)
		}

		buf.Reset()
		want = "deleted"

		if err := deleteDiskImage(buf, tc.ProjectID, imageName); err != nil {
			t.Errorf("deleteDiskImage got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, want) {
			t.Errorf("deleteDiskImage got %q, want %q", got, want)
		}
	})

	err = deleteDisk(ctx, tc.ProjectID, zone, diskName)
	if err != nil {
		t.Errorf("deleteDisk got err: %v", err)
	}

	err = deleteSnapshot(ctx, tc.ProjectID, snapshotName)
	if err != nil {
		t.Errorf("deleteSnapshot got err: %v", err)
	}
}
