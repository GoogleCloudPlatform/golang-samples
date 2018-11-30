// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package productsearch contains samples for Google Cloud Vision API Product Search.
package productsearch

// [START vision_product_search_list_reference_images]

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/api/iterator"
	vision "cloud.google.com/go/vision/apiv1"
	visionpb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

func listReferenceImages(w io.Writer, projectID string, location string, productID string) error {
	ctx := context.Background()
	c, err := vision.NewProductSearchClient(ctx)
	if err != nil {
		return err
	}

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
			return err
		}

		fmt.Fprintln(w, "Reference image name:", resp.Name)
		fmt.Fprintln(w, "Reference image uri:", resp.Uri)
		fmt.Fprintln(w, "Reference image bounding polygons: ", resp.BoundingPolys, "\n")
	}

	return nil
}

// [END vision_product_search_list_reference_images]
