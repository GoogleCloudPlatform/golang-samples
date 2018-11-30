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

func getProduct(w io.Writer, projectID string, location string, productID string) error {
	ctx := context.Background()
	c, err := vision.NewProductSearchClient(ctx)
	if err != nil {
		return err
	}

	req := &visionpb.GetProductRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/products/%s", projectID, location, productID),
	}

	resp, err := c.GetProduct(ctx, req)
	if err != nil {
		return err
	}

	fmt.Fprintln(w, "Product name:", resp.Name)
	fmt.Fprintln(w, "Product display name:", resp.DisplayName)
	fmt.Fprintln(w, "Product category:", resp.ProductCategory)
	fmt.Fprintln(w, "Product labels:", resp.ProductLabels)

	return nil
}

// [END vision_product_search_get_product]
