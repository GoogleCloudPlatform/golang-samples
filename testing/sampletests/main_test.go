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

package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTestsToRegionTags(t *testing.T) {
	uniqueRegionTags, testRegionTags, err := testsToRegionTags("./fakesamples")
	if err != nil {
		t.Fatalf("testsToRegionTags got err: %v", err)
	}

	wantRegionTags := map[string]struct{}{
		"fakesamples_package_decl_not_tested": {},
		"fakesamples_tested":                  {},
		"fakesamples_not_tested":              {},
		"samemodule_unchecked":                {},
	}
	if diff := cmp.Diff(uniqueRegionTags, wantRegionTags); diff != "" {
		t.Errorf("testsToRegionTags got uniqueRegionTags (%+v), want (%+v). Diff: %v", uniqueRegionTags, wantRegionTags, diff)
	}
	wantTestRegionTags := map[string]map[string]map[string]struct{}{
		"github.com/GoogleCloudPlatform/golang-samples/testing/sampletests/fakesamples": {
			"TestHello": {
				"fakesamples_tested": {},
			},
		},
	}
	if diff := cmp.Diff(testRegionTags, wantTestRegionTags); diff != "" {
		t.Errorf("testsToRegionTags got testRegionTags (%+v), want (%+v). Diff: %v", testRegionTags, wantTestRegionTags, diff)
	}
}

func TestTestCoverage(t *testing.T) {
	got, err := testCoverage("./fakesamples")
	if err != nil {
		t.Fatalf("testCoverage got err: %v", err)
	}
	// We should only find one test file with a single range.
	if gotLen, wantLen := len(got), 1; gotLen != wantLen {
		t.Fatalf("testCoverage found %d test files, want %d", gotLen, wantLen)
	}
	for _, ranges := range got {
		if gotLen, wantLen := len(ranges), 1; gotLen != wantLen {
			t.Fatalf("testCoverage found %d test ranges, want %d", gotLen, wantLen)
		}
		gotRange := ranges[0]
		want := testRange{
			pkgPath:  "github.com/GoogleCloudPlatform/golang-samples/testing/sampletests/fakesamples",
			testName: "TestHello",
			start:    25, // If hello.go changes, this test will intentionally break.
			end:      27,
		}
		if gotRange != want {
			t.Errorf("testCoverage found incorrect range: got %+v, want %+v", gotRange, want)
		}
	}
}
