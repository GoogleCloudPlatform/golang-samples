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

/*The sampletests command adds the region tags tested by each test to the XML
properties of that test case. It reads JUnit XML from stdin and writes JUnit
XML to stdout.

For example, if TestFoo tests the regions foo_hello_world and
foo_hello_gopher, the TestFoo element will have the following property:
    <property name="region_tags" value="foo_hello_world,foo_hello_gopher"></property>

sampletests only looks at direct function calls by tests, not transitive calls.

There are some duplicate region tags, but they aren't tracked anywhere else,
so it's OK if they are "applied" to more than one test.

sampletests only looks in the current module, which matches the behavior of
`go test`. So, if you run `go test ./...` and sampletests in the same directory
they should both find the same set of packages.

The test coverage over all regions is printed to stderr at the end. The coverage
is based on the entire module, not just the tests that happen to be in the
given XML input. The XML may not be for for all tests in the module.

Warnings are printed to stderr for invalid region tags (e.g. mis-matched START
and END tags).

The -enable_xml flag can be used to disable XML processing and only print
warnings and coverage.

To get the number of unique region tags in the repo manually, run:
	grep -ERho '\[START .+' | sort -u | wc -l
*/
package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"go/ast"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/jstemmer/go-junit-report/formatter"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/types/typeutil"
)

var (
	startRe = regexp.MustCompile("\\[START ([[:word:]]+)\\]")
	endRe   = regexp.MustCompile("\\[END ([[:word:]]+)\\]")
)

// testRange includes the name of the test and the line number ranges it tests.
// testRange does not include the name of the file being tested - file names are
// tracked separately.
type testRange struct {
	pkgPath    string
	testName   string
	start, end int
}

// regionTag represents a single region tag. There may be many regions in a
// single file, regions with duplicate names in the same file or separate files,
// or any other combination.
type regionTag struct {
	filePath   string
	name       string
	start, end int
}

func main() {
	enableXML := flag.Bool("enable_xml", true, "Enable XML processing")
	flag.Parse()

	// Get a set of all region tags and a map from test names to region tags.
	uniqueRegionTags, testRegionTags, err := testsToRegionTags(".")
	if err != nil {
		log.Fatal(err)
	}

	testedRegionTags := map[string]struct{}{}
	for _, regionsTested := range testRegionTags {
		for _, regions := range regionsTested {
			for r := range regions {
				testedRegionTags[r] = struct{}{}
			}
		}
	}
	fmt.Fprintf(os.Stderr, "%d/%d (%.2f%%) of region tags are tested.\n", len(testedRegionTags), len(uniqueRegionTags), 100*float64(len(testedRegionTags))/float64(len(uniqueRegionTags)))

	if *enableXML {
		if err := processXML(os.Stdin, os.Stdout, testRegionTags); err != nil {
			log.Fatal(err)
		}
	}
}

// testsToRegionTags gets region tag info.
//
// allRegions is a set of all regions in the current module.
// testedRegions is a map from package path -> test name -> set of regions.
func testsToRegionTags(dir string) (allRegions map[string]struct{}, testedRegions map[string]map[string]map[string]struct{}, err error) {
	testFileRanges, err := testCoverage(dir)
	if err != nil {
		return nil, nil, err
	}

	regionFileRanges, err := regionTags(dir)
	if err != nil {
		return nil, nil, err
	}

	allRegions = map[string]struct{}{}
	for _, ranges := range regionFileRanges {
		for regionName := range ranges {
			allRegions[regionName] = struct{}{}
		}
	}

	// testedRegions is a map from package path -> test name -> set of
	// region tags.
	testedRegions = map[string]map[string]map[string]struct{}{}

	for file, testRanges := range testFileRanges {
		for _, tr := range testRanges {
			if testedRegions[tr.pkgPath] == nil {
				testedRegions[tr.pkgPath] = map[string]map[string]struct{}{}
			}
			if testedRegions[tr.pkgPath][tr.testName] == nil {
				testedRegions[tr.pkgPath][tr.testName] = map[string]struct{}{}
			}

			for regionName, regions := range regionFileRanges[file] {
				for _, region := range regions {
					// If the test range start falls within the region, it's
					// tested.
					switch {
					case
						tr.start >= region.start && tr.start <= region.end,
						tr.end >= region.start && tr.end <= region.end,
						region.start >= tr.start && region.start <= tr.end,
						region.end >= tr.start && region.end <= tr.end:

						testedRegions[tr.pkgPath][tr.testName][regionName] = struct{}{}
					}
				}
			}
		}
	}

	return allRegions, testedRegions, nil
}

// testCoverage returns a map from file path to a slice of testRanges, which
// include the name of the test and the line numbers of the file path tested.
//
// testCoverage only looks at direct function calls from the given test, it
// does not look at transitive calls. This may lead to missing some region tags.
// See https://github.com/GoogleCloudPlatform/golang-samples/issues/1402.
func testCoverage(dir string) (map[string][]testRange, error) {
	result := map[string][]testRange{}

	config := &packages.Config{
		Mode:  packages.NeedName | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo,
		Tests: true,
	}

	pkgs, err := packages.Load(config, dir+"/...")
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
			// Use path[1] because we don't want the token, we want the actual
			// function call.
			ast.Inspect(path[1], func(node ast.Node) bool {
				call, ok := node.(*ast.CallExpr)
				if !ok {
					return true
				}
				callee := typeutil.StaticCallee(pkg.TypesInfo, call)
				if callee == nil {
					return true
				}
				if callee.Pkg() != obj.Pkg() && callee.Pkg().Name()+"_test" != obj.Pkg().Name() {
					return true
				}
				calleeScope := callee.Scope()
				calleePos := pkg.Fset.Position(calleeScope.Pos())
				calleeEnd := pkg.Fset.Position(calleeScope.End())
				result[calleePos.Filename] = append(result[calleePos.Filename], testRange{
					pkgPath:  callee.Pkg().Path(),
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

// regionTags returns a map from file name -> region name -> []*regionTag.
// Regions may be duplicated within a single file. So, we
// can't assume there is only one *regionTag with a given name.
//
// regionTags logs warnings for malformed region tags.
func regionTags(dir string) (map[string]map[string][]*regionTag, error) {
	result := map[string]map[string][]*regionTag{}
	// Iterate through every file.
	// We use filepath.Walk instead of go/packages because we want to find _all_
	// regions, not just those in .go files.
	//
	// Note: implicitly, if a file isn't tested, neither are its regions. But,
	// we walk through every file so we can get all region tags to compute
	// region tag coverage.
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			// If this is the argument directory, always check it.
			if info.Name() == dir {
				return nil
			}
			// Skip the .git directory.
			if info.Name() == ".git" {
				return filepath.SkipDir
			}
			// If the directory contains a go.mod file, skip it since those
			// files aren't in the current module.
			if _, err := os.Stat(filepath.Join(path, "go.mod")); err == nil {
				return filepath.SkipDir
			}
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

		result[path] = map[string][]*regionTag{}

		lines := strings.Split(string(src), "\n")
		for lineNum, line := range lines {
			startRegion := startRe.FindStringSubmatch(line)
			if len(startRegion) > 0 {
				name := startRegion[1]
				region := &regionTag{
					filePath: path,
					name:     name,
					start:    lineNum + 1,
					// Don't know the end yet.
				}
				result[path][name] = append(result[path][name], region)
			}

			endRegion := endRe.FindStringSubmatch(line)
			if len(endRegion) > 0 {
				name := endRegion[1]
				// This end corresponds to the most recent START with this name.
				// The same region name may show in the same file more than
				// once.
				i := len(result[path][name]) - 1
				if i < 0 {
					fmt.Fprintf(os.Stderr, "WARNING: found region tag END without START: %v:%v %v\n", path, lineNum, name)
				} else {
					result[path][name][i].end = lineNum + 1
				}
			}
		}
		for _, regions := range result[path] {
			for _, region := range regions {
				if region.end == 0 {
					fmt.Fprintf(os.Stderr, "WARNING: found region tag START without END: %v:%v %v\n", region.filePath, region.start, region.name)
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func processXML(in io.Reader, out io.Writer, testRegionTags map[string]map[string]map[string]struct{}) error {
	suites := &formatter.JUnitTestSuites{}
	dec := xml.NewDecoder(in)
	if err := dec.Decode(suites); err != nil {
		return fmt.Errorf("Decode: %v", err)
	}
	for i := range suites.Suites {
		suite := &suites.Suites[i]
		pkgPath := suite.Name
		for j := range suite.TestCases {
			testCase := &suite.TestCases[j]
			testName := testCase.Name
			regionsTested := testRegionTags[pkgPath][testName]
			regions := []string{}
			for r := range regionsTested {
				regions = append(regions, r)
			}
			sort.Strings(regions)
			if len(regions) > 0 {
				regionsString := strings.Join(regions, ",")
				testCase.Properties = []formatter.JUnitProperty{
					{Name: "region_tags", Value: regionsString},
				}
			}
		}
	}

	enc := xml.NewEncoder(out)
	enc.Indent("", "\t")
	if err := enc.Encode(suites); err != nil {
		return fmt.Errorf("Encode: %v", err)
	}
	fmt.Fprintln(out)
	return nil
}
