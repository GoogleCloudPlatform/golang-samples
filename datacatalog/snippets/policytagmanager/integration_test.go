// Copyright 2020 Google LLC
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

// Package policytagmanager demonstrates interactions with the Policy
// Tag Manager client, used for managing Policy Tags.  This functionality
// underpins features in other services such as BigQuery's ability to
// define column-level access control.
package policytagmanager

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestPolicyTagManager(t *testing.T) {
	t.Skip("deprecated service")
	tc := testutil.SystemTest(t)

	location := "us"
	// If you wish to capture output, change the output to os.Stdout.
	// Normal operation should use an instance of ioutil.Discard.
	output := ioutil.Discard

	taxonomyName := fmt.Sprintf("example-taxonomy-%d", time.Now().UnixNano())
	taxID, err := createTaxonomy(output, tc.ProjectID, location, taxonomyName)
	if err != nil {
		t.Errorf("createTaxonomy: %v", err)
	}
	defer deleteTaxonomy(taxID)

	if err := getTaxonomy(output, taxID); err != nil {
		t.Errorf("getTaxonomy: %v", err)
	}

	if err := listTaxonomies(output, tc.ProjectID, location); err != nil {
		t.Errorf("listTaxonomies: %v", err)
	}

	// Create some policy tags
	displayName := "PII Tag"
	tagOne, err := createPolicyTag(output, taxID, displayName, "")
	if err != nil {
		t.Errorf("createPolicyTag(%s): %v", displayName, err)
	}

	displayName = "Child PII Tag"
	tagTwo, err := createPolicyTag(output, taxID, displayName, tagOne)
	if err != nil {
		t.Errorf("createPolicyTag(%s): %v", displayName, err)
	}

	if err := getPolicyTag(output, tagOne); err != nil {
		t.Errorf("getPolicyTag(%s): %v", tagOne, err)
	}

	probedPermissions := []string{
		"datacatalog.categories.fineGrainedGet",
	}
	// check before setting policy
	var buf bytes.Buffer
	if err := testIAMPermissions(&buf, tagOne, probedPermissions); err != nil {
		t.Errorf("testIAMPermissions(%s): %v", tagOne, err)
	}
	wantedResp := "of the 1 permissions probed, caller has 0 permissions"
	if !strings.Contains(buf.String(), wantedResp) {
		t.Errorf("unexpected output (%q) did not contain (%s)", buf.String(), wantedResp)
	}
	buf.Reset()

	if err := setIAMPolicy(output, tagOne, "allAuthenticatedUsers"); err != nil {
		t.Errorf("setIAMPolicy(%s): %v", tagOne, err)
	}

	if err := testIAMPermissions(&buf, tagOne, probedPermissions); err != nil {
		t.Errorf("testIAMPermissions(%s): %v", tagOne, err)
	}
	wantedResp = "of the 1 permissions probed, caller has 1 permissions: datacatalog.categories.fineGrainedGet"
	if !strings.Contains(buf.String(), wantedResp) {
		t.Errorf("unexpected output (%q) did not contain (%s)", buf.String(), wantedResp)
	}
	buf.Reset()

	if err := getIAMPolicy(output, tagOne); err != nil {
		t.Errorf("getIAMPolicy(%s): %v", tagOne, err)
	}

	if err := listPolicyTags(output, taxID); err != nil {
		t.Errorf("listPolicyTags: %v", err)
	}

	// delete a Policy tag
	if err := deletePolicyTag(tagTwo); err != nil {
		t.Errorf("deletePolicyTag(%s): %v", tagTwo, err)
	}

	if err := deleteTaxonomy(taxID); err != nil {
		t.Errorf("deleteTaxonomy: %v", err)
	}

}
