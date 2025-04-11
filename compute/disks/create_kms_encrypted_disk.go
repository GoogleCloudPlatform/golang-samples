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

// [START compute_create_kms_encrypted_disk]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/protobuf/proto"
)

// createKmsEncryptedDisk creates a new disk in a project in given zone.
// If you do not provide diskLink or imageLink, an empty disk will be created.
func createKmsEncryptedDisk(
	w io.Writer,
	projectID, zone, diskName, diskType string,
	diskSizeGb int64,
	kmsKeyLink, imageLink, diskLink, snapshotLink string,
) error {
	// projectID := "your_project_id"
	// zone := "us-west3-b" // should match diskType below
	// diskName := "your_disk_name"
	// diskType := "zones/us-west3/diskTypes/pd-ssd"
	// diskSizeGb := 120
	// kmsKeyLink := "projects/your_kms_project_id/locations/us-central1/keyRings/your_key_ring/cryptoKeys/your_key"
	// // Only use one of these at a time
	// diskLink := "projects/your_project_id/global/disks/disk_name"
	// imageLink := "projects/your_project_id/global/images/image_name"
	// snapshotLink := "projects/your_project_id/global/snapshots/snapshot_name"

	ctx := context.Background()
	disksClient, err := compute.NewDisksRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewDisksRESTClient: %w", err)
	}
	defer disksClient.Close()

	req := &computepb.InsertDiskRequest{
		Project: projectID,
		Zone:    zone,
		DiskResource: &computepb.Disk{
			Name:   proto.String(diskName),
			Zone:   proto.String(zone),
			Type:   proto.String(diskType),
			SizeGb: proto.Int64(diskSizeGb),
			DiskEncryptionKey: &computepb.CustomerEncryptionKey{
				KmsKeyName: &kmsKeyLink,
			},
		},
	}

	// if a source disk, image or snapshot has been specified, apply it.
	// These arguments are mutually exclusive.
	if diskLink != "" {
		req.DiskResource.SourceDisk = proto.String(diskLink)
	} else if imageLink != "" {
		req.DiskResource.SourceImage = proto.String(imageLink)
	} else if snapshotLink != "" {
		req.DiskResource.SourceSnapshot = proto.String(snapshotLink)
	}

	op, err := disksClient.Insert(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to create disk: %w", err)
	}

	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %w", err)
	}

	fmt.Fprintf(w, "Disk created\n")

	return nil
}

// [END compute_create_kms_encrypted_disk]
