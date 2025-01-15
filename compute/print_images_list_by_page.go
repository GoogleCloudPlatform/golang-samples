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

// [START compute_images_list_page]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/api/iterator"
	"google.golang.org/protobuf/proto"
)

// printImagesListByPage prints a list of all non-deprecated image names available in a given project,
// divided into pages as returned by the Compute Engine API.
func printImagesListByPage(w io.Writer, projectID string, pageSize uint32) error {
	// projectID := "your_project_id"
	// pageSize := 10
	ctx := context.Background()
	imagesClient, err := compute.NewImagesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewImagesRESTClient: %w", err)
	}
	defer imagesClient.Close()

	// Listing only non-deprecated images to reduce the size of the reply.
	req := &computepb.ListImagesRequest{
		Project:    projectID,
		MaxResults: proto.Uint32(pageSize),
		Filter:     proto.String("deprecated.state != DEPRECATED"),
	}

	// Use the `iterator.NewPage` to have more granular control of iteration over
	// paginated results from the API. Each time you want to access the
	// next page, the library retrieves that page from the API.
	it := imagesClient.List(ctx, req)
	p := iterator.NewPager(it, int(pageSize), "" /* start from the beginning */)
	for page := 0; ; page++ {
		var items []*computepb.Image
		pageToken, err := p.NextPage(&items)
		if err != nil {
			return fmt.Errorf("iterator paging failed: %w", err)
		}
		fmt.Fprintf(w, "Page %d: %v\n", page, items)
		if pageToken == "" {
			break
		}
	}

	return nil
}

// [END compute_images_list_page]
