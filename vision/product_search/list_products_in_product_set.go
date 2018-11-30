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

	"google.golang.org/api/iterator"
	vision "cloud.google.com/go/vision/apiv1"
	visionpb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

func listProductsInProductSet(w io.Writer, projectID string, location string, productSetID string) error {
	ctx := context.Background()
	c, err := vision.NewProductSearchClient(ctx)
	if err != nil {
		return err
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
			return err
		}

		fmt.Fprintln(w, "Product name:", resp.Name)
		fmt.Fprintln(w, "Product display name:", resp.DisplayName)
		fmt.Fprintln(w, "Product category:", resp.ProductCategory)
		fmt.Fprintln(w, "Product labels:", resp.ProductLabels, "\n")
	}

	return nil
}

// [END vision_product_search_list_products_in_product_set]
