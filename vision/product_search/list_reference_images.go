// Copyright 2019 Google LLC
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

// Package productsearch contains samples for Google Cloud Vision API Product Search.
package productsearch

// [START vision_product_search_list_reference_images]

import (
	"context"
	"fmt"
	"io"

	vision "cloud.google.com/go/vision/apiv1"
	"google.golang.org/api/iterator"
	visionpb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

// listReferenceImages lists reference images of a product.
func listReferenceImages(w io.Writer, projectID string, location string, productID string) error {
	ctx := context.Background()
	c, err := vision.NewProductSearchClient(ctx)
	if err != nil {
		return fmt.Errorf("NewProductSearchClient: %v", err)
	}
	defer c.Close()

	req := &visionpb.ListReferenceImagesRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s/products/%s", projectID, location, productID),
	}

	it := c.ListReferenceImages(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("Next: %v", err)
		}

		fmt.Fprintf(w, "Reference image name: %s\n", resp.Name)
		fmt.Fprintf(w, "Reference image uri: %s\n", resp.Uri)
		fmt.Fprintf(w, "Reference image bounding polygons: %s\n", resp.BoundingPolys)
	}

	return nil
}

// [END vision_product_search_list_reference_images]
