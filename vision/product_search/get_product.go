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

// [START vision_product_search_get_product]

import (
	"context"
	"fmt"
	"io"

	vision "cloud.google.com/go/vision/apiv1"
	"cloud.google.com/go/vision/v2/apiv1/visionpb"
)

// getProduct gets a product.
func getProduct(w io.Writer, projectID string, location string, productID string) error {
	ctx := context.Background()
	c, err := vision.NewProductSearchClient(ctx)
	if err != nil {
		return fmt.Errorf("NewProductSearchClient: %w", err)
	}
	defer c.Close()

	req := &visionpb.GetProductRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/products/%s", projectID, location, productID),
	}

	resp, err := c.GetProduct(ctx, req)
	if err != nil {
		return fmt.Errorf("GetProduct: %w", err)
	}

	fmt.Fprintf(w, "Product name: %s\n", resp.Name)
	fmt.Fprintf(w, "Product display name: %s\n", resp.DisplayName)
	fmt.Fprintf(w, "Product category: %s\n", resp.ProductCategory)
	fmt.Fprintf(w, "Product labels: %s\n", resp.ProductLabels)

	return nil
}

// [END vision_product_search_get_product]
