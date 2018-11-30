// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package productsearch contains samples for Google Cloud Vision API Product Search.
package productsearch

// [START vision_product_search_update_product_labels]

import (
	"context"
	"fmt"
	"io"

	vision "cloud.google.com/go/vision/apiv1"
	visionpb "google.golang.org/genproto/googleapis/cloud/vision/v1"
	field_mask "google.golang.org/genproto/protobuf/field_mask"
)

func updateProductLabels(w io.Writer, projectID string, location string, productID string, key string, value string) error {
	ctx := context.Background()
	c, err := vision.NewProductSearchClient(ctx)
	if err != nil {
		return err
	}

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
					Key: key,
					Value: value,
				},
			},
		},
	}

	resp, err := c.UpdateProduct(ctx, req)
	if err != nil {
		return err
	}

	fmt.Fprintln(w, "Product name:", resp.Name)
	fmt.Fprintln(w, "Updated product labels:", resp.ProductLabels)

	return nil
}

// [END vision_product_search_update_product_labels]
