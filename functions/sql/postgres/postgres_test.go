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

	if got, want := rr.Body.String(), "Now:"; !strings.Contains(got, want) {
		t.Fatalf("PostgresDemo got %q, want to contain %q", got, want)
	}
}
