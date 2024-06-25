// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package search

import "testing"

func TestSearch(t *testing.T) {
	t.Skip("See http://github.com/GoogleCloudPlatform/golang-samples/issues/3569")
	// Customize this for your project
	projectID := "my-project-id"
	location := "us-central1"
	searchEngineID := "my-search-engine-id"
	query := "my-query"

	err := search(projectID, location, searchEngineID, query)

	if err != nil {
		t.Errorf("search() error = %v", err)
	}
}
