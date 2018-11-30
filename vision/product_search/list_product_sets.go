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

	"google.golang.org/api/iterator"
	vision "cloud.google.com/go/vision/apiv1"
	visionpb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

// listProductSets lists product sets.
func listProductSets(w io.Writer, projectID string, location string) error {
	ctx := context.Background()
	c, err := vision.NewProductSearchClient(ctx)
	if err != nil {
		fmt.Errorf("NewProductSearchClient: %v", err)
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
			fmt.Errorf("Next: %v", err)
		}

		fmt.Fprintln(w, "Product set name:", resp.Name)
		fmt.Fprintln(w, "Product set display name:", resp.DisplayName)
		fmt.Fprintln(w, "Product set index time:")
		fmt.Fprintln(w, "  seconds: ", resp.IndexTime.Seconds)
		fmt.Fprintln(w, "  nanos: ", resp.IndexTime.Nanos, "\n")
	}

	return nil
}

// [END vision_product_search_list_product_sets]
