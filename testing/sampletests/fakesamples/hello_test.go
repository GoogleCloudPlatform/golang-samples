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

package hello_test

import (
	"testing"

	hello "github.com/GoogleCloudPlatform/golang-samples/testing/sampletests/fakesamples"
)

func TestHello(t *testing.T) {
	if got, want := hello.Hello(), "Hello!"; got != want {
		t.Errorf("hello got %q, want %q", got, want)
	}
}
