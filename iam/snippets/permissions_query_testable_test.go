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

package snippets

import (
	"bytes"
	"os"
	"testing"
)

func TestQueryTestablePermissions(t *testing.T) {
	buf := &bytes.Buffer{}
	project := os.Getenv("GOLANG_SAMPLES_PROJECT_ID")
	name := "//cloudresourcemanager.googleapis.com/projects/" + project
	permissions, err := queryTestablePermissions(buf, name)
	if err != nil {
		t.Fatalf("queryTestablePermissions: %v", err)
	}
	if len(permissions) < 1 {
		t.Fatalf("queryTestablePermissions: expected at least 1 item")
	}
}
