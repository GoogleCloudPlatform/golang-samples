// Copyright 2020 Google LLC
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

package security

import (
	"bytes"
	"os"
	"testing"
)

func TestMakeGetRequest(t *testing.T) {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		t.Skip("BASE_URL not set")
	}
	url := baseURL + "/HelloHTTP"
	var b bytes.Buffer
	if err := makeGetRequest(&b, url); err != nil {
		t.Fatalf("makeGetRequest: %v", err)
	}
	got := b.String()
	if got != "Hello, World!" {
		t.Fatalf("got %s, want %s", got, "Hello, World!")
	}
}
