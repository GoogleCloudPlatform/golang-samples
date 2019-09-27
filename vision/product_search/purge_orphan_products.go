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

// [START vision_product_search_purge_orphan_products]
import (
	"context"
	"fmt"
	"io"

	vision "cloud.google.com/go/vision/apiv1"
	visionpb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

// purgeOrphanProducts deletes all products not in any product sets.
func purgeOrphanProducts(w io.Writer, projectID string, location string) error {
	// projectID := "your-gcp-project-id"
	// location := "us-west1"

	ctx := context.Background()
	c, err := vision.NewProductSearchClient(ctx)
	if err != nil {
		return fmt.Errorf("NewProductSearchClient: %v", err)
	}
	defer c.Close()

	req := &visionpb.PurgeProductsRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, location),
		Target: &visionpb.PurgeProductsRequest_DeleteOrphanProducts{
			DeleteOrphanProducts: true,
		},
		Force: true,
	}

	// The purge operation is async.
	op, err := c.PurgeProducts(ctx, req)
	if err != nil {
		return fmt.Errorf("NewProductSearchClient: %v", err)
	}
	fmt.Fprintf(w, "Processing operation name: %q\n", op.Name())

	if err := op.Wait(ctx); err != nil {
		return fmt.Errorf("Wait: %v", err)
	}

	fmt.Fprintf(w, "Orphan products deleted.\n")

	return nil
}

// [END vision_product_search_purge_orphan_products]
