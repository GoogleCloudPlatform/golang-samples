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
// Tag Manager Client, used for managing Policy Tags.  This construct
// underpins features in other services such as BigQuery's ability to
// define column-level access control.
package policytagmanager

import (
	"os"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestPolicyTagManager(t *testing.T) {
	tc := testutil.SystemTest(t)

	location := "us"
	// If you wish to capture output, change the output to os.Stdout.
	// Normal operation should use an instance of ioutil.Discard.
	output := os.Stdout

	taxID, err := createTaxonomy(tc.ProjectID, location, output)
	if err != nil {
		t.Errorf("createTaxonomy: %v", err)
	}

	if err := listTaxonomies(tc.ProjectID, location, output); err != nil {
		t.Errorf("listTaxonomies: %v", err)
	}

	// Create some policy tags
	displayName := "PII Tag"
	tagOne, err := createPolicyTag(taxID, displayName, "", output)
	if err != nil {
		t.Errorf("createPolicyTag(%s): %v", displayName, err)
	}

	displayName = "Child PII Tag"
	tagTwo, err := createPolicyTag(taxID, displayName, tagOne, output)
	if err != nil {
		t.Errorf("createPolicyTag(%s): %v", displayName, err)
	}

	if err := listPolicyTags(taxID, output); err != nil {
		t.Errorf("listPolicyTags: %v", err)
	}

	// delete a Policy tag
	if err := deletePolicyTag(tagTwo, output); err != nil {
		t.Errorf("deletePolicy(%s): %v", tagTwo, err)
	}

	if err := deleteTaxonomy(taxID, output); err != nil {
		t.Errorf("deleteTaxonomy: %v", err)
	}

}
