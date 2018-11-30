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

func TestGetSimilarProductsURI(t *testing.T) {
	tc := testutil.SystemTest(t)

	const location = "us-west1"
	const productSetID = "indexed_product_set_id_for_testing"
	const productCategory = "apparel"
	const productID1 = "indexed_product_id_for_testing_1"
	const productID2 = "indexed_product_id_for_testing_2"
	const imageURI = "gs://cloud-samples-data/vision/product_search/shoes_1.jpg"
	const filter = ""

	var buf bytes.Buffer

	if err := getSimilarProductsURI(&buf, tc.ProjectID, location, productSetID, productCategory, imageURI, filter); err != nil {
		t.Fatalf("getSimilarProductsURI: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, productID1) || !strings.Contains(got, productID2) {
		t.Errorf("Product IDs %s %s not returned", productID1, productID2)
	}
}
