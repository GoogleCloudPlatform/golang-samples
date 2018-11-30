// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package productsearch contains samples for Google Cloud Vision API Product Search.
package productsearch

// [START vision_product_search_get_similar_products_gcs]

import (
	"context"
	"fmt"
	"io"

	vision "cloud.google.com/go/vision/apiv1"
	visionpb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

// getSimilarProductsURI searches for products from a product set similar to products in an image file on GCS.
func getSimilarProductsURI(w io.Writer, projectID string, location string, productSetID string, productCategory string, imageURI string, filter string) error {
	ctx := context.Background()
	c, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		fmt.Errorf("NewImageAnnotatorClient: %v", err)
	}

	image := vision.NewImageFromURI(imageURI)

	ictx := &visionpb.ImageContext{
		ProductSearchParams: &visionpb.ProductSearchParams{
			ProductSet: fmt.Sprintf("projects/%s/locations/%s/productSets/%s", projectID, location, productSetID),
			ProductCategories: []string{productCategory},
			Filter: filter,
		},
	}

	response, err := c.ProductSearch(ctx, image, ictx)
	if err != nil {
		fmt.Errorf("ProductSearch: %v", err)
	}

	fmt.Fprintln(w, "Product set index time:")
	fmt.Fprintln(w, "  seconds: ", response.IndexTime.Seconds)
	fmt.Fprintln(w, "  nanos: ", response.IndexTime.Nanos, "\n")

	fmt.Fprintln(w, "Search results:")
	for _, result := range response.Results {
		fmt.Fprintln(w, "Score(Confidence): ", result.Score)
		fmt.Fprintln(w, "Image name: ", result.Image)

		fmt.Fprintln(w, "Prodcut name: ", result.Product.Name)
		fmt.Fprintln(w, "Product display name: ", result.Product.DisplayName)
		fmt.Fprintln(w, "Product labels: ", result.Product.ProductLabels, "\n")
	}

	return nil
}

// [END vision_product_search_get_similar_products_gcs]
