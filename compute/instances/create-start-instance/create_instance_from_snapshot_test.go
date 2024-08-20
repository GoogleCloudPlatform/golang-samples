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

func TestComputeCreateInstanceFromSnapshotSnippets(t *testing.T) {
	ctx := context.Background()
	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	tc := testutil.SystemTest(t)
	zone := "europe-central2-b"
	instanceName := "test-" + fmt.Sprint(seededRand.Int())
	instanceName2 := "test-" + fmt.Sprint(seededRand.Int())
	diskName := "test-disk-" + fmt.Sprint(seededRand.Int())
	snapshotName := "test-snapshot-" + fmt.Sprint(seededRand.Int())
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

	buf := &bytes.Buffer{}

	if err := createInstanceFromPublicImage(buf, tc.ProjectID, zone, instanceName); err != nil {
		t.Errorf("createInstanceFromPublicImage got err: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("createInstanceFromPublicImage got %q, want %q", got, expectedResult)
	}

	buf.Reset()

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

	diskSnapshotLink := fmt.Sprintf("projects/%s/global/snapshots/%s", tc.ProjectID, snapshotName)

	if err := createInstanceFromSnapshot(buf, tc.ProjectID, zone, instanceName2, diskSnapshotLink); err != nil {
		t.Errorf("createInstanceFromSnapshot got err: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("createInstanceFromSnapshot got %q, want %q", got, expectedResult)
	}

	err = deleteInstance(ctx, tc.ProjectID, zone, instanceName)
	if err != nil {
		t.Errorf("deleteInstance got err: %v", err)
	}
	err = deleteInstance(ctx, tc.ProjectID, zone, instanceName2)
	if err != nil {
		t.Errorf("deleteInstance got err: %v", err)
	}

	buf.Reset()

	if err := createInstanceWithSnapshottedDataDisk(buf, tc.ProjectID, zone, instanceName, diskSnapshotLink); err != nil {
		t.Errorf("createInstanceWithSnapshottedDataDisk got err: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("createInstanceWithSnapshottedDataDisk got %q, want %q", got, expectedResult)
	}
	err = deleteInstance(ctx, tc.ProjectID, zone, instanceName)
	if err != nil {
		t.Errorf("deleteInstance got err: %v", err)
	}

	buf.Reset()

	deleteSnapshotReq := &computepb.DeleteSnapshotRequest{
		Project:  tc.ProjectID,
		Snapshot: snapshotName,
	}
	op, err = snapshotsClient.Delete(ctx, deleteSnapshotReq)
	if err != nil {
		t.Errorf("unable to delete disk snapshot: %v", err)
	}

	if err = op.Wait(ctx); err != nil {
		t.Errorf("unable to wait for the operation: %v", err)
	}

	deleteDiskReq := &computepb.DeleteDiskRequest{
		Project: tc.ProjectID,
		Disk:    diskName,
		Zone:    zone,
	}
	op, err = disksClient.Delete(ctx, deleteDiskReq)
	if err != nil {
		t.Errorf("unable to delete disk: %v", err)
	}

	if err = op.Wait(ctx); err != nil {
		t.Errorf("unable to wait for the operation: %v", err)
	}
}
