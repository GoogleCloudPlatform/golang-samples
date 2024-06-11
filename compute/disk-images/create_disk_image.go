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

// [START compute_windows_image_create]
// [START compute_images_create]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
)

// Creates a disk image from an existing disk
func createImageFromDisk(
	w io.Writer,
	projectID, zone, sourceDiskName, imageName string,
	storageLocations []string,
	forceCreate bool,
) error {
	// projectID := "your_project_id"
	// zone := "us-central1-a"
	// sourceDiskName := "your_disk_name"
	// imageName := "my_image"
	// // If storageLocations empty, automatically selects the closest one to the source
	// storageLocations = []string{}
	// // If forceCreate is set to `true`, proceeds even if the disk is attached to
	// // a running instance. This may compromise integrity of the image!
	// forceCreate = false

	ctx := context.Background()
	disksClient, err := compute.NewDisksRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewDisksRESTClient: %w", err)
	}
	defer disksClient.Close()
	imagesClient, err := compute.NewImagesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewImagesRESTClient: %w", err)
	}
	defer imagesClient.Close()

	// Get the source disk
	source_req := &computepb.GetDiskRequest{
		Disk:    sourceDiskName,
		Project: projectID,
		Zone:    zone,
	}

	disk, err := disksClient.Get(ctx, source_req)
	if err != nil {
		return fmt.Errorf("unable to get source disk: %w", err)
	}

	// Create the image
	req := computepb.InsertImageRequest{
		ForceCreate: &forceCreate,
		ImageResource: &computepb.Image{
			Name:             &imageName,
			SourceDisk:       disk.SelfLink,
			StorageLocations: storageLocations,
		},
		Project: projectID,
	}

	op, err := imagesClient.Insert(ctx, &req)

	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %w", err)
	}

	fmt.Fprintf(w, "Disk image %s created\n", imageName)

	return nil
}

// [END compute_windows_image_create]
// [END compute_images_create]
