// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package productsearch contains samples for Google Cloud Vision API Product Search.
package productsearch

// [START vision_product_search_list_product_sets]

import (
	"context"
	"fmt"
	"io"

	vision "cloud.google.com/go/vision/apiv1"
	"google.golang.org/api/iterator"
	visionpb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

// listProductSets lists product sets.
func listProductSets(w io.Writer, projectID string, location string) error {
	ctx := context.Background()
	c, err := vision.NewProductSearchClient(ctx)
	if err != nil {
		return fmt.Errorf("NewProductSearchClient: %v", err)
	}

	req := &visionpb.ListProductSetsRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, location),
	}

	it := c.ListProductSets(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("Next: %v", err)
		}

		fmt.Fprintf(w, "Product set name: %s\n", resp.Name)
		fmt.Fprintf(w, "Product set display name: %s\n", resp.DisplayName)
		fmt.Fprintf(w, "Product set index time:\n")
		fmt.Fprintf(w, "seconds: %d\n", resp.IndexTime.Seconds)
		fmt.Fprintf(w, "nanos: %d\n", resp.IndexTime.Nanos)
	}

	return nil
}

// [END vision_product_search_list_product_sets]
