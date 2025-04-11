// Copyright 2021 Google LLC
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

package publish

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	"cloud.google.com/go/pubsublite"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/GoogleCloudPlatform/golang-samples/pubsublite/internal/psltest"
	"github.com/google/uuid"
)

const (
	region       = "europe-west1"
	zone         = "europe-west1-b"
	topicPrefix  = "publish-test-"
	messageCount = 10
)

func TestPublish(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()
	admin, err := pubsublite.NewAdminClient(context.Background(), region)
	if err != nil {
		t.Fatalf("pubsublite.NewAdminClient: %v", err)
	}
	defer admin.Close()
	psltest.Cleanup(t, admin, tc.ProjectID, region, topicPrefix, []string{zone})

	topicID := topicPrefix + uuid.NewString()
	topicPath := fmt.Sprintf("projects/%s/locations/%s/topics/%s", tc.ProjectID, zone, topicID)
	psltest.MustCreateTopic(ctx, t, admin, topicPath)
	defer admin.DeleteTopic(ctx, topicPath)

	t.Run("WithBatchSettings", func(t *testing.T) {
		buf := new(bytes.Buffer)
		err := publishWithBatchSettings(buf, tc.ProjectID, zone, topicID, messageCount)
		if err != nil {
			t.Fatalf("publishWithBatchSettings: %v", err)
		}

		got := buf.String()
		want := fmt.Sprintf("Published %d messages with batch settings", messageCount)
		if !strings.Contains(got, want) {
			t.Errorf("got %q\nwant to contain %q", got, want)
		}
	})

	t.Run("WithOrderingKey", func(t *testing.T) {
		buf := new(bytes.Buffer)
		err := publishWithOrderingKey(buf, tc.ProjectID, zone, topicID, messageCount)
		if err != nil {
			t.Fatalf("publishWithOrderingKey: %v", err)
		}

		got := buf.String()
		want := fmt.Sprintf("Published %d messages with ordering key", messageCount)
		if !strings.Contains(got, want) {
			t.Errorf("got %q\nwant to contain %q", got, want)
		}
	})

	t.Run("WithCustomAttributes", func(t *testing.T) {
		buf := new(bytes.Buffer)
		err := publishWithCustomAttributes(buf, tc.ProjectID, zone, topicID)
		if err != nil {
			t.Fatalf("publishWithCustomAttributes: %v", err)
		}

		got := buf.String()
		want := "Published a message with custom attributes"
		if !strings.Contains(got, want) {
			t.Errorf("got %q\nwant to contain %q", got, want)
		}
	})
}
