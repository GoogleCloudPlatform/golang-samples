// Copyright 2023 Google LLC
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

package snippets

import (
	"bytes"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestGenerateEmbeddings(t *testing.T) {
	tc := testutil.SystemTest(t)

	prompt := "hello, say something nice."
	projectID := tc.ProjectID
	location := "us-central1"
	publisher := "google"
	model := "textembedding-gecko"

	var buf bytes.Buffer
	if err := generateEmbeddings(&buf, prompt, projectID, location, publisher, model); err != nil {
		t.Fatal(err)
	}

	if got := buf.String(); !strings.Contains(got, "embeddings generated:") {
		t.Error("generated embeddings content not found in response")
	}

}
