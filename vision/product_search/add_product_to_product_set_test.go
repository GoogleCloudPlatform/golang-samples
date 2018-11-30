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

func TestAddProductToProductSet(t *testing.T) {
	tc := testutil.SystemTest(t)

	const location = "us-west1"
	const productSetDisplayName = "fake_product_set_display_name_for_testing"
	const productSetID = "fake_product_set_id_for_testing"
	const productDisplayName = "fake_product_display_name_for_testing"
	const productCategory = "homegoods"
	const productID = "fake_product_id_for_testing"

	var buf bytes.Buffer

	// Create fake product set and product.
	if err := createProductSet(&buf, tc.ProjectID, location, productSetID, productSetDisplayName); err != nil {
		t.Fatalf("createProductSet: %v", err)
	}
	if err := createProduct(&buf, tc.ProjectID, location, productID, productDisplayName, productCategory); err != nil {
		t.Fatalf("createProduct: %v", err)
	}

	// Make sure the product is not in the product set.
	buf.Reset()
	if err := listProductsInProductSet(&buf, tc.ProjectID, location, productSetID); err != nil {
		t.Fatalf("listProductsInProductSet: %v", err)
	}
	if got := buf.String(); strings.Contains(got, productID) {
		t.Errorf("Product ID %s already in product set", productID)
	}

	// Add product to product set.
	if err := addProductToProductSet(&buf, tc.ProjectID, location, productID, productSetID); err != nil {
		t.Fatalf("addProductToProductSet: %v", err)
	}

	// Check if the product is in the product set now.
	buf.Reset()
	if err := listProductsInProductSet(&buf, tc.ProjectID, location, productSetID); err != nil {
		t.Fatalf("listProductsInProductSet: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, productID) {
		t.Errorf("Product ID %s is not in product set", productID)
	}

	// Clean up.
	deleteProduct(&buf, tc.ProjectID, location, productID)
	deleteProductSet(&buf, tc.ProjectID, location, productSetID)
}
