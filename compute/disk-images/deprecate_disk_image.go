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

// [START compute_images_set_deprecation_status]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/protobuf/proto"
)

// Geg a disk image from the given project
func deprecateDiskImage(
	w io.Writer,
	projectID, imageName string,
) error {
	// projectID := "your_project_id"
	// imageName := "my_image"

	deprecationStatus := &computepb.DeprecationStatus{
		State: proto.String(computepb.DeprecationStatus_DEPRECATED.String()),
	}

	ctx := context.Background()
	imagesClient, err := compute.NewImagesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewImagesRESTClient: %w", err)
	}
	defer imagesClient.Close()

	source_req := &computepb.DeprecateImageRequest{
		Project:                   projectID,
		Image:                     imageName,
		DeprecationStatusResource: deprecationStatus,
	}

	op, err := imagesClient.Deprecate(ctx, source_req)
	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %w", err)
	}

	fmt.Fprintf(w, "Disk image %s deprecated\n", imageName)

	return nil
}

// [END compute_images_set_deprecation_status]
