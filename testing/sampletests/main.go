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

// The sampletests command adds the region tags tested by each test to the test
// name. It reads from stdin and writes to stdout.
//
// For example, if TestFoo tests the regions foo_hello_world and
// foo_hello_gopher, the result will be:
//     TestFoo -> TestFoo__foo_hello_world__foo_hello_gopher
//
// It only looks at direct function calls by tests, not transitive calls.
//
// There are some duplicate region tags, but they aren't tracked anywhere else,
// so it's OK if they are "applied" to more than one test.
//
// There are some duplicate test names, which means it's possible to get names
// like:
//     TestCreate__spanner_create__bigtable_create
// Again, this is an acceptable tradeoff. It only happens when both regions are
// tested. If only one is tested, only one would show up. Failures are
// attributed to both snippets, though.
//
// The test coverage over all regions is printed at the end.
//
// To get the number of unique region tags in the repo manually, run:
//     grep -RoPh '\[START \K(.+)\]' | sort -u | wc -l
package main

import (
	"fmt"
	"go/ast"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/types/typeutil"
)

var (
	startRe = regexp.MustCompile("\\[START ([^\\]]+)\\]")
	endRe   = regexp.MustCompile("\\[END ([^\\]]+)\\]")
)

// testRange includes the name of the test and the line number ranges it tests.
// testRange does not include the name of the file being tested - file names are
// tracked separately.
type testRange struct {
	testName   string
	start, end int
}

func main() {
	// Get a set of all region tags and a map from test names to region tags.
	uniqueRegionTags, testRegionTags, err := testsToRegionTags()
	if err != nil {
		log.Fatal(err)
	}

	// Read all input, make replacements.
	inputBytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("ioutil.ReadAll: %v", err)
	}
	input := string(inputBytes)

	testedRegionTags := map[string]struct{}{}
	newTestNames := map[string]string{}
	for test, regionsTested := range testRegionTags {
		regions := []string{}
		for r := range regionsTested {
			regions = append(regions, r)
			testedRegionTags[r] = struct{}{}
		}
		newName := test
		if len(regions) > 0 {
			newName = fmt.Sprintf("%v__%v", test, strings.Join(regions, "__"))
		}
		newTestNames[test] = newName
	}

	for test, newName := range newTestNames {
		// Include whitespace to avoid replacing prefixes (TestFoo of
		// TestFooBar).
		input = strings.Replace(input, " "+test+" ", " "+newName+" ", -1)
		input = strings.Replace(input, " "+test+"\n", " "+newName+"\n", -1)
	}

	fmt.Println(input)
	fmt.Printf("%d/%d (%.2f%%) of region tags are tested.\n", len(testedRegionTags), len(uniqueRegionTags), 100*float64(len(testedRegionTags))/float64(len(uniqueRegionTags)))
}

// testsToRegionTags returns the total unique region tags in the current
// directory, a map from test to sets of regions, and an error.
func testsToRegionTags() (unique map[string]struct{}, testRegionTags map[string]map[string]struct{}, err error) {
	// Get map from file to []testRange.
	testFileRanges, err := testCoverage()
	if err != nil {
		return nil, nil, err
	}

	// testRegionTags is a map from test to sets of regions.
	// Uses a set of regions instead of a slice to avoid duplication.
	testRegionTags = map[string]map[string]struct{}{}

	// Initialize the map values of testRegionTags.
	for _, t := range testFileRanges {
		for _, r := range t {
			testRegionTags[r.testName] = map[string]struct{}{}
		}
	}

	uniqueRegionTags := map[string]struct{}{}

	// Iterate through every file.
	// Note: implicitly, if a file isn't tested, neither are its regions. But,
	// we walk through every file so we can get all region tags to compute
	// region tag coverage.
	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		path, err = filepath.Abs(path)
		if err != nil {
			return err
		}

		src, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		lines := strings.Split(string(src), "\n")
		for lineNum, line := range lines {
			if regions := startRe.FindAllStringSubmatch(line, -1); len(regions) > 0 {
				region := regions[0][1]
				uniqueRegionTags[region] = struct{}{}
				for _, r := range testFileRanges[path] {
					if r.start >= lineNum && r.end >= lineNum {
						testRegionTags[r.testName][region] = struct{}{}
					}
				}
			}
			if regions := endRe.FindAllStringSubmatch(line, -1); len(regions) > 0 {
				region := regions[0][1]
				for _, r := range testFileRanges[path] {
					if r.start >= lineNum && r.end >= lineNum {
						testRegionTags[r.testName][region] = struct{}{}
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return uniqueRegionTags, testRegionTags, nil
}

// testCoverage returns a map from file path to a slice of testRanges, which
// include the name of the test and the line numbers of the file path tested.
//
// testCoverage only looks at direct function calls from the given test, it
// does not look at transitive calls. This may lead to missing some region tags.
func testCoverage() (map[string][]testRange, error) {
	result := map[string][]testRange{}

	config := &packages.Config{
		Mode:  packages.NeedName | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo,
		Tests: true,
	}

	pkgs, err := packages.Load(config, "./...")
	if err != nil {
		return nil, fmt.Errorf("packages.Load: %v", err)
	}
	for _, pkg := range pkgs {
		for id, obj := range pkg.TypesInfo.Defs {
			if !strings.HasPrefix(id.Name, "Test") {
				continue
			}
			var file *ast.File
			for _, f := range pkg.Syntax {
				if f.Pos() <= id.Pos() && id.Pos() <= f.End() {
					file = f
				}
			}
			if file == nil {
				return nil, fmt.Errorf("file for %q not found", id.Name)
			}
			path, exact := astutil.PathEnclosingInterval(file, id.Pos(), id.End())
			if !exact {
				return nil, fmt.Errorf("PathEnclosingInterval got not exact path for %q in %v", id.Name, file)
			}
			ast.Inspect(path[1], func(node ast.Node) bool {
				call, ok := node.(*ast.CallExpr)
				if !ok {
					return true
				}
				callee := typeutil.StaticCallee(pkg.TypesInfo, call)
				if callee == nil {
					return true
				}
				if callee.Pkg() != obj.Pkg() {
					return true
				}
				calleeScope := callee.Scope()
				calleePos := pkg.Fset.Position(calleeScope.Pos())
				calleeEnd := pkg.Fset.Position(calleeScope.End())
				result[calleePos.Filename] = append(result[calleePos.Filename], testRange{
					testName: id.Name,
					start:    calleePos.Line,
					end:      calleeEnd.Line,
				})
				return true
			})
		}
	}

	return result, nil
}
