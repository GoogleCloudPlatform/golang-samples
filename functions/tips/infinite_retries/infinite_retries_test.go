// Copyright 2018 Google LLC. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

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
