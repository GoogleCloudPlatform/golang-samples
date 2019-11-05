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

// [START vision_product_search_purge_products_in_product_set]
import (
	"context"
	"fmt"
	"io"

	vision "cloud.google.com/go/vision/apiv1"
	visionpb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

// purgeProductsInProductSet deletes all products in a product set.
func purgeProductsInProductSet(w io.Writer, projectID string, location string, productSetID string) error {
	// projectID := "your-gcp-project-id"
	// location := "us-west1"
	// productSetID := "sampleProductSetID"

	ctx := context.Background()
	c, err := vision.NewProductSearchClient(ctx)
	if err != nil {
		return fmt.Errorf("NewProductSearchClient: %v", err)
	}
	defer c.Close()

	req := &visionpb.PurgeProductsRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, location),
		Target: &visionpb.PurgeProductsRequest_ProductSetPurgeConfig{
			ProductSetPurgeConfig: &visionpb.ProductSetPurgeConfig{
				ProductSetId: productSetID,
			},
		},
		Force: true,
	}

	// The purge operation is async.
	op, err := c.PurgeProducts(ctx, req)
	if err != nil {
		return fmt.Errorf("PurgeProducts: %v", err)
	}
	fmt.Fprintf(w, "Processing operation name: %q\n", op.Name())

	if err := op.Wait(ctx); err != nil {
		return fmt.Errorf("Wait: %v", err)
	}

	fmt.Fprintf(w, "Deleted products in product set.\n")

	return nil
}

// [END vision_product_search_purge_products_in_product_set]
