// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package productsearch

import (
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestMain(m *testing.M) {
	tc, ok := testutil.ContextMain(m)
	if !ok {
		log.Printf("Could not get test context")
		return
	}

	// Create a fake product to test if the service account has access.
	// The Product Search API does not support using a service account in a
	// different project (the setup of the golang-samples tests).
	// This is not perfect, since it will hide failures in the createProduct
	// sample and doesn't run tests on all projects.
	const location = "us-west1"
	const productSetDisplayName = "fake_product_set_display_name_for_testing"
	const productSetID = "fake_product_set_id_for_testing"
	const productDisplayName = "fake_product_display_name_for_testing"
	const productCategory = "homegoods"
	const productID = "fake_product_id_for_testing"

	var buf bytes.Buffer

	if err := createProduct(&buf, tc.ProjectID, location, productID, productDisplayName, productCategory); err != nil {
		log.Printf("Skipping product_search tests: Could not create product: %v", err)
		return
	}

	deleteProduct(&buf, tc.ProjectID, location, productID)

	// We successfuly created a project. Run the tests.
	os.Exit(m.Run())
}
