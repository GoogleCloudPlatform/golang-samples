// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package productsearch contains samples for Google Cloud Vision API Product Search.
package productsearch

// [START vision_product_search_create_product]

import (
	"context"
	"fmt"
	"io"

	vision "cloud.google.com/go/vision/apiv1"
	visionpb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

// createProduct creates a product.
func createProduct(w io.Writer, projectID string, location string, productID string, productDisplayName string, productCategory string) error {
	ctx := context.Background()
	c, err := vision.NewProductSearchClient(ctx)
	if err != nil {
		return fmt.Errorf("NewProductSearchClient: %v", err)
	}

	req := &visionpb.CreateProductRequest{
		Parent:    fmt.Sprintf("projects/%s/locations/%s", projectID, location),
		ProductId: productID,
		Product: &visionpb.Product{
			DisplayName:     productDisplayName,
			ProductCategory: productCategory,
		},
	}

	resp, err := c.CreateProduct(ctx, req)
	if err != nil {
		return fmt.Errorf("CreateProduct: %v", err)
	}

	fmt.Fprintf(w, "Product name: %s\n", resp.Name)

	return nil
}

// [END vision_product_search_create_product]
