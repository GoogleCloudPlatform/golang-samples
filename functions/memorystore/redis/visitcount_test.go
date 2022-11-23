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

package visitcount

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	miniredis "github.com/alicebob/miniredis/v2"
)

func TestVisitCount(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("miniredis.Run: %v", err)
	}
	defer s.Close()

	os.Setenv("REDISHOST", s.Host())
	os.Setenv("REDISPORT", s.Port())

	req := httptest.NewRequest("GET", "/", strings.NewReader(""))
	rr := httptest.NewRecorder()

	visitCount(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("VisitCount got status %v, want %v", rr.Code, http.StatusOK)
	}
}
