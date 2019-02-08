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

package newsletter

// [START utility_example_1]
import (
	"testing"

	"google.golang.org/appengine/aetest"
)

func TestWithContext(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	// Run code and tests requiring the context.Context using ctx.
	// [START_EXCLUDE]
	check(t, ctx)
	// [END_EXCLUDE]
}

// [END utility_example_1]

func check(t *testing.T, ctx interface{}) {
}
