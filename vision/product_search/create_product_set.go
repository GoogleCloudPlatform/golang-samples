// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package productsearch contains samples for Google Cloud Vision API Product Search.
package productsearch

// [START vision_product_search_create_product_set]

import (
	"context"
	"fmt"
	"io"

	vision "cloud.google.com/go/vision/apiv1"
	visionpb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

// createProductSet creates a product set.
func createProductSet(w io.Writer, projectID string, location string, productSetID string, productSetDisplayName string) error {
	ctx := context.Background()
	c, err := vision.NewProductSearchClient(ctx)
	if err != nil {
		return fmt.Errorf("NewProductSearchClient: %v", err)
	}

	req := &visionpb.CreateProductSetRequest{
		Parent:       fmt.Sprintf("projects/%s/locations/%s", projectID, location),
		ProductSetId: productSetID,
		ProductSet: &visionpb.ProductSet{
			DisplayName: productSetDisplayName,
		},
	}

	resp, err := c.CreateProductSet(ctx, req)
	if err != nil {
		return fmt.Errorf("CreateProductSet: %v", err)
	}

	fmt.Fprintf(w, "Product set name: %s\n", resp.Name)

	return nil
}

// [END vision_product_search_create_product_set]
