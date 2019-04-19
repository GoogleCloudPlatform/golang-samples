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

// [START vision_product_search_get_product_set]

import (
	"context"
	"fmt"
	"io"

	vision "cloud.google.com/go/vision/apiv1"
	visionpb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

// getProductSet gets a product set.
func getProductSet(w io.Writer, projectID string, location string, productSetID string) error {
	ctx := context.Background()
	c, err := vision.NewProductSearchClient(ctx)
	if err != nil {
		return fmt.Errorf("NewProductSearchClient: %v", err)
	}
	defer c.Close()

	req := &visionpb.GetProductSetRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/productSets/%s", projectID, location, productSetID),
	}

	resp, err := c.GetProductSet(ctx, req)
	if err != nil {
		return fmt.Errorf("GetProductSet: %v", err)
	}

	fmt.Fprintf(w, "Product set name: %s\n", resp.Name)
	fmt.Fprintf(w, "Product set display name: %s\n", resp.DisplayName)
	fmt.Fprintf(w, "Product set index time:\n")
	fmt.Fprintf(w, "seconds: %d\n", resp.IndexTime.Seconds)
	fmt.Fprintf(w, "nanos: %d\n", resp.IndexTime.Nanos)

	return nil
}

// [END vision_product_search_get_product_set]
