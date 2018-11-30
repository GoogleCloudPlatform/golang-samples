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

func TestDeleteProductSet(t *testing.T) {
	tc := testutil.SystemTest(t)

	const location = "us-west1"
	const productSetDisplayName = "fake_product_set_display_name_for_testing"
	const productSetID = "fake_product_set_id_for_testing"

	var buf bytes.Buffer

	// Create a fake product set.
	if err := createProductSet(&buf, tc.ProjectID, location, productSetID, productSetDisplayName); err != nil {
		t.Fatalf("createProductSet: %v", err)
	}

	// Confirm the product set exists.
	if err := listProductSets(&buf, tc.ProjectID, location); err != nil {
		t.Fatalf("listProductSets: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, productSetID) {
		t.Errorf("Product set ID %s does not exist", productSetID)
	}

	// Delete the product set.
	deleteProductSet(&buf, tc.ProjectID, location, productSetID)

	// Confirm the product set has been deleted.
	buf.Reset()
	if err := listProductSets(&buf, tc.ProjectID, location); err != nil {
		t.Fatalf("listProductSets: %v", err)
	}
	if got := buf.String(); strings.Contains(got, productSetID) {
		t.Errorf("Product set ID %s still exists", productSetID)
	}
}
