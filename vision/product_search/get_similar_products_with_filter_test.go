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

func TestGetSimilarProductsWithFilter(t *testing.T) {
	tc := testutil.SystemTest(t)

	const location = "us-west1"
	const productSetID = "indexed_product_set_id_for_testing"
	const productCategory = "apparel"
	const productID1 = "indexed_product_id_for_testing_1"
	const productID2 = "indexed_product_id_for_testing_2"
	const file = "../testdata/shoes_1.jpg"
	const filter = "style=womens"

	var buf bytes.Buffer

	if err := getSimilarProducts(&buf, tc.ProjectID, location, productSetID, productCategory, file, filter); err != nil {
		t.Fatalf("getSimilarProducts: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, productID1) || strings.Contains(got, productID2) {
		t.Errorf("Product ID %s should be the only one returned", productID1)
	}
}
