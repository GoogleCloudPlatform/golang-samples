// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package counter

import (
	"context"

	"google.golang.org/appengine/delay"
)

func example() {
	ctx := context.Background()

	// [START deferred_tasks]
	var expensiveFunc = delay.Func("some-arbitrary-key", func(ctx context.Context, a string, b int) {
		// Do something expensive!
	})

	// Somewhere else.
	expensiveFunc.Call(ctx, "Hello, world!", 42)
	// [END deferred_tasks]
}
