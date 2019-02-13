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

func TestCreateReferenceImage(t *testing.T) {
	tc := testutil.SystemTest(t)

	const location = "us-west1"
	const productDisplayName = "fake_product_display_name_for_testing"
	const productCategory = "homegoods"
	const productID = "fake_product_id_for_testing"
	const referenceImageID = "fake_reference_image_id_for_testing"
	const gcsURI = "gs://cloud-samples-data/vision/product_search/shoes_1.jpg"

	var buf bytes.Buffer

	// Create a fake product.
	if err := createProduct(&buf, tc.ProjectID, location, productID, productDisplayName, productCategory); err != nil {
		t.Fatalf("createProduct: %v", err)
	}

	// Make sure the reference image to be created does not already exist.
	if err := listReferenceImages(&buf, tc.ProjectID, location, productID); err != nil {
		t.Fatalf("listReferenceImages: %v", err)
	}
	if got := buf.String(); strings.Contains(got, referenceImageID) {
		t.Errorf("Reference image ID %s already exists", referenceImageID)
	}

	// Create reference image.
	if err := createReferenceImage(&buf, tc.ProjectID, location, productID, referenceImageID, gcsURI); err != nil {
		t.Fatalf("createReferenceImage: %v", err)
	}

	// Check if the reference image exists now.
	buf.Reset()
	if err := listReferenceImages(&buf, tc.ProjectID, location, productID); err != nil {
		t.Fatalf("listReferenceImages: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, referenceImageID) {
		t.Errorf("Reference image ID %s does not exist", referenceImageID)
	}

	// Clean up.
	deleteProduct(&buf, tc.ProjectID, location, productID)
}
