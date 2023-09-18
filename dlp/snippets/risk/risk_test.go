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

// Package risk contains example snippets using the DLP API to create risk jobs.
package risk

import (
	"context"
	"testing"

	"cloud.google.com/go/pubsub"
)

const (
	riskTopicName        = "dlp-risk-test-topic-"
	riskSubscriptionName = "dlp-risk-test-sub-"
)

func cleanupPubsub(t *testing.T, client *pubsub.Client, topicName, subName string) {
	ctx := context.Background()
	topic := client.Topic(topicName)
	if exists, err := topic.Exists(ctx); err != nil {
		t.Logf("Exists: %v", err)
		return
	} else if exists {
		if err := topic.Delete(ctx); err != nil {
			t.Logf("Delete: %v", err)
		}
	}

	s := client.Subscription(subName)
	if exists, err := s.Exists(ctx); err != nil {
		t.Logf("Exists: %v", err)
		return
	} else if exists {
		if err := s.Delete(ctx); err != nil {
			t.Logf("Delete: %v", err)
		}
	}
}
