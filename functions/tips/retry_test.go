// Copyright 2018 Google LLC. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package tips

import (
	"context"
	"testing"
)

func TestRetryPubSub(t *testing.T) {
	misconfigured = true // Ensures MisconfiguredDataClient returns an error.

	if err := RetryPubSub(context.Background(), PubSubMessage{}); err != nil {
		t.Errorf("RetryPubSub: got %v, want nil", err)
	}

	misconfigured = false
	if err := RetryPubSub(context.Background(), PubSubMessage{}); err == nil {
		t.Errorf("RetryPubSub: got nil, want an error")
	}
}
