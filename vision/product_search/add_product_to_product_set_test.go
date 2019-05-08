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

func TestAddProductToProductSet(t *testing.T) {
	tc := testutil.SystemTest(t)

	const location = "us-west1"
	const productSetDisplayName = "fake_product_set_display_name_for_testing"
	const productSetID = "fake_product_set_id_for_testing"
	const productDisplayName = "fake_product_display_name_for_testing"
	const productCategory = "homegoods"
	const productID = "fake_product_id_for_testing"

	var buf bytes.Buffer

	// Ensure re-used resource names don't exist prior to test start.
	if err := getProductSet(&buf, tc.ProjectID, location, productSetID); err == nil {
		deleteProductSet(&buf, tc.ProjectID, location, productSetID)
	}
	if err := getProduct(&buf, tc.ProjectID, location, productID); err == nil {
		deleteProduct(&buf, tc.ProjectID, location, productID)
	}

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
