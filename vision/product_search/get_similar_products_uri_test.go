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

func TestGetSimilarProductsURI(t *testing.T) {
	// TODO(#2924): Re-enable once resource availability is hammered out.
	t.Skip()
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
