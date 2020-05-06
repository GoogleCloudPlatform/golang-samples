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

package tips

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestScopeDemo(t *testing.T) {
	req := httptest.NewRequest("GET", "/", strings.NewReader(""))
	rr := httptest.NewRecorder()
	ScopeDemo(rr, req)

	want := `Global: "slow", Local: "fast"`
	if got := rr.Body.String(); got != want {
		t.Errorf("ScopeDemo got %q, want %q", got, want)
	}
}
