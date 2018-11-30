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

func TestCreateProduct(t *testing.T) {
	tc := testutil.SystemTest(t)

	const location = "us-west1"
	const productDisplayName = "fake_product_display_name_for_testing"
	const productCategory = "homegoods"
	const productID = "fake_product_id_for_testing"

	var buf bytes.Buffer

	// Make sure the product to be created does not already exist.
	if err := listProducts(&buf, tc.ProjectID, location); err != nil {
		t.Fatalf("listProducts: %v", err)
	}
	if got := buf.String(); strings.Contains(got, productID) {
		t.Errorf("Product ID %s already exists", productID)
	}

	// Create a fake product.
	if err := createProduct(&buf, tc.ProjectID, location, productID, productDisplayName, productCategory); err != nil {
		t.Fatalf("createProduct: %v", err)
	}

	// Check if the product exists now.
	buf.Reset()
	if err := listProducts(&buf, tc.ProjectID, location); err != nil {
		t.Fatalf("listProducts: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, productID) {
		t.Errorf("Product ID %s does not exist", productID)
	}

	// Clean up.
	deleteProduct(&buf, tc.ProjectID, location, productID)
}
