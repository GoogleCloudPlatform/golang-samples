// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package productsearch contains samples for Google Cloud Vision API Product Search.
package productsearch

// [START vision_product_search_get_product]

import (
	"context"
	"fmt"
	"io"

	vision "cloud.google.com/go/vision/apiv1"
	visionpb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

// getProduct gets a product.
func getProduct(w io.Writer, projectID string, location string, productID string) error {
	ctx := context.Background()
	c, err := vision.NewProductSearchClient(ctx)
	if err != nil {
		return fmt.Errorf("NewProductSearchClient: %v", err)
	}

	req := &visionpb.GetProductRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/products/%s", projectID, location, productID),
	}

	resp, err := c.GetProduct(ctx, req)
	if err != nil {
		return fmt.Errorf("GetProduct: %v", err)
	}

	fmt.Fprintf(w, "Product name: %s\n", resp.Name)
	fmt.Fprintf(w, "Product display name: %s\n", resp.DisplayName)
	fmt.Fprintf(w, "Product category: %s\n", resp.ProductCategory)
	fmt.Fprintf(w, "Product labels: %s\n", resp.ProductLabels)

	return nil
}

// [END vision_product_search_get_product]
