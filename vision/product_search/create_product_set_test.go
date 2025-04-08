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

func TestCreateProductSet(t *testing.T) {
	tc := testutil.SystemTest(t)

	const location = "us-west1"
	const productSetDisplayName = "fake_product_set_display_name_for_testing"
	const productSetID = "fake_product_set_id_for_testing"

	var buf bytes.Buffer

	// Ensure re-used resource names don't exist prior to test start.
	if err := getProductSet(&buf, tc.ProjectID, location, productSetID); err == nil {
		deleteProductSet(&buf, tc.ProjectID, location, productSetID)
	}

	// Create a fake product set.
	if err := createProductSet(&buf, tc.ProjectID, location, productSetID, productSetDisplayName); err != nil {
		t.Fatalf("createProductSet: %v", err)
	}

	// Check if the product set exists now.
	buf.Reset()
	if err := listProductSets(&buf, tc.ProjectID, location); err != nil {
		t.Fatalf("listProductSets: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, productSetID) {
		t.Errorf("Product set ID %s does not exist", productSetID)
	}

	// Clean up.
	deleteProductSet(&buf, tc.ProjectID, location, productSetID)
}
