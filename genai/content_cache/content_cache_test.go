// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package content_cache

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestContentCaching(t *testing.T) {
	tc := testutil.SystemTest(t)

	t.Setenv("GOOGLE_GENAI_USE_VERTEXAI", "1")
	t.Setenv("GOOGLE_CLOUD_LOCATION", "us-central1")
	t.Setenv("GOOGLE_CLOUD_PROJECT", tc.ProjectID)

	buf := new(bytes.Buffer)

	// 1) Create a content cache. The name of the cache created will be used in the next test steps.
	cacheName, err := createContentCache(buf)
	if err != nil {
		t.Fatalf("createContentCache: %v", err.Error())
	}

	fmt.Println(cacheName)

	// 2) Update the cached content by its name.
	buf.Reset()
	err = updateContentCache(buf, cacheName)
	if err != nil {
		t.Errorf("updateContentCache: %v", err.Error())
	}

	// 3) - Use
	// 4) - Delete
}
