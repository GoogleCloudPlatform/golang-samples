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

package main

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/pubsublite"
	"cloud.google.com/go/pubsublite/pscompat"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/GoogleCloudPlatform/golang-samples/pubsublite/internal/psltest"
	"github.com/google/uuid"
)

const (
	region         = "us-east1"
	zone           = "us-east1-c"
	resourcePrefix = "quickstart-subscriber-"
)

func TestQuickstartSubscriber(t *testing.T) {
	tc := testutil.SystemTest(t)
	m := testutil.BuildMain(t)

	if !m.Built() {
		t.Fatalf("failed to build app")
	}

	ctx := context.Background()
	admin, err := pubsublite.NewAdminClient(ctx, region)
	if err != nil {
		t.Fatalf("pubsublite.NewAdminClient: %v", err)
	}
	defer admin.Close()
	psltest.Cleanup(t, admin, tc.ProjectID, region, resourcePrefix, []string{zone})

	resourceID := resourcePrefix + uuid.NewString()
	topicPath := fmt.Sprintf("projects/%s/locations/%s/topics/%s", tc.ProjectID, zone, resourceID)
	psltest.MustCreateTopic(ctx, t, admin, topicPath)
	defer admin.DeleteTopic(ctx, topicPath)

	subscriptionPath := fmt.Sprintf("projects/%s/locations/%s/subscriptions/%s", tc.ProjectID, zone, resourceID)
	psltest.MustCreateSubscription(ctx, t, admin, topicPath, subscriptionPath)
	defer admin.DeleteSubscription(ctx, subscriptionPath)

	publishMessages(ctx, t, topicPath, 10)

	stdOut, stdErr, err := m.Run(nil, 10*time.Minute,
		"--project_id", tc.ProjectID,
		"--zone", zone,
		"--subscription_id", resourceID,
	)

	if err != nil {
		t.Errorf("stdout: %v", string(stdOut))
		t.Errorf("stderr: %v", string(stdErr))
		t.Errorf("execution failed: %v", err)
	}

	if got, want := string(stdOut), "Received 10 messages"; !strings.Contains(got, want) {
		t.Errorf("got %q\nwant to contain %q", got, want)
	}
}

func publishMessages(ctx context.Context, t *testing.T, topicPath string, messageCount int) {
	publisher, err := pscompat.NewPublisherClient(ctx, topicPath)
	if err != nil {
		t.Fatalf("pscompat.NewPublisherClient error: %v", err)
	}
	defer publisher.Stop()

	var results []*pubsub.PublishResult
	for i := 0; i < messageCount; i++ {
		r := publisher.Publish(ctx, &pubsub.Message{
			Data: []byte(fmt.Sprintf("msg-%d", i)),
		})
		results = append(results, r)
	}

	for _, r := range results {
		if _, err := r.Get(ctx); err != nil {
			t.Fatalf("Publish error: %v", err)
		}
	}
	t.Logf("Published %d messages", messageCount)
}
