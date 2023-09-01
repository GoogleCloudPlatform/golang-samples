// Copyright 2023 Google LLC
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

package language_v2

import (
	"bytes"
	"strings"
	"testing"
)

func TestAnalyzeEntities(t *testing.T) {
	buf := new(bytes.Buffer)

	text := "Google is located in Mountain View."
	err := analyzeEntities(buf, text)
	if err != nil {
		t.Fatalf("TestAnalyzeEntities: %v", err)
	}

	got := buf.String()
	if want := "entities:"; !strings.Contains(got, want) {
		t.Fatalf("got %q, want %q", got, want)
	}
	if want := "language_code:"; !strings.Contains(got, want) {
		t.Fatalf("got %q, want %q", got, want)
	}
	if want := "language_supported:"; !strings.Contains(got, want) {
		t.Fatalf("got %q, want %q", got, want)
	}
}
