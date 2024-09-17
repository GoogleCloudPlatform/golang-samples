// Copyright 2024 Google LLC
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

package discoveryengine

import (
	"os"
	"testing"
)

func TestSearch(t *testing.T) {
	// Note: This test assumes a pre-populated Vertex AI Search data store and app in the
	// below-referenced project.
	projectID := os.Getenv("GOLANG_SAMPLES_PROJECT_ID")
	location := "global"
	dataStoreID := "test-data-store"
	appID := "test-app"

	createDataStore(projectID, location, dataStoreID)
	search(projectID, location, appID, "test")
}
