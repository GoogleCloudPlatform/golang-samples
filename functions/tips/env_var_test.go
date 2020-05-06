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
	"fmt"
	"net/http/httptest"
	"os"
	"testing"
)

func TestEnvVar(t *testing.T) {
	tests := []struct {
		foo string
	}{
		{"bar"},
		{},
	}
	for _, test := range tests {
		os.Setenv("FOO", test.foo)
		rr := httptest.NewRecorder()
		EnvVar(rr, nil)

		want := fmt.Sprintf("FOO: %q", test.foo)
		if got := rr.Body.String(); got != want {
			t.Errorf("EnvVar(%s) got %q, want %q", test.foo, got, want)
		}
	}
}
