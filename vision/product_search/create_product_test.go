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

func TestCreateProduct(t *testing.T) {
	tc := testutil.SystemTest(t)

	const location = "us-west1"
	const productDisplayName = "fake_product_display_name_for_testing"
	const productCategory = "homegoods"
	const productID = "fake_product_id_for_testing"

	var buf bytes.Buffer

	// Ensure re-used resource names don't exist prior to test start.
	if err := getProduct(&buf, tc.ProjectID, location, productID); err == nil {
		deleteProduct(&buf, tc.ProjectID, location, productID)
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
