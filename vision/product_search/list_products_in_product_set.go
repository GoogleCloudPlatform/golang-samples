// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package productsearch contains samples for Google Cloud Vision API Product Search.
package productsearch

// [START vision_product_search_list_products_in_product_set]

import (
	"context"
	"fmt"
	"io"

	vision "cloud.google.com/go/vision/apiv1"
	"google.golang.org/api/iterator"
	visionpb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

// listProductsInProductSet lists products in a product set.
func listProductsInProductSet(w io.Writer, projectID string, location string, productSetID string) error {
	ctx := context.Background()
	c, err := vision.NewProductSearchClient(ctx)
	if err != nil {
		return fmt.Errorf("NewProductSearchClient: %v", err)
	}

	req := &visionpb.ListProductsInProductSetRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/productSets/%s", projectID, location, productSetID),
	}

	it := c.ListProductsInProductSet(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("Next: %v", err)
		}

		fmt.Fprintf(w, "Product name: %s\n", resp.Name)
		fmt.Fprintf(w, "Product display name: %s\n", resp.DisplayName)
		fmt.Fprintf(w, "Product category: %s\n", resp.ProductCategory)
		fmt.Fprintf(w, "Product labels: %s\n", resp.ProductLabels)
	}

	return nil
}

// [END vision_product_search_list_products_in_product_set]
