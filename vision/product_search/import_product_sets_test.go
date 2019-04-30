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

func TestImportProductSets(t *testing.T) {
	tc := testutil.SystemTest(t)

	const location = "us-west1"
	const gcsURI = "gs://cloud-samples-data/vision/product_search/product_sets.csv"
	const productSetID = "fake_product_set_id_for_testing"
	const productID1 = "fake_product_id_for_testing_1"
	const productID2 = "fake_product_id_for_testing_2"
	const imageURI1 = "shoes_1.jpg"
	const imageURI2 = "shoes_2.jpg"

	var buf bytes.Buffer

	// Ensure re-used resource names don't exist prior to test start.
	if err := getProductSet(&buf, tc.ProjectID, location, productSetID); err == nil {
		deleteProductSet(&buf, tc.ProjectID, location, productSetID)
	}
	if err := getProduct(&buf, tc.ProjectID, location, productID1); err == nil {
		deleteProduct(&buf, tc.ProjectID, location, productID1)
	}
	if err := getProduct(&buf, tc.ProjectID, location, productID2); err == nil {
		deleteProduct(&buf, tc.ProjectID, location, productID2)
	}

	buf.Reset()

	// Import product set.
	if err := importProductSets(&buf, tc.ProjectID, location, gcsURI); err != nil {
		t.Fatalf("importProductSets: %v", err)
	}

	// Check if the product set exists now.
	buf.Reset()
	if err := listProductSets(&buf, tc.ProjectID, location); err != nil {
		t.Fatalf("listProductSets: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, productSetID) {
		t.Errorf("Product set ID %s does not exist", productSetID)
	}

	// Check if the products exist.
	buf.Reset()
	if err := listProducts(&buf, tc.ProjectID, location); err != nil {
		t.Fatalf("listProducts: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, productID1) || !strings.Contains(got, productID2) {
		t.Errorf("Product IDs %s %s do not exist", productID1, productID2)
	}

	// Check if the products are in the product set.
	buf.Reset()
	if err := listProductsInProductSet(&buf, tc.ProjectID, location, productSetID); err != nil {
		t.Fatalf("listProductsInProductSet: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, productID1) || !strings.Contains(got, productID2) {
		t.Errorf("Product IDs %s %s do not exist in product set", productID1, productID2)
	}

	// check if the reference images exsit.
	buf.Reset()
	if err := listReferenceImages(&buf, tc.ProjectID, location, productID1); err != nil {
		t.Fatalf("listReferenceImages: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, imageURI1) {
		t.Errorf("Reference image uri %s does not exist in product set", imageURI1)
	}

	buf.Reset()
	if err := listReferenceImages(&buf, tc.ProjectID, location, productID2); err != nil {
		t.Fatalf("listReferenceImages: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, imageURI2) {
		t.Errorf("Reference image uri %s does not exist in product set", imageURI2)
	}

	// Clean up.
	deleteProduct(&buf, tc.ProjectID, location, productID1)
	deleteProduct(&buf, tc.ProjectID, location, productID2)
	deleteProductSet(&buf, tc.ProjectID, location, productSetID)
}
