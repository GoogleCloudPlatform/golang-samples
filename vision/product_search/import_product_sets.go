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

// [START vision_product_search_import_product_images]

// [START vision_product_search_tutorial_import]

import (
	"context"
	"fmt"
	"io"

	vision "cloud.google.com/go/vision/apiv1"
	visionpb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

// [END vision_product_search_tutorial_import]

// importProductSets creates a product set using information in a csv file on GCS.
func importProductSets(w io.Writer, projectID string, location string, gcsURI string) error {
	ctx := context.Background()
	c, err := vision.NewProductSearchClient(ctx)
	if err != nil {
		return fmt.Errorf("NewProductSearchClient: %v", err)
	}
	defer c.Close()

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
