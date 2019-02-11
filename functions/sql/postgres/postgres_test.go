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

package sql

import (
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPostgresDemo(t *testing.T) {
	if connectionName == "" {
		t.Skip("Postgres database not configured")
	}
	rr := httptest.NewRecorder()
	PostgresDemo(rr, nil)
	out, err := ioutil.ReadAll(rr.Result().Body)
	if err != nil {
		t.Fatalf("ioutil.ReadAll: %v", err)
	}
	want := "Now:"
	if got := string(out); !strings.Contains(got, want) {
		t.Fatalf("PostgresDemo got %q, want to contain %q", got, want)
	}
}
