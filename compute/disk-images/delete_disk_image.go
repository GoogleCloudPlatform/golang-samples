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

// [START compute_images_delete]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
)

// Deletes the specified disk image
func deleteDiskImage(
	w io.Writer,
	projectID, imageName string,
) error {
	// projectID := "your_project_id"
	// imageName := "your_image"

	ctx := context.Background()
	imagesClient, err := compute.NewImagesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewImagesRESTClient: %w", err)
	}
	defer imagesClient.Close()

	req := computepb.DeleteImageRequest{
		Image:   imageName,
		Project: projectID,
	}

	op, err := imagesClient.Delete(ctx, &req)

	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %w", err)
	}

	fmt.Fprintf(w, "Disk image %s deleted\n", imageName)

	return nil
}

// [END compute_images_delete]
