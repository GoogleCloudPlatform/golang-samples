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

// [START compute_disk_clone_encrypted_disk]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/protobuf/proto"
)

// Creates a zonal non-boot persistent disk in a project with the copy of data from an existing disk.
// The encryption key must be the same for the source disk and the new disk.
// The disk type and size may differ.
func createDiskFromCustomerEncryptedDisk(
	w io.Writer,
	projectID, zone, diskName, diskType string,
	diskSizeGb int64,
	diskLink, encryptionKey string,
) error {
	// projectID := "your_project_id"
	// zone := "us-west3-b" // should match diskType below
	// diskName := "your_disk_name"
	// diskType := "zones/us-west3/diskTypes/pd-ssd"
	// diskSizeGb := 120
	// diskLink := "projects/your_project_id/global/disks/disk_name"
	// encryptionKey := "SGVsbG8gZnJvbSBHb29nbGUgQ2xvdWQgUGxhdGZvcm0=" // in base64

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
			Name:       proto.String(diskName),
			Zone:       proto.String(zone),
			Type:       proto.String(diskType),
			SizeGb:     proto.Int64(diskSizeGb),
			SourceDisk: proto.String(diskLink),
			DiskEncryptionKey: &computepb.CustomerEncryptionKey{
				RawKey: &encryptionKey,
			},
		},
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

// [END compute_disk_clone_encrypted_disk]
