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

package main

import (
	"context"
	"log"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestMain(t *testing.T) {
	m := testutil.BuildMain(t)
	defer m.Cleanup()

	if !m.Built() {
		t.Errorf("failed to build app")
	}
}

func TestHandler(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{
			name:     "basic test",
			expected: "Incremented sidecar_sample_counter metric!\n",
		},
	}

	for _, test := range tests {
		ctx := context.Background()
		shutdown := setupCounter(ctx)
		defer shutdown(ctx)

		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
			log.Printf("defaulting to port %s", port)
		}

		req := httptest.NewRequest("GET", "http://localhost:"+port, nil)
		rr := httptest.NewRecorder()
		handler(rr, req)
		if rr.Body.String() != test.expected {
			t.Errorf("unexpected output: '%s', expected '%s'", rr.Body.String(), test.expected)
		}
	}
}
