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

func createRegionDisk(ctx context.Context, projectId, region, diskName string) error {
	regionDisksClient, err := compute.NewRegionDisksRESTClient(ctx)
	if err != nil {
		return err
	}
	defer regionDisksClient.Close()
	req := &computepb.InsertRegionDiskRequest{
		Project: projectId,
		Region:  region,
		DiskResource: &computepb.Disk{
			Name:   proto.String(diskName),
			SizeGb: proto.Int64(200),
			ReplicaZones: []string{
				fmt.Sprintf("projects/%s/zones/europe-central2-a", projectId),
				fmt.Sprintf("projects/%s/zones/europe-central2-b", projectId),
			},
		},
	}

	op, err := regionDisksClient.Insert(ctx, req)
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

func deleteRegionDisk(ctx context.Context, projectId, region, diskName string) error {
	regionDisksClient, err := compute.NewRegionDisksRESTClient(ctx)
	if err != nil {
		return err
	}
	defer regionDisksClient.Close()
	req := &computepb.DeleteRegionDiskRequest{
		Project: projectId,
		Region:  region,
		Disk:    diskName,
	}

	op, err := regionDisksClient.Delete(ctx, req)
	if err != nil {
		return err
	}

	return op.Wait(ctx)
}

func TestComputeSnapshotsSnippets(t *testing.T) {
	ctx := context.Background()
	var r *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	tc := testutil.SystemTest(t)
	zone := "europe-central2-b"
	location := "europe-central2"
	snapshotName := fmt.Sprintf("test-snapshot-%v-%v", time.Now().Format("01-02-2006"), r.Int())
	diskName := fmt.Sprintf("test-disk-%v-%v", time.Now().Format("01-02-2006"), r.Int())
	sourceImage := "projects/debian-cloud/global/images/family/debian-11"
	want := "Snapshot created"

	buf := &bytes.Buffer{}

	err := createDisk(ctx, tc.ProjectID, zone, diskName, sourceImage)
	if err != nil {
		t.Fatalf("createDisk got err: %v", err)
	}

	if err := createSnapshot(buf, tc.ProjectID, diskName, snapshotName, zone, "", location, ""); err != nil {
		t.Fatalf("createSnapshot got err: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("createSnapshot got %q, want %q", got, want)
	}

	buf.Reset()
	want = "Snapshot deleted"

	if err := deleteSnapshot(buf, tc.ProjectID, snapshotName); err != nil {
		t.Errorf("deleteSnapshot got err: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("deleteSnapshot got %q, want %q", got, want)
	}

	err = deleteDisk(ctx, tc.ProjectID, zone, diskName)
	if err != nil {
		t.Errorf("deleteDisk got err: %v", err)
	}

	buf.Reset()
	want = "Snapshot created"

	err = createRegionDisk(ctx, tc.ProjectID, location, diskName)
	if err != nil {
		t.Fatalf("createRegionDisk got err: %v", err)
	}

	if err := createSnapshot(buf, tc.ProjectID, diskName, snapshotName, "", location, location, ""); err != nil {
		t.Fatalf("createSnapshot got err: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("createSnapshot got %q, want %q", got, want)
	}

	buf.Reset()
	want = "Snapshot deleted"

	if err := deleteSnapshot(buf, tc.ProjectID, snapshotName); err != nil {
		t.Errorf("deleteSnapshot got err: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("deleteSnapshot got %q, want %q", got, want)
	}

	err = deleteRegionDisk(ctx, tc.ProjectID, location, diskName)
	if err != nil {
		t.Errorf("deleteRegionalDisk got err: %v", err)
	}
}
