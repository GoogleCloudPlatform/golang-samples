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

// [START vision_product_search_create_reference_image]

import (
	"context"
	"fmt"
	"io"

	vision "cloud.google.com/go/vision/apiv1"
	visionpb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

// createReferenceImage creates a reference image for a product.
func createReferenceImage(w io.Writer, projectID string, location string, productID string, referenceImageID string, gcsURI string) error {
	ctx := context.Background()
	c, err := vision.NewProductSearchClient(ctx)
	if err != nil {
		return fmt.Errorf("NewProductSearchClient: %v", err)
	}
	defer c.Close()

	req := &visionpb.CreateReferenceImageRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s/products/%s", projectID, location, productID),
		ReferenceImage: &visionpb.ReferenceImage{
			Uri: gcsURI,
		},
		ReferenceImageId: referenceImageID,
	}

	resp, err := c.CreateReferenceImage(ctx, req)
	if err != nil {
		return fmt.Errorf("CreateReferenceImage: %v", err)
	}

	fmt.Fprintf(w, "Reference image name: %s\n", resp.Name)
	fmt.Fprintf(w, "Reference image uri: %s\n", resp.Uri)

	return nil
}

// [END vision_product_search_create_reference_image]
