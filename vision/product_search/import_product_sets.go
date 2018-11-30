// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package productsearch contains samples for Google Cloud Vision API Product Search.
package productsearch

// [START vision_product_search_import_product_images]

import (
	"context"
	"fmt"
	"io"

	vision "cloud.google.com/go/vision/apiv1"
	visionpb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

// importProductSets creates a product set using information in a csv file on GCS.
func importProductSets(w io.Writer, projectID string, location string, gcsURI string) error {
	ctx := context.Background()
	c, err := vision.NewProductSearchClient(ctx)
	if err != nil {
		return fmt.Errorf("NewProductSearchClient: %v", err)
	}

	req := &visionpb.ImportProductSetsRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, location),
		InputConfig: &visionpb.ImportProductSetsInputConfig{
			Source: &visionpb.ImportProductSetsInputConfig_GcsSource{
				GcsSource: &visionpb.ImportProductSetsGcsSource{
					CsvFileUri: gcsURI,
				},
			},
		},
	}

	op, err := c.ImportProductSets(ctx, req)
	if err != nil {
		return fmt.Errorf("ImportProductSets: %v", err)
	}

	fmt.Fprintf(w, "Processing operation name: %s\n", op.Name())

	resp, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("Wait: %v", err)
	}

	fmt.Fprintf(w, "processing done.\n")

	for i, status := range resp.Statuses {
		// `0` is the coee for OK in google.rpc.code
		fmt.Fprintf(w, "Status of processing line %d of the csv: %d\n", i, status.Code)

		if status.Code == 0 {
			fmt.Fprintf(w, "Reference image name: %s\n", resp.ReferenceImages[i].Name)
		} else {
			fmt.Fprintf(w, "Status code not OK: %s\n", status.Message)
		}
	}

	return nil
}

// [END vision_product_search_import_product_images]
