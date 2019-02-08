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

func TestUpdateProductLabels(t *testing.T) {
	tc := testutil.SystemTest(t)

	const location = "us-west1"
	const productDisplayName = "fake_product_display_name_for_testing"
	const productCategory = "homegoods"
	const productID = "fake_product_id_for_testing"
	const key = "fake_key_for_testing"
	const value = "fake_value_for_testing"

	var buf bytes.Buffer

	// Create a fake product.
	if err := createProduct(&buf, tc.ProjectID, location, productID, productDisplayName, productCategory); err != nil {
		t.Fatalf("createProduct: %v", err)
	}

	// Make sure the label to be added to the product does not already exist.
	if err := getProduct(&buf, tc.ProjectID, location, productID); err != nil {
		t.Fatalf("getProduct: %v", err)
	}
	if got := buf.String(); strings.Contains(got, key) || strings.Contains(got, value) {
		t.Errorf("Key-value %s %s already exists", key, value)
	}

	// Update product labels.
	buf.Reset()
	if err := updateProductLabels(&buf, tc.ProjectID, location, productID, key, value); err != nil {
		t.Fatalf("updateProductLabels: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, key) || !strings.Contains(got, value) {
		t.Errorf("Label %s %s does not exist", key, value)
	}

	// Clean up.
	deleteProduct(&buf, tc.ProjectID, location, productID)
}
