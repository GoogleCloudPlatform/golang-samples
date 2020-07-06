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

// [START fakesamples_package_decl_not_tested]

// Package hello is used only for testing. It's not in a testdata directory
// to work with go/packages.
package hello

// [END fakesamples_package_decl_not_tested]

// [START fakesamples_tested_0]
// [START fakesamples_tested_1]

// Hello returns hello. Only for testing. May change at any time.
func Hello() string {

	// [END fakesamples_tested_0]
	// [START fakesamples_tested_3]
	// [START fakesamples_tested_2]
	return "Hello!"
	// [END fakesamples_tested_2]

}

// [END fakesamples_tested_1]
// [END fakesamples_tested_3]

// [START fakesamples_not_tested]

func notTested() string {
	return "This function isn't tested!"
}

// [END fakesamples_not_tested]

// [START fakesamples_indirect_test]

// IndirectlyTested returns a string. Only for testing. May change at any time.
func IndirectlyTested() string {
	return "This function is tested via a function reference rather than a direct call"
}

// [END fakesamples_indirect_test]
