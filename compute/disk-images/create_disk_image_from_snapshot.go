// Copyright 2024 Google LLC
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

// [START compute_images_create_from_snapshot]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
)

// Creates a disk image from an existing image
func createImageFromSnapshot(
	w io.Writer,
	projectID, snapshotName, imageName string,
) error {
	// projectID := "your_project_id"
	// snapshotName := "your_image_name"
	// imageName := "my_image"

	// // If storageLocations empty, automatically selects the closest one to the source
	storageLocations := []string{}

	ctx := context.Background()
	imagesClient, err := compute.NewImagesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewImagesRESTClient: %w", err)
	}
	defer imagesClient.Close()

	snapshotsClient, err := compute.NewSnapshotsRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewSnapshotsRESTClient: %w", err)
	}
	defer snapshotsClient.Close()

	// Get the source image
	source_req := &computepb.GetSnapshotRequest{
		Snapshot: snapshotName,
		Project:  projectID,
	}

	snapshot, err := snapshotsClient.Get(ctx, source_req)
	if err != nil {
		return fmt.Errorf("unable to get source image: %w", err)
	}

	// Create the image
	req := computepb.InsertImageRequest{
		ImageResource: &computepb.Image{
			Name:             &imageName,
			SourceSnapshot:   snapshot.SelfLink,
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

// [END compute_images_create_from_snapshot]
