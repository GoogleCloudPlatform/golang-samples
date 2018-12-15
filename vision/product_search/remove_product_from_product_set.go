// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package productsearch contains samples for Google Cloud Vision API Product Search.
package productsearch

// [START vision_product_search_remove_product_from_product_set]

import (
	"context"
	"fmt"
	"io"

	vision "cloud.google.com/go/vision/apiv1"
	visionpb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

// removeProductFromProductSet removes a product from a product set.
func removeProductFromProductSet(w io.Writer, projectID string, location string, productID string, productSetID string) error {
	ctx := context.Background()
	c, err := vision.NewProductSearchClient(ctx)
	if err != nil {
		return fmt.Errorf("NewProductSearchClient: %v", err)
	}

	req := &visionpb.RemoveProductFromProductSetRequest{
		Name:    fmt.Sprintf("projects/%s/locations/%s/productSets/%s", projectID, location, productSetID),
		Product: fmt.Sprintf("projects/%s/locations/%s/products/%s", projectID, location, productID),
	}

	if err = c.RemoveProductFromProductSet(ctx, req); err != nil {
		return fmt.Errorf("NewProductSearchClient: %v", err)
	}

	fmt.Fprintf(w, "Product removed from product set.\n")

	return nil
}

// [END vision_product_search_remove_product_from_product_set]
