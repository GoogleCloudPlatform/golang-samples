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
	"bytes"
	"os"
	"strings"
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
		"fakesamples_tested_0":                {},
		"fakesamples_tested_1":                {},
		"fakesamples_tested_2":                {},
		"fakesamples_tested_3":                {},
		"fakesamples_not_tested":              {},
		"fakesamples_indirect_test":           {},
		"samemodule_not_tested":               {},
	}
	if diff := cmp.Diff(uniqueRegionTags, wantRegionTags); diff != "" {
		t.Errorf("testsToRegionTags got uniqueRegionTags (%+v), want (%+v). Diff: %v", uniqueRegionTags, wantRegionTags, diff)
	}
	wantTestRegionTags := map[string]map[string]map[string]struct{}{
		"github.com/GoogleCloudPlatform/golang-samples/testing/sampletests/fakesamples": {
			"TestHello": {
				"fakesamples_tested_0": {},
				"fakesamples_tested_1": {},
				"fakesamples_tested_2": {},
				"fakesamples_tested_3": {},
			},
			"TestIndirectlyTested": {
				"fakesamples_indirect_test": {},
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
	for _, gotRanges := range got {
		want := []testRange{
			{
				pkgPath:  "github.com/GoogleCloudPlatform/golang-samples/testing/sampletests/fakesamples",
				testName: "TestIndirectlyTested",
				start:    51, // If hello.go changes, this test will intentionally break.
				end:      53,
			},
			{
				pkgPath:  "github.com/GoogleCloudPlatform/golang-samples/testing/sampletests/fakesamples",
				testName: "TestHello",
				start:    27, // If hello.go changes, this test will intentionally break.
				end:      35,
			},
		}
		if len(gotRanges) != len(want) {
			t.Fatalf("testCoverage found %d ranges, want %d", len(gotRanges), len(want))
		}
		// Don't rely on the order of the slice.
		for _, wantRange := range want {
			found := false
			for _, gotRange := range gotRanges {
				if wantRange == gotRange {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("testCoverage found incorrect ranges: got %v, want to contain %v", gotRanges, wantRange)
			}
		}
	}
}

func TestRegionTags(t *testing.T) {
	got, err := regionTags("./fakesamples")
	if err != nil {
		t.Fatalf("regionTags got err: %v", err)
	}
	gotRegions := map[string]regionTag{}
	for _, region := range got {
		for name, regionsWithName := range region {
			if len(regionsWithName) != 1 {
				t.Fatalf("regionTags found %d for region %q, want 1", len(regionsWithName), name)
			}
			gotRegions[name] = *regionsWithName[0]
		}
	}
	gotTestRegion := gotRegions["fakesamples_tested_1"]
	if gotStart, want := gotTestRegion.start, 24; gotStart != want {
		t.Errorf("regionTags %v got start %v, want %v", gotTestRegion.name, gotStart, want)
	}
	if gotEnd, want := gotTestRegion.end, 37; gotEnd != want {
		t.Errorf("regionTags %v got end %v, want %v", gotTestRegion.name, gotEnd, want)
	}
}

func TestProcessXML(t *testing.T) {
	_, testRegionTags, err := testsToRegionTags("./fakesamples")
	if err != nil {
		t.Fatalf("testsToRegionTags got err: %v", err)
	}

	input, err := os.Open("testdata/raw_log.xml")
	if err != nil {
		t.Fatalf("os.Open: %v", err)
	}

	buf := &bytes.Buffer{}

	processXML(input, buf, testRegionTags)

	want := `<property name="region_tags" value="fakesamples_tested_0,fakesamples_tested_1,fakesamples_tested_2,fakesamples_tested_3"></property>`
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("processXML got\n\n----\n%v\n----\nWant to contain:\n\n----\n%v", got, want)
	}
}
