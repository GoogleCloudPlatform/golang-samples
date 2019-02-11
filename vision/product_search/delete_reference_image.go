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

// [START vision_product_search_delete_reference_image]

import (
	"context"
	"fmt"
	"io"

	vision "cloud.google.com/go/vision/apiv1"
	visionpb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

// deleteReferenceImage deletes a reference image from a product.
func deleteReferenceImage(w io.Writer, projectID string, location string, productID string, referenceImageID string) error {
	ctx := context.Background()
	c, err := vision.NewProductSearchClient(ctx)
	if err != nil {
		return fmt.Errorf("NewProductSearchClient: %v", err)
	}

	req := &visionpb.DeleteReferenceImageRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/products/%s/referenceImages/%s", projectID, location, productID, referenceImageID),
	}

	if err = c.DeleteReferenceImage(ctx, req); err != nil {
		return fmt.Errorf("NewProductSearchClient: %v", err)
	}

	fmt.Fprintf(w, "Reference image deleted from product.\n")

	return nil
}

// [END vision_product_search_delete_reference_image]
