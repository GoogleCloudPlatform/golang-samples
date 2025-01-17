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

package counter

import (
	"context"

	"google.golang.org/appengine/delay"
)

func example() {
	ctx := context.Background()

	// [START gae_deferred_tasks]
	var expensiveFunc = delay.Func("some-arbitrary-key", func(ctx context.Context, a string, b int) {
		// Do something expensive!
	})

	// Somewhere else.
	expensiveFunc.Call(ctx, "Hello, world!", 42)
	// [START gae_deferred_tasks]
}
