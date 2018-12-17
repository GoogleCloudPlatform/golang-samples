// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

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
