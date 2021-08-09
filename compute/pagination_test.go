// Copyright 2021 Google LLC
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
	"strings"
	"testing"
)

func TestPaginationSnippets(t *testing.T) {
	testProjectId := "windows-sql-cloud"
	buf := &bytes.Buffer{}

	if err := printImagesList(buf, testProjectId); err != nil {
		t.Errorf("printImagesList got err: %v", err)
	}

	result := strings.Split(buf.String(), "\n")

	if len(result) < 4 {
		t.Errorf("printImagesList returns incorrect amount of images")
	}

	buf.Reset()

	if err := printImagesListByPage(buf, testProjectId, 3); err != nil {
		t.Errorf("printImagesListByList got err: %v", err)
	}

	expectedResult := "Page 0:"
	expectedResult2 := "Page 1:"
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("printImagesListByList got %q, want %q", got, expectedResult)
	}
	if got := buf.String(); !strings.Contains(got, expectedResult2) {
		t.Errorf("printImagesListByList got %q, want %q", got, expectedResult2)
	}
}
