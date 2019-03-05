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
	"time"

	"cloud.google.com/go/functions/metadata"
)

func TestFiniteRetryPubSub(t *testing.T) {
	tests := []struct {
		name    string
		created time.Time
		wantErr bool
	}{
		{
			name:    "current",
			created: time.Now(),
			wantErr: true,
		},
		{
			// More than 10 seconds in the past.
			name:    "1 day old",
			created: time.Now().AddDate(0, 0, -1),
			wantErr: false,
		},
		{
			// Dates in the future are silently accepted.
			name:    "future",
			created: time.Now().AddDate(0, 0, 1),
			wantErr: true,
		},
	}

	for _, test := range tests {
		meta := &metadata.Metadata{
			EventID:   "event ID",
			Timestamp: test.created,
		}
		ctx := metadata.NewContext(context.Background(), meta)

		err := FiniteRetryPubSub(ctx, PubSubMessage{[]byte("message")})
		gotErr := err != nil
		if gotErr != test.wantErr {
			t.Errorf("FiniteRetryPubSub(%s): got retry(%t), want retry(%t)", test.name, gotErr, test.wantErr)
		}
	}
}
