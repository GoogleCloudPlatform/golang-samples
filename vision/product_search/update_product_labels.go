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

// [START vision_product_search_update_product_labels]

import (
	"context"
	"fmt"
	"io"

	vision "cloud.google.com/go/vision/apiv1"
	"cloud.google.com/go/vision/v2/apiv1/visionpb"
	field_mask "google.golang.org/genproto/protobuf/field_mask"
)

// updateProductLabels updates product labels of a product.
func updateProductLabels(w io.Writer, projectID string, location string, productID string, key string, value string) error {
	ctx := context.Background()
	c, err := vision.NewProductSearchClient(ctx)
	if err != nil {
		return fmt.Errorf("NewProductSearchClient: %w", err)
	}
	defer c.Close()

	req := &visionpb.UpdateProductRequest{
		UpdateMask: &field_mask.FieldMask{
			Paths: []string{
				"product_labels",
			},
		},
		Product: &visionpb.Product{
			Name: fmt.Sprintf("projects/%s/locations/%s/products/%s", projectID, location, productID),
			ProductLabels: []*visionpb.Product_KeyValue{
				{
					Key:   key,
					Value: value,
				},
			},
		},
	}

	resp, err := c.UpdateProduct(ctx, req)
	if err != nil {
		return fmt.Errorf("UpdateProduct: %w", err)
	}

	fmt.Fprintf(w, "Product name: %s\n", resp.Name)
	fmt.Fprintf(w, "Updated product labels: %s\n", resp.ProductLabels)

	return nil
}

// [END vision_product_search_update_product_labels]
