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

package productsearch

import (
	"bytes"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestDeleteReferenceImage(t *testing.T) {
	tc := testutil.SystemTest(t)

	const location = "us-west1"
	const productDisplayName = "fake_product_display_name_for_testing"
	const productCategory = "homegoods"
	const productID = "fake_product_id_for_testing"
	const referenceImageID = "fake_reference_image_id_for_testing"
	const gcsURI = "gs://cloud-samples-data/vision/product_search/shoes_1.jpg"

	var buf bytes.Buffer

	// Ensure re-used resource names don't exist prior to test start.
	if err := getReferenceImage(&buf, tc.ProjectID, location, productID, referenceImageID); err == nil {
		deleteReferenceImage(&buf, tc.ProjectID, location, productID, referenceImageID)
	}
	if err := getProduct(&buf, tc.ProjectID, location, productID); err == nil {
		deleteProduct(&buf, tc.ProjectID, location, productID)
	}
	buf.Reset()

	// Create a fake product.
	if err := createProduct(&buf, tc.ProjectID, location, productID, productDisplayName, productCategory); err != nil {
		t.Fatalf("createProduct: %v", err)
	}

	// Create reference image.
	if err := createReferenceImage(&buf, tc.ProjectID, location, productID, referenceImageID, gcsURI); err != nil {
		t.Fatalf("createReferenceImage: %v", err)
	}

	// Confirm the reference image exists.
	buf.Reset()
	if err := listReferenceImages(&buf, tc.ProjectID, location, productID); err != nil {
		t.Fatalf("listReferenceImages: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, referenceImageID) {
		t.Errorf("Reference image ID %s does not exist", referenceImageID)
	}

	// Delete reference image.
	if err := deleteReferenceImage(&buf, tc.ProjectID, location, productID, referenceImageID); err != nil {
		t.Fatalf("deleteReferenceImage: %v", err)
	}

	// Check if the reference image is deleted.
	buf.Reset()
	if err := listReferenceImages(&buf, tc.ProjectID, location, productID); err != nil {
		t.Fatalf("listReferenceImages: %v", err)
	}
	if got := buf.String(); strings.Contains(got, referenceImageID) {
		t.Errorf("Reference image ID %s still exists", referenceImageID)
	}

	// Clean up.
	deleteProduct(&buf, tc.ProjectID, location, productID)
}
