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

// [START vision_product_search_create_product]

import (
	"context"
	"fmt"
	"io"

	vision "cloud.google.com/go/vision/apiv1"
	"cloud.google.com/go/vision/v2/apiv1/visionpb"
)

// createProduct creates a product.
func createProduct(w io.Writer, projectID string, location string, productID string, productDisplayName string, productCategory string) error {
	ctx := context.Background()
	c, err := vision.NewProductSearchClient(ctx)
	if err != nil {
		return fmt.Errorf("NewProductSearchClient: %w", err)
	}
	defer c.Close()

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
		return fmt.Errorf("CreateProduct: %w", err)
	}

	fmt.Fprintf(w, "Product name: %s\n", resp.Name)

	return nil
}

// [END vision_product_search_create_product]
