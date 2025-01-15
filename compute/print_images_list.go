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

// [START compute_images_list]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/api/iterator"
	"google.golang.org/protobuf/proto"
)

// printImagesList prints a list of all non-deprecated image names available in given project.
func printImagesList(w io.Writer, projectID string) error {
	// projectID := "your_project_id"
	ctx := context.Background()
	imagesClient, err := compute.NewImagesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewImagesRESTClient: %w", err)
	}
	defer imagesClient.Close()

	// Listing only non-deprecated images to reduce the size of the reply.
	req := &computepb.ListImagesRequest{
		Project:    projectID,
		MaxResults: proto.Uint32(3),
		Filter:     proto.String("deprecated.state != DEPRECATED"),
	}

	// Although the `MaxResults` parameter is specified in the request, the iterator returned
	// by the `list()` method hides the pagination mechanic. The library makes multiple
	// requests to the API for you, so you can simply iterate over all the images.
	it := imagesClient.List(ctx, req)
	for {
		image, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "- %s\n", image.GetName())
	}
	return nil
}

// [END compute_images_list]
