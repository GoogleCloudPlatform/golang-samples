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

// Smoke tests for main.go

package main

import (
	"os"
	"testing"
	"time"
)

func setup(t *testing.T) string {
	orgId := os.Getenv("GCLOUD_ORGANIZATION")
	if orgId == "" {
		t.Skip("GCLOUD_ORGANIZATION not set")
	}
	return orgId
}

func TestAllAssetsReturned(t *testing.T) {
	org := setup(t)
	found, err := ListAllAssets(org)
	if err != nil {
		t.Errorf("foo failed: %v", err)
	}
	if found < 59 {
		t.Errorf("Not enough results found: %d", found)
	}
}

func TestAllProjectsReturned(t *testing.T) {
	org := setup(t)
	found, err := ListAllProjectAssets(org)
	if err != nil {
		t.Errorf("foo failed: %v", err)
	}
	if found > 4 || found == 0 {
		t.Errorf("Unexpected number of results: %d", found)
	}
}

func TestBeforeDateNoAssetsReturned(t *testing.T) {
	org := setup(t)
	var nothingInstant = time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)

	found, err := ListAllProjectAssetsAtTime(org, nothingInstant)
	if err != nil {
		t.Errorf("foo failed: %v", err)
	}
	if found > 0 {
		t.Errorf("Results found for at %v: %d", nothingInstant, found)
	}

	var somethingInstant = time.Date(2019, 3, 15, 0, 0, 0, 0, time.UTC)
	found, err = ListAllProjectAssetsAtTime(org, somethingInstant)
	if err != nil {
		t.Errorf("foo failed: %v", err)
	}
	if found < 3 {
		t.Errorf("Not enough results at %v: %d", somethingInstant, found)
	}
}

func TestProjectsWithStatusReturned(t *testing.T) {
	org := setup(t)
	found, err := ListAllProjectAssetsWithStateChanges(org)
	if err != nil {
		t.Errorf("foo failed: %v", err)
	}
	if found > 4 || found == 0 {
		t.Errorf("Unexpected number of results: %d", found)
	}
}

/*
func TestListAssetsNoFilterOrDate(t *testing.T) {
	la := setup(t)
	assertTrue(59 >= snippets.listAssets(null, null).size())
}

func TestListAssetsWithFilterAndInstance(t *testing.T) {
	la := setup(t)
	assertTrue(
		3 >= snippets.listAssets(PROJECT_ASSET_FILTERS, SOMETHING_INSTANCE).size())
}

func TestChangesReturnsValues(t *testing.T) {
	ImmutableList < ListAssetsResult > result =
		snippets.listAssetAndStatusChanges(
			Duration.ofDays(3), AssetSnippets.PROJECT_ASSET_FILTERS, SOMETHING_INSTANCE)
	assertTrue("Result: "+result.toString(), result.toString().contains("ADDED"))
	assertTrue(3 >= result.size())
}
*/
