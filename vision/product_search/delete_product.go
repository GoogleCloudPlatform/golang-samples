// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package productsearch contains samples for Google Cloud Vision API Product Search.
package productsearch

// [START vision_product_search_delete_product]

import (
	"context"
	"fmt"
	"io"

	vision "cloud.google.com/go/vision/apiv1"
	visionpb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

func deleteProduct(w io.Writer, projectID string, location string, productID string) error {
	ctx := context.Background()
	c, err := vision.NewProductSearchClient(ctx)
	if err != nil {
		return err
	}

	req := &visionpb.DeleteProductRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/products/%s", projectID, location, productID),
	}

	err = c.DeleteProduct(ctx, req)
	if err != nil {
		return err
	}

	fmt.Fprintln(w, "Product deleted.")

	return nil
}

// [END vision_product_search_delete_product]
