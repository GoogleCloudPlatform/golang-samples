// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestCreateProducts(t *testing.T) {
	var projectId = os.Getenv("GOLANG_SAMPLES_PROJECT_ID")
	const location = "us-west1"
	const productDisplayName = "fake_product_display_name_for_testing"
	const productCategory = "homegoods"
	const productId = "fake_product_id_for_testing"

	var buf bytes.Buffer

	// Make sure the product to be created does not already exist.
	err := listProducts(&buf, projectId, location)
	if err != nil {
		t.Fatal(err)
	}
	if got := buf.String(); strings.Contains(got, productId) {
		t.Errorf("Product ID %s already exists", productId)
	}

	// Create a fake product.
	err = createProduct(&buf, projectId, location, productId, productDisplayName, productCategory)
	if err != nil {
		t.Fatal(err)
	}

	// Check if the product exists now.
	buf.Reset()
	err = listProducts(&buf, projectId, location)
	if err != nil {
		t.Fatal(err)
	}
	if got := buf.String(); !strings.Contains(got, productId) {
		t.Errorf("Product ID %s does not exist", productId)
	}

	// Clean up.
	deleteProduct(&buf, projectId, location, productId)
}
