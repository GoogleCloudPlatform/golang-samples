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

// [START compute_images_get_from_family]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
)

// Geg a disk image from family for the given project
func getDiskImageFromFamily(
	w io.Writer,
	projectID, family string,
) (*computepb.Image, error) {
	// projectID := "your_project_id"
	// family := "my_family"

	ctx := context.Background()
	imagesClient, err := compute.NewImagesRESTClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewImagesRESTClient: %w", err)
	}
	defer imagesClient.Close()

	source_req := &computepb.GetFromFamilyImageRequest{
		Project: projectID,
		Family:  family,
	}

	newestImage, err := imagesClient.GetFromFamily(ctx, source_req)
	if err != nil {
		return nil, fmt.Errorf("unable to get image: %w", err)
	}

	fmt.Fprintf(w, "Newest disk image was found: %s\n", *newestImage.Name)

	return newestImage, nil
}

// [END compute_images_get_from_family]
