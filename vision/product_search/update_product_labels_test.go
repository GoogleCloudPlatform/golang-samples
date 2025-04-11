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

func TestUpdateProductLabels(t *testing.T) {
	tc := testutil.SystemTest(t)

	const location = "us-west1"
	const productDisplayName = "fake_product_display_name_for_testing"
	const productCategory = "homegoods"
	const productID = "fake_product_id_for_testing"
	const key = "fake_key_for_testing"
	const value = "fake_value_for_testing"

	var buf bytes.Buffer

	// Ensure re-used resource names don't exist prior to test start.
	if err := getProduct(&buf, tc.ProjectID, location, productID); err == nil {
		deleteProduct(&buf, tc.ProjectID, location, productID)
	}
	buf.Reset()

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
