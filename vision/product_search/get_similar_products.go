// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package productsearch contains samples for Google Cloud Vision API Product Search.
package productsearch

// [START vision_product_search_get_similar_products]

import (
	"context"
	"fmt"
	"io"
	"os"

	vision "cloud.google.com/go/vision/apiv1"
	"cloud.google.com/go/vision/v2/apiv1/visionpb"
)

// getSimilarProducts searches for products from a product set similar to products in an image file.
func getSimilarProducts(w io.Writer, projectID string, location string, productSetID string, productCategory string, file string, filter string) error {
	ctx := context.Background()
	c, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return fmt.Errorf("NewImageAnnotatorClient: %w", err)
	}
	defer c.Close()

	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("Open: %w", err)
	}
	defer f.Close()

	image, err := vision.NewImageFromReader(f)
	if err != nil {
		return fmt.Errorf("NewImageFromReader: %w", err)
	}

	ictx := &visionpb.ImageContext{
		ProductSearchParams: &visionpb.ProductSearchParams{
			ProductSet:        fmt.Sprintf("projects/%s/locations/%s/productSets/%s", projectID, location, productSetID),
			ProductCategories: []string{productCategory},
			Filter:            filter,
		},
	}

	response, err := c.ProductSearch(ctx, image, ictx)
	if err != nil {
		return fmt.Errorf("ProductSearch: %w", err)
	}

	fmt.Fprintf(w, "Product set index time:\n")
	fmt.Fprintf(w, "seconds: %d\n", response.IndexTime.Seconds)
	fmt.Fprintf(w, "nanos: %d\n", response.IndexTime.Nanos)

	fmt.Fprintf(w, "Search results:\n")
	for _, result := range response.Results {
		fmt.Fprintf(w, "Score(Confidence): %f\n", result.Score)
		fmt.Fprintf(w, "Image name: %s\n", result.Image)

		fmt.Fprintf(w, "Prodcut name: %s\n", result.Product.Name)
		fmt.Fprintf(w, "Product display name: %s\n", result.Product.DisplayName)
		fmt.Fprintf(w, "Product labels: %s\n", result.Product.ProductLabels)
	}

	return nil
}

// [END vision_product_search_get_similar_products]
