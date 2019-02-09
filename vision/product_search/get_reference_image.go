// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package productsearch contains samples for Google Cloud Vision API Product Search.
package productsearch

// [START vision_product_search_get_reference_image]

import (
	"context"
	"fmt"
	"io"

	vision "cloud.google.com/go/vision/apiv1"
	visionpb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

// getReferenceImage gets a reference image.
func getReferenceImage(w io.Writer, projectID string, location string, productID string, referenceImageID string) error {
	ctx := context.Background()
	c, err := vision.NewProductSearchClient(ctx)
	if err != nil {
		return fmt.Errorf("NewProductSearchClient: %v", err)
	}

	req := &visionpb.GetReferenceImageRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/products/%s/referenceImages/%s", projectID, location, productID, referenceImageID),
	}

	resp, err := c.GetReferenceImage(ctx, req)
	if err != nil {
		return fmt.Errorf("GetReferenceImage: %v", err)
	}

	fmt.Fprintf(w, "Reference image name: %s\n", resp.Name)
	fmt.Fprintf(w, "Reference image uri: %s\n", resp.Uri)
	fmt.Fprintf(w, "Reference image bounding polygons: %s\n", resp.BoundingPolys)

	return nil
}

// [END vision_product_search_get_reference_image]
