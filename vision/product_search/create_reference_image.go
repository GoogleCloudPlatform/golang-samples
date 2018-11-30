// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

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
