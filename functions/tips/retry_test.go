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
